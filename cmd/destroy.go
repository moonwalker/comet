package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

var (
	destroyCmd = &cobra.Command{
		Use:   "destroy <stack> [component...]",
		Short: "Destroy previously-created infrastructure",
		Run:   destroy,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	rootCmd.AddCommand(destroyCmd)
}

func destroy(cmd *cobra.Command, args []string) {
	run(args, true, func(component *schema.Component, executor schema.Executor) {
		err := executor.Destroy(component)
		if err != nil {
			log.Fatal(err)
		}
	})
}
