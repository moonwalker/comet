package schema

import (
	"fmt"
	"path"

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
func (c *Component) EnsurePath(config *Config, copy bool) error {
	if len(config.WorkDir) > 0 {
		dest := path.Join(config.WorkDir, c.Stack, c.Name)
		if copy {
			err := cp.Copy(c.Path, dest)
			if err != nil {
				return err
			}
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
func (c *Component) ResolveVars(config *Config, stacks *Stacks, executor Executor) error {
	tdata := map[string]interface{}{
		"component": c.Name,
	}

	t, err := NewTemplater(config, stacks, executor, c.Stack)
	if err != nil {
		return err
	}

	// template backend
	c.Backend.Config, err = t.Map(c.Backend.Config, tdata)
	if err != nil {
		return err
	}

	// template vars
	c.Inputs, err = t.Map(c.Inputs, tdata)
	if err != nil {
		return err
	}

	// template providers
	c.Providers, err = t.Map(c.Providers, tdata)
	if err != nil {
		return err
	}

	return nil
}
