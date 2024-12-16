package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/exec"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/parser"
)

var (
	kubeCmd = &cobra.Command{
		Use:     "kube [stack]",
		Short:   "Kubeconfig for stack",
		Aliases: []string{"kubeconfig"},
		RunE:    kube,
		Args:    cobra.ExactArgs(1),
	}
)

func init() {
	rootCmd.AddCommand(kubeCmd)
}

func kube(cmd *cobra.Command, args []string) error {
	executor, err := exec.GetExecutor(config)
	if err != nil {
		log.Fatal(err)
	}

	stacks, err := parser.LoadStacks(config.StacksDir)
	if err != nil {
		return err
	}

	stack, err := stacks.GetStack(args[0])
	if stack == nil {
		return err
	}

	kubeconfig, err := stack.Kubeconfig.Render(config, stacks, executor, stack.Name)
	if err != nil {
		return err
	}

	fmt.Println(kubeconfig)
	return nil
}
