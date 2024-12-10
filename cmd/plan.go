package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

var (
	planCmd = &cobra.Command{
		Use:   "plan <stack> [component]",
		Short: "Show changes required by the current configuration",
		Run:   plan,
		Args:  cobra.RangeArgs(1, 2),
	}
)

func init() {
	rootCmd.AddCommand(planCmd)
}

func plan(cmd *cobra.Command, args []string) {
	run(args, false, func(component *schema.Component, executor schema.Executor) {
		_, err := executor.Plan(component)
		if err != nil {
			log.Fatal(err)
		}
	})
}
