package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/bootstrap"
	"github.com/moonwalker/comet/internal/log"
)

var (
	bootstrapForce bool

	bootstrapCmd = &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap secrets and dependencies",
		Long: `Bootstrap secrets and dependencies from configured sources.

This command fetches secrets (like SOPS keys) from remote sources (1Password, etc.)
and saves them locally for fast access. This is a one-time setup step.

After bootstrapping, your comet commands will be fast since secrets are cached locally.`,
		Run: runBootstrap,
	}

	bootstrapStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show bootstrap status",
		Run:   bootstrapStatus,
	}

	bootstrapClearCmd = &cobra.Command{
		Use:   "clear",
		Short: "Clear bootstrap state",
		Run:   bootstrapClear,
	}
)

func init() {
	rootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.AddCommand(bootstrapStatusCmd)
	bootstrapCmd.AddCommand(bootstrapClearCmd)

	bootstrapCmd.Flags().BoolVarP(&bootstrapForce, "force", "f", false, "Force re-run all steps")
}

func runBootstrap(cmd *cobra.Command, args []string) {
	if len(config.Bootstrap) == 0 {
		log.Info("No bootstrap configuration found in comet.yaml")
		log.Info("\nTo configure bootstrap, add a 'bootstrap' section:")
		log.Info(`
bootstrap:
  - name: sops-key
    type: secret
    source: op://vault/item/field
    target: ~/.config/sops/age/keys.txt
    mode: "0600"
`)
		return
	}

	log.Info(fmt.Sprintf("Bootstrap configuration: %d step(s)", len(config.Bootstrap)))

	runner, err := bootstrap.NewRunner(config, bootstrapForce)
	if err != nil {
		log.Fatal(err)
	}

	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}
}

func bootstrapStatus(cmd *cobra.Command, args []string) {
	if len(config.Bootstrap) == 0 {
		log.Info("No bootstrap configuration found")
		return
	}

	state, err := bootstrap.LoadState()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nBootstrap Configuration:")
	fmt.Printf("  Steps: %d\n", len(config.Bootstrap))
	if !state.LastRun.IsZero() {
		fmt.Printf("  Last run: %s\n", formatTime(state.LastRun))
	}

	fmt.Println("\nStep Status:")

	completed := 0
	for _, step := range config.Bootstrap {
		stepState := state.GetStep(step.Name)

		if stepState == nil {
			fmt.Printf("  ⚪ %-20s Not run yet\n", step.Name)
			continue
		}

		switch stepState.Status {
		case "completed":
			completed++
			fmt.Printf("  ✅ %-20s Completed (%s)\n", step.Name, formatTime(stepState.CompletedAt))
			if step.Target != "" {
				targetExists := fileExists(expandPath(step.Target))
				if targetExists {
					fmt.Printf("     %-20s Target: %s (exists)\n", "", step.Target)
				} else {
					fmt.Printf("     %-20s Target: %s (missing!)\n", "", step.Target)
				}
			}
		case "failed":
			fmt.Printf("  ❌ %-20s Failed (%s)\n", step.Name, formatTime(stepState.LastAttempt))
			fmt.Printf("     %-20s Error: %s\n", "", stepState.Error)
		case "skipped":
			fmt.Printf("  ⏭️  %-20s Skipped\n", step.Name)
		default:
			fmt.Printf("  ⚪ %-20s Unknown status: %s\n", step.Name, stepState.Status)
		}
	}

	fmt.Printf("\nOverall: %d/%d steps completed\n", completed, len(config.Bootstrap))

	if completed < len(config.Bootstrap) {
		fmt.Println("\nRun 'comet bootstrap' to complete setup")
		fmt.Println("Run 'comet bootstrap --force' to re-run all steps")
	}
}

func bootstrapClear(cmd *cobra.Command, args []string) {
	if err := bootstrap.Clear(); err != nil {
		log.Fatal(err)
	}
	log.Info("✅ Bootstrap state cleared")
	log.Info("Run 'comet bootstrap' to set up again")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func expandPath(path string) string {
	if path == "" {
		return ""
	}
	if path[0] == '~' {
		home, _ := os.UserHomeDir()
		return home + path[1:]
	}
	return os.ExpandEnv(path)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}

	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		mins := int(duration.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}
}
