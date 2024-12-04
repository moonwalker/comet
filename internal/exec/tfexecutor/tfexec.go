package tfexecutor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/terraform-exec/tfexec"

	"github.com/moonwalker/comet/internal/exec/execintf"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

type executor struct {
	cmd string
}

func init() {
	os.Setenv("TF_IN_AUTOMATION", "true")
}

func NewExecutor(cmd string) (execintf.Executor, error) {
	return &executor{cmd}, nil
}

func (e *executor) ResolveVars(component *schema.Component, stacks *schema.Stacks) error {
	// lookup for component ref in vars
	for _, v := range component.Vars {
		componentRef := schema.TryComponentRefFromJSON(v)
		if componentRef == nil {
			// not component ref, continue
			continue
		}

		// component ref found, resolve it through stacks
		referencedStack, err := stacks.GetStack(componentRef.Stack)
		if err != nil {
			return err
		}
		referencedComponent := referencedStack.ComponentByName(componentRef.Component)
		if referencedComponent != nil {
			referencedComponentState, err := e.Output(referencedComponent)
			if err != nil {
				return err
			}

			if referencedComponentProperty, ok := referencedComponentState[componentRef.Property]; ok {
				referencedComponentValue := referencedComponentProperty.Value
				component.Vars[componentRef.Property] = referencedComponentValue
			}
		}
	}

	return nil
}

func (e *executor) Output(component *schema.Component) (map[string]execintf.OutputMeta, error) {
	log.Debug("output", "component", component.Name)

	tf, err := tfexec.NewTerraform(component.Path, e.cmd)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = tf.Init(ctx, tfexec.Upgrade(false))
	if err != nil {
		return nil, err
	}

	tfoutput, err := tf.Output(ctx)
	if err != nil {
		return nil, err
	}

	if len(tfoutput) == 0 {
		return nil, fmt.Errorf("empty state for: %s, provision that first", component.Name)
	}

	output := make(map[string]execintf.OutputMeta, len(tfoutput))
	for k, v := range tfoutput {
		output[k] = execintf.OutputMeta{
			Sensitive: v.Sensitive,
			Type:      v.Type,
			Value:     v.Value,
		}
	}

	return output, nil
}

func (e *executor) Plan(component *schema.Component) (bool, error) {
	log.Debug("plan", "component", component.Name)

	varsfile, err := prepareProvision(component, false)
	if err != nil {
		return false, err
	}

	tf, err := tfexec.NewTerraform(component.Path, e.cmd)
	if err != nil {
		return false, err
	}

	err = tf.Init(context.Background(), tfexec.Reconfigure(true))
	if err != nil {
		return false, err
	}

	planfile := fmt.Sprintf("%s-%s.planfile", component.Stack, component.Name)
	return tf.Plan(context.Background(), tfexec.VarFile(varsfile), tfexec.Out(planfile))
}

func (e *executor) Apply(component *schema.Component) error {
	log.Debug("apply", "component", component.Name)

	varsfile, err := prepareProvision(component, false)
	if err != nil {
		return err
	}

	tf, err := tfexec.NewTerraform(component.Path, e.cmd)
	if err != nil {
		return err
	}

	err = tf.Init(context.Background(), tfexec.Reconfigure(true))
	if err != nil {
		return err
	}

	return tf.Apply(context.Background(), tfexec.VarFile(varsfile))
}

func (e *executor) Destroy(component *schema.Component) error {
	log.Debug("destroy", "component", component.Name)

	varsfile, err := prepareProvision(component, false)
	if err != nil {
		return err
	}

	tf, err := tfexec.NewTerraform(component.Path, e.cmd)
	if err != nil {
		return err
	}

	err = tf.Init(context.Background(), tfexec.Reconfigure(true))
	if err != nil {
		return err
	}

	return tf.Destroy(context.Background(), tfexec.VarFile(varsfile))
}

func prepareProvision(component *schema.Component, generateBackend bool) (string, error) {
	varsfile := fmt.Sprintf("%s-%s.tfvars.json", component.Stack, component.Name)
	err := writeJSON(component.Vars, component.Path, varsfile)
	if err != nil {
		return "", err
	}

	if generateBackend {
		_, err := writeBackend(component)
		if err != nil {
			return "", err
		}
	}

	return varsfile, nil
}

func writeBackend(component *schema.Component) (string, error) {
	backend := map[string]interface{}{
		"terraform": map[string]interface{}{
			"backend": map[string]interface{}{
				component.Backend.Type: component.Backend.Data,
			},
		},
	}

	backendfile := "backend.tf.json"
	err := writeJSON(backend, component.Path, backendfile)
	if err != nil {
		return "", err
	}
	return backendfile, nil
}

func writeJSON(v any, dir string, filename string) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(dir, filename), b, 0644)
}
