package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

	// Expand target path
	targetPath := expandPath(step.Target)

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

	// Write file
	if err := os.WriteFile(targetPath, []byte(value), mode); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Info(fmt.Sprintf("💾 Saved to: %s (mode: %o)", targetPath, mode))
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

// expandPath expands ~ and environment variables in paths
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}
	return os.ExpandEnv(path)
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
