package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

var (
	applyCmd = &cobra.Command{
		Use:   "apply <stack> [component...]",
		Short: "Create or update infrastructure",
		Run:   apply,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

func apply(cmd *cobra.Command, args []string) {
	run(args, false, func(component *schema.Component, executor schema.Executor) {
		err := executor.Apply(component)
		if err != nil {
			log.Fatal(err)
		}
	})
}
