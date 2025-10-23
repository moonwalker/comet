package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"filippo.io/age"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
	"github.com/moonwalker/comet/internal/secrets"
)

// Runner executes bootstrap steps
type Runner struct {
	config *schema.Config
	state  *State
	force  bool
}

// NewRunner creates a new bootstrap runner
func NewRunner(config *schema.Config, force bool) (*Runner, error) {
	state, err := LoadState()
	if err != nil {
		return nil, fmt.Errorf("failed to load bootstrap state: %w", err)
	}

	return &Runner{
		config: config,
		state:  state,
		force:  force,
	}, nil
}

// Run executes all bootstrap steps
func (r *Runner) Run() error {
	if len(r.config.Bootstrap) == 0 {
		log.Info("No bootstrap steps configured")
		return nil
	}

	log.Info(fmt.Sprintf("Running %d bootstrap step(s)...", len(r.config.Bootstrap)))

	for _, step := range r.config.Bootstrap {
		if err := r.runStep(step); err != nil {
			if step.Optional {
				log.Warn(fmt.Sprintf("Optional step '%s' failed: %v", step.Name, err))
				continue
			}
			return fmt.Errorf("step '%s' failed: %w", step.Name, err)
		}
	}

	if err := r.state.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	log.Info("✅ Bootstrap complete!")
	return nil
}

// runStep executes a single bootstrap step
func (r *Runner) runStep(step *schema.BootstrapStep) error {
	log.Info(fmt.Sprintf("Step: %s", step.Name))

	// Check if step needs to run
	if !r.shouldRun(step) {
		log.Info(fmt.Sprintf("⏭️  %s: Already completed (skip)", step.Name))
		return nil
	}

	// Execute based on type
	var err error
	switch step.Type {
	case "secret":
		err = r.runSecretStep(step)
	case "command":
		err = r.runCommandStep(step)
	case "check":
		err = r.runCheckStep(step)
	default:
		err = fmt.Errorf("unknown step type: %s", step.Type)
	}

	// Update state
	stepState := &StepState{
		Target:      step.Target,
		LastAttempt: time.Now(),
	}

	if err != nil {
		stepState.Status = "failed"
		stepState.Error = err.Error()
		r.state.UpdateStep(step.Name, stepState)
		return err
	}

	stepState.Status = "completed"
	stepState.CompletedAt = time.Now()
	r.state.UpdateStep(step.Name, stepState)

	return nil
}

// shouldRun determines if a step needs to be executed
func (r *Runner) shouldRun(step *schema.BootstrapStep) bool {
	// Force always runs
	if r.force {
		return true
	}

	// Check state
	stepState := r.state.GetStep(step.Name)
	if stepState == nil || stepState.Status != "completed" {
		return true
	}

	// For secret steps, check if target file exists
	if step.Type == "secret" && step.Target != "" {
		targetPath := expandPath(step.Target)
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			log.Debug("Target file missing, needs bootstrap", "step", step.Name, "target", targetPath)
			return true
		}
	}

	// Check if custom check command passes
	if step.Check != "" {
		cmd := exec.Command("sh", "-c", step.Check)
		if err := cmd.Run(); err != nil {
			log.Debug("Check command failed, needs bootstrap", "step", step.Name)
			return true
		}
	}

	return false
}

// runSecretStep fetches a secret and saves it to a file
func (r *Runner) runSecretStep(step *schema.BootstrapStep) error {
	log.Info(fmt.Sprintf("⏳ Fetching secret from: %s", step.Source))

	start := time.Now()
	value, err := secrets.Get(step.Source)
	if err != nil {
		return fmt.Errorf("failed to fetch secret: %w", err)
	}
	duration := time.Since(start)

	log.Debug("Secret fetched", "duration", duration)

	// Determine target path - use default for SOPS age keys if not specified
	targetPath := step.Target
	if targetPath == "" {
		// Auto-detect default path for common secret types
		if isSopsAgeKeySource(step.Source) {
			targetPath = getDefaultSopsAgePath()
			log.Debug("Using default SOPS age key path", "path", targetPath)
		} else {
			return fmt.Errorf("target path is required for secret type: %s", step.Source)
		}
	}

	// Expand target path (handles ~, env vars, and platform-specific SOPS paths)
	targetPath = expandPath(targetPath)

	// Create parent directory
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	// Parse file mode
	mode := os.FileMode(0600) // Default
	if step.Mode != "" {
		modeInt, err := strconv.ParseUint(step.Mode, 8, 32)
		if err != nil {
			return fmt.Errorf("invalid file mode %s: %w", step.Mode, err)
		}
		mode = os.FileMode(modeInt)
	}

	// Format the value based on the secret type
	var formattedValue string
	if isSopsAgeKeySource(step.Source) {
		// For SOPS age keys, format with public key comment
		formattedValue, err = formatAgeKey(value)
		if err != nil {
			log.Warn("Could not parse age key for formatting, saving as-is", "error", err)
			formattedValue = value
		}
	} else {
		formattedValue = value
	}

	// Ensure value ends with newline (Unix text file convention)
	if len(formattedValue) > 0 && !strings.HasSuffix(formattedValue, "\n") {
		formattedValue = formattedValue + "\n"
	}

	// For SOPS age keys, check if the key already exists in the file
	if isSopsAgeKeySource(step.Source) {
		shouldAppend, err := shouldAppendAgeKey(targetPath, formattedValue)
		if err != nil {
			return fmt.Errorf("failed to check existing keys: %w", err)
		}
		if !shouldAppend {
			log.Info(fmt.Sprintf("ℹ️  Key already exists in: %s", targetPath))
			log.Info(fmt.Sprintf("✅ %s completed (%.2fs)", step.Name, duration.Seconds()))
			return nil
		}

		// Append to existing file
		f, err := os.OpenFile(targetPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, mode)
		if err != nil {
			return fmt.Errorf("failed to open file for append: %w", err)
		}
		defer f.Close()

		if _, err := f.WriteString(formattedValue); err != nil {
			return fmt.Errorf("failed to append to file: %w", err)
		}

		log.Info(fmt.Sprintf("💾 Appended to: %s (mode: %o)", targetPath, mode))
	} else {
		// For non-age-key secrets, just write/overwrite
		if err := os.WriteFile(targetPath, []byte(formattedValue), mode); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		log.Info(fmt.Sprintf("💾 Saved to: %s (mode: %o)", targetPath, mode))
	}

	log.Info(fmt.Sprintf("✅ %s completed (%.2fs)", step.Name, duration.Seconds()))

	return nil
}

// runCommandStep executes a shell command
func (r *Runner) runCommandStep(step *schema.BootstrapStep) error {
	log.Info(fmt.Sprintf("⏳ Executing: %s", step.Command))

	cmd := exec.Command("sh", "-c", step.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	log.Info(fmt.Sprintf("✅ %s completed", step.Name))
	return nil
}

// runCheckStep checks if required binaries exist
func (r *Runner) runCheckStep(step *schema.BootstrapStep) error {
	// Parse command as comma-separated binary names
	binaries := strings.Split(step.Command, ",")

	var missing []string
	for _, binary := range binaries {
		binary = strings.TrimSpace(binary)
		if _, err := exec.LookPath(binary); err != nil {
			missing = append(missing, binary)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required binaries: %s", strings.Join(missing, ", "))
	}

	log.Info(fmt.Sprintf("✅ All required binaries found"))
	return nil
}

// expandPath expands ~ and environment variables in paths.
// It also handles SOPS age key path resolution on macOS.
func expandPath(path string) string {
	// Handle special case for SOPS age keys on macOS
	// SOPS uses different default paths depending on XDG_CONFIG_HOME:
	// - If XDG_CONFIG_HOME is set: $XDG_CONFIG_HOME/sops/age/keys.txt
	// - On macOS without XDG_CONFIG_HOME: ~/Library/Application Support/sops/age/keys.txt
	// - On Linux without XDG_CONFIG_HOME: ~/.config/sops/age/keys.txt
	if strings.Contains(path, "sops/age/keys.txt") {
		resolvedPath := resolveSopsAgePath(path)
		if resolvedPath != "" {
			return resolvedPath
		}
	}

	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}
	return os.ExpandEnv(path)
}

// resolveSopsAgePath resolves the SOPS age key path to match what the SOPS library expects.
// This ensures bootstrap saves the key where SOPS will actually look for it.
func resolveSopsAgePath(path string) string {
	const sopsAgeKeyPath = "sops/age/keys.txt"

	// If path doesn't contain the SOPS age key path, don't modify it
	if !strings.Contains(path, sopsAgeKeyPath) {
		return ""
	}

	// Check if XDG_CONFIG_HOME is set
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		// Use XDG_CONFIG_HOME if set (works on all platforms)
		return filepath.Join(xdgConfigHome, sopsAgeKeyPath)
	}

	// Platform-specific defaults when XDG_CONFIG_HOME is not set
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	if runtime.GOOS == "darwin" {
		// macOS: ~/Library/Application Support/sops/age/keys.txt
		// This matches what os.UserConfigDir() returns on macOS
		return filepath.Join(home, "Library", "Application Support", sopsAgeKeyPath)
	}

	// Linux/others: ~/.config/sops/age/keys.txt
	return filepath.Join(home, ".config", sopsAgeKeyPath)
}

// isSopsAgeKeySource checks if the source is likely a SOPS age key
func isSopsAgeKeySource(source string) bool {
	// Common patterns for SOPS age keys in secret managers
	lower := strings.ToLower(source)
	return strings.Contains(lower, "sops") &&
		(strings.Contains(lower, "age") || strings.Contains(lower, "key"))
}

// getDefaultSopsAgePath returns the default SOPS age key path for the current platform
func getDefaultSopsAgePath() string {
	// Check if XDG_CONFIG_HOME is set
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "sops", "age", "keys.txt")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to a reasonable default
		return "~/.config/sops/age/keys.txt"
	}

	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "sops", "age", "keys.txt")
	}

	return filepath.Join(home, ".config", "sops", "age", "keys.txt")
}

// NeedsBootstrap checks if any bootstrap steps need to be run
func NeedsBootstrap(config *schema.Config) bool {
	if len(config.Bootstrap) == 0 {
		return false
	}

	runner, err := NewRunner(config, false)
	if err != nil {
		return true // Assume needs bootstrap if we can't check
	}

	for _, step := range config.Bootstrap {
		if runner.shouldRun(step) {
			return true
		}
	}

	return false
}

// formatAgeKey formats an age secret key with a public key comment
func formatAgeKey(secretKey string) (string, error) {
	// Remove whitespace
	secretKey = strings.TrimSpace(secretKey)

	// Parse the age identity to get the public key
	identity, err := age.ParseX25519Identity(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse age identity: %w", err)
	}

	// Get the recipient (public key)
	recipient := identity.Recipient()

	// Format with public key comment
	formatted := fmt.Sprintf("# public key: %s\n%s", recipient.String(), secretKey)

	return formatted, nil
}

// shouldAppendAgeKey checks if an age key should be appended to the file
// Returns true if the key doesn't exist, false if it already exists
func shouldAppendAgeKey(filePath, formattedKey string) (bool, error) {
	// Parse the new key to get its public key
	newIdentity, err := age.ParseX25519Identity(strings.TrimSpace(extractSecretKey(formattedKey)))
	if err != nil {
		return false, fmt.Errorf("failed to parse new key: %w", err)
	}
	newPublicKey := newIdentity.Recipient().String()

	// Check if file exists
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, we should create it
			return true, nil
		}
		return false, err
	}
	defer file.Close()

	// Parse all existing identities from the file
	existingIdentities, err := age.ParseIdentities(file)
	if err != nil {
		// If file exists but has no valid keys (or is empty), we can append
		if strings.Contains(err.Error(), "no secret keys found") {
			return true, nil
		}
		return false, fmt.Errorf("failed to parse existing keys: %w", err)
	}

	// Check if our public key already exists
	for _, identity := range existingIdentities {
		if x25519Identity, ok := identity.(*age.X25519Identity); ok {
			existingPublicKey := x25519Identity.Recipient().String()
			if existingPublicKey == newPublicKey {
				// Key already exists
				return false, nil
			}
		}
	}

	// Key doesn't exist, we should append
	return true, nil
}

// extractSecretKey extracts just the AGE-SECRET-KEY line from formatted content
func extractSecretKey(content string) string {
	re := regexp.MustCompile(`(?m)^AGE-SECRET-KEY-[A-Z0-9]+$`)
	match := re.FindString(content)
	return match
}
