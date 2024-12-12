package schema

import (
	"fmt"
	"path"

	"dario.cat/mergo"
	cp "github.com/otiai10/copy"
)

type (
	Component struct {
		Stack     string                 `json:"stack"`
		Backend   Backend                `json:"backend"`
		Appends   map[string][]string    `json:"appends"`
		Name      string                 `json:"name"`
		Path      string                 `json:"path"`
		Inputs    map[string]interface{} `json:"inputs"`
		Providers map[string]interface{} `json:"providers"`
	}
)

// copy component to workdir if needed
func (c *Component) EnsurePath(config *Config) error {
	if len(config.WorkDir) > 0 {
		dest := path.Join(config.WorkDir, c.Stack, c.Name)
		err := cp.Copy(c.Path, dest)
		if err != nil {
			return err
		}
		c.Path = dest
	}

	return nil
}

// property ref template to resolve later
func (c *Component) PropertyRef(property string) string {
	return fmt.Sprintf(`{{ (state "%s" "%s").%s }}`, c.Stack, c.Name, property)
}

// resolve templates in component
func (c *Component) ResolveVars(stacks *Stacks, executor Executor) error {
	stack, err := stacks.GetStack(c.Stack)
	if err != nil {
		return err
	}

	tdata := map[string]interface{}{"stack": stack.Name, "component": c.Name}
	err = mergo.Merge(&tdata, stack.Options)
	if err != nil {
		return err
	}

	funcMap := map[string]interface{}{
		"state": stateFunc(stacks, executor),
	}

	// template backend
	c.Backend.Config, err = tpl(c.Backend.Config, tdata, funcMap)
	if err != nil {
		return err
	}

	// template vars
	c.Inputs, err = tpl(c.Inputs, tdata, funcMap)
	if err != nil {
		return err
	}

	// template providers
	c.Providers, err = tpl(c.Providers, tdata, funcMap)
	if err != nil {
		return err
	}

	return nil
}

func stateFunc(stacks *Stacks, executor Executor) func(stack, component string) any {
	return func(stack, component string) any {
		refStack, err := stacks.GetStack(stack)
		if err != nil {
			return nil
		}

		refComponent, err := refStack.GetComponent(component)
		if err != nil {
			return nil
		}

		refState, err := executor.Output(refComponent)
		if err != nil {
			return nil
		}

		res := map[string]string{}
		for k, v := range refState {
			res[k] = v.String()
		}

		return res
	}
}
