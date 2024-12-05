package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/cli"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/parser"
)

var (
	listCmd = &cobra.Command{
		Use:     "list [stack]",
		Short:   "List stacks or components",
		Aliases: []string{"ls"},
		RunE:    list,
		Args:    cobra.MaximumNArgs(1),
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	stacks, err := parser.LoadStacks(config.StacksDir)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		cli.PrintStacksList(stacks)
		return nil
	}

	stack, err := stacks.GetStack(args[0])
	if stack == nil {
		return err
	}

	comps := stack.Components
	if len(comps) == 0 {
		log.Info("no components found")
		return nil
	}

	cli.PrintComponentsList(comps)
	return nil
}
