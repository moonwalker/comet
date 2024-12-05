package cmd

import (
	"github.com/moonwalker/comet/internal/exec"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/parser"
	"github.com/moonwalker/comet/internal/schema"
)

func run(args []string, cb func(*schema.Component, schema.Executor)) {
	executor, err := exec.GetExecutor(config)
	if err != nil {
		log.Fatal(err)
	}

	stacks, err := parser.LoadStacks(config.StacksDir)
	if err != nil {
		log.Fatal(err)
	}

	stack, err := stacks.GetStack(args[0])
	if err != nil {
		log.Fatal(err)
	}

	var componentName string
	if len(args) == 2 {
		componentName = args[1]
	}

	components, err := stack.GetComponents(componentName)
	if err != nil {
		log.Fatal(err)
	}

	for _, component := range components {
		err := component.EnsurePath(config)
		if err != nil {
			log.Fatal(err)
		}

		err = component.ResolveVars(stacks, executor)
		if err != nil {
			log.Fatal(err)
		}

		cb(component, executor)
	}
}
