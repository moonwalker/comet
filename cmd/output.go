package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

var (
	outputCmd = &cobra.Command{
		Use:   "output <stack> [component]",
		Short: "Show output values from components",
		Run:   output,
		Args:  cobra.RangeArgs(1, 2),
	}
)

func init() {
	rootCmd.AddCommand(outputCmd)
}

func output(cmd *cobra.Command, args []string) {
	run(args, false, func(component *schema.Component, executor schema.Executor) {
		_, err := executor.Output(component)
		if err != nil {
			log.Fatal(err)
		}
	})
}
