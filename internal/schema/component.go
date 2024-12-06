package schema

import (
	"fmt"
	"path"

	cp "github.com/otiai10/copy"
)

type (
	Component struct {
		Stack   *Stack                 `json:"stack"`
		Name    string                 `json:"name"`
		Path    string                 `json:"path"`
		Backend Backend                `json:"backend"`
		Vars    map[string]interface{} `json:"vars"`
	}
)

// copy component to workdir if needed
func (c *Component) EnsurePath(config *Config) error {
	if len(config.WorkDir) > 0 {
		dest := path.Join(config.WorkDir, c.Stack.Name, c.Name)
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
	return fmt.Sprintf(`{{ (state "%s" "%s").%s }}`, c.Stack.Name, c.Name, property)
}

// resolve templates in component backend and vars
func (c *Component) ResolveVars(stacks *Stacks, executor Executor) error {
	tdata := map[string]interface{}{"stack": c.Stack.Name, "component": c.Name}

	funcMap := map[string]interface{}{
		"state": stateFunc(stacks, executor),
	}

	// set backend from stack's backend template
	c.Backend.Data = tpl(c.Stack.Backend.Data, tdata, funcMap)

	// template vars
	c.Vars = tpl(c.Vars, tdata, funcMap)

	return nil
}

func stateFunc(stacks *Stacks, executor Executor) func(stack, component string) map[string]string {
	return func(stack, component string) map[string]string {
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
