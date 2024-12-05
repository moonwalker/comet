package cmd

import (
	"encoding/json"
	"fmt"
	"path"

	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/cli"
	"github.com/moonwalker/comet/internal/exec"
	"github.com/moonwalker/comet/internal/exec/execintf"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
	"github.com/moonwalker/comet/internal/stacks"
)

var (
	listCmd = &cobra.Command{
		Use:     "list [stack]",
		Short:   "List stacks or components",
		Aliases: []string{"ls"},
		RunE:    list,
		Args:    cobra.MaximumNArgs(1),
	}
	planCmd = &cobra.Command{
		Use:   "plan <stack> [component]",
		Short: "Show changes required by the current configuration",
		Run:   plan,
		Args:  cobra.RangeArgs(1, 2),
	}
	applyCmd = &cobra.Command{
		Use:   "apply <stack> [component]",
		Short: "Create or update infrastructure",
		Run:   apply,
		Args:  cobra.RangeArgs(1, 2),
	}
	destroyCmd = &cobra.Command{
		Use:   "destroy <stack> [component]",
		Short: "Destroy previously-created infrastructure",
		Run:   destroy,
		Args:  cobra.RangeArgs(1, 2),
	}
	outputCmd = &cobra.Command{
		Use:   "output <stack> [component]",
		Short: "Show output values from components",
		Run:   output,
		Args:  cobra.RangeArgs(1, 2),
	}
)

func init() {
	rootCmd.AddCommand(listCmd, planCmd, applyCmd, destroyCmd, outputCmd)
}

func list(cmd *cobra.Command, args []string) error {
	stacks, err := stacks.LoadStacks(config.StacksDir)
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

func destroy(cmd *cobra.Command, args []string) {
	run(args, func(component *schema.Component, executor execintf.Executor) {
		err := executor.Destroy(component)
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

func run(args []string, cb func(*schema.Component, execintf.Executor)) {
	stacks, err := stacks.LoadStacks(config.StacksDir)
	if err != nil {
		log.Fatal(err)
	}

	executor, err := exec.GetExecutor(config)
	if err != nil {
		log.Fatal(err)
	}

	stack, err := stacks.GetStack(args[0])
	if stack == nil {
		log.Fatal(err)
	}

	var componentName string
	if len(args) == 2 {
		componentName = args[1]
	}

	for _, component := range stack.Components {
		if len(componentName) > 0 && componentName != component.Name {
			continue
		}

		if config.UseWorkDir {
			dest := path.Join(config.WorkDir, stack.Name, component.Name)
			err := cp.Copy(component.Path, dest)
			if err != nil {
				log.Fatal(err)
			}
			component.Path = dest
		}

		err = executor.ResolveVars(component, stacks)
		if err != nil {
			log.Fatal(err)
		}

		cb(component, executor)
	}
}
