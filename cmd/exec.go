package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/cli"
	"github.com/moonwalker/comet/internal/config"
	"github.com/moonwalker/comet/internal/exec"
	"github.com/moonwalker/comet/internal/exec/execintf"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
	"github.com/moonwalker/comet/internal/stacks"
)

var (
	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List stacks or components",
		Aliases: []string{"ls"},
		Run:     list,
		Args:    cobra.MaximumNArgs(2),
	}
	planCmd = &cobra.Command{
		Use:   "plan",
		Short: "Show changes required by the current configuration",
		Run:   plan,
	}
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Create or update infrastructure",
		Run:   apply,
	}
	outputCmd = &cobra.Command{
		Use:   "output",
		Short: "Shows infrastructure state",
		Run:   output,
	}
	destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy previously-created infrastructure",
		Run:   destroy,
	}
)

func init() {
	rootCmd.AddCommand(listCmd, planCmd, applyCmd, outputCmd, destroyCmd)
}

func list(cmd *cobra.Command, args []string) {
	stacks, err := stacks.LoadStacks(config.Settings.StacksDir)
	if err != nil {
		log.Fatal(err)
	}

	cli.PrintStacksList(stacks)
}

func plan(cmd *cobra.Command, args []string) {
	run(args, func(component *schema.Component, executor execintf.Executor) {
		_, err := executor.Plan(component)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func apply(cmd *cobra.Command, args []string) {
	run(args, func(component *schema.Component, executor execintf.Executor) {
		err := executor.Apply(component)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func output(cmd *cobra.Command, args []string) {
	run(args, func(component *schema.Component, executor execintf.Executor) {
		output, err := executor.Output(component)
		if err != nil {
			log.Fatal(err)
		}
		b, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		log.Info("output", "component", component.Name)
		fmt.Println(string(b))
	})
}

func destroy(cmd *cobra.Command, args []string) {

}

func run(args []string, cb func(*schema.Component, execintf.Executor)) {
	if len(args) < 1 {
		log.Fatal(fmt.Errorf("stack name is required"))
	}

	stacks, err := stacks.LoadStacks(config.Settings.StacksDir)
	if err != nil {
		log.Fatal(err)
	}

	executor, err := exec.GetExecutor(config.Settings.Command)
	if err != nil {
		log.Fatal(err)
	}

	stackName := args[0]
	stack := stacks.GetStack(stackName)
	if stack == nil {
		log.Fatal(fmt.Errorf("stack not found: %s", stackName))
	}

	var componentName string
	if len(args) == 2 {
		componentName = args[1]
	}

	for _, component := range stack.Components {
		if len(componentName) > 0 && componentName != component.Name {
			continue
		}

		err := component.CopyToWorkDir()
		if err != nil {
			log.Fatal(err)
		}

		err = executor.ResolveVars(component, stacks)
		if err != nil {
			log.Fatal(err)
		}

		cb(component, executor)
	}
}
