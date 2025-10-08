package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

var (
	initCmd = &cobra.Command{
		Use:   "init <stack> [component...]",
		Short: "Initialize backends and providers",
		Long: `Initialize Terraform/OpenTofu backends and download required providers.

This command prepares components for use without running plan/apply operations.
Useful for read-only operations like 'comet output' or troubleshooting.`,
		Run:  initialize,
		Args: cobra.MinimumNArgs(1),
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}

func initialize(cmd *cobra.Command, args []string) {
	run(args, false, func(component *schema.Component, executor schema.Executor) {
		err := executor.Init(component)
		if err != nil {
			log.Fatal(err)
		}
	})
}
