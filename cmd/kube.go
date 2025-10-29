package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/exec"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/parser"
)

var (
	save          bool
	kubeconfigCmd = &cobra.Command{
		Use:   "kubeconfig [stack]",
		Short: "Kubeconfig for stack",
		RunE:  kubeconfig,
		Args:  cobra.ExactArgs(1),
	}
)

func init() {
	kubeconfigCmd.Flags().BoolVarP(&save, "save", "s", save, "save kubeconfig")

	rootCmd.AddCommand(kubeconfigCmd)
}

func kubeconfig(cmd *cobra.Command, args []string) error {
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

	// Apply stack-specific environment variables and defer cleanup
	cleanup := stack.ApplyEnvs()
	defer cleanup()

	if save {
		return stack.Kubeconfig.Save(config, stacks, executor, stack.Name)
	}

	return stack.Kubeconfig.Write(os.Stdout, config, stacks, executor, stack.Name)
}
