package tf

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/hashicorp/terraform-exec/tfexec"

	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

var (
	errCmdNotFound = "command not found: %s"
	errEmptyState  = "empty state for: %s"
	backendFile    = "backend.tf.json"
	varsFileFmt    = "%s-%s.tfvars.json"
	planFileFmt    = "%s-%s.planfile"
)

type executor struct {
	config *schema.Config
}

func NewExecutor(config *schema.Config) (*executor, error) {
	_, err := exec.LookPath(config.Command)
	if err != nil {
		return nil, fmt.Errorf(errCmdNotFound, config.Command)
	}

	return &executor{config}, nil
}

func (e *executor) Plan(component *schema.Component) (bool, error) {
	log.Debug("plan", "component", component.Name)

	varsfile, err := prepareProvision(component, e.config.GenerateBackend)
	if err != nil {
		return false, err
	}

	tf, err := tfexec.NewTerraform(component.Path, e.config.Command)
	if err != nil {
		return false, err
	}

	tf.SetSkipProviderVerify(true)
	tf.SetStdout(os.Stdout)

	err = tf.Init(context.Background(), tfexec.Reconfigure(true))
	if err != nil {
		return false, err
	}

	planfile := fmt.Sprintf(planFileFmt, component.Stack.Name, component.Name)
	return tf.Plan(context.Background(), tfexec.VarFile(varsfile), tfexec.Out(planfile))
}

func (e *executor) Apply(component *schema.Component) error {
	log.Debug("apply", "component", component.Name)

	varsfile, err := prepareProvision(component, e.config.GenerateBackend)
	if err != nil {
		return err
	}

	tf, err := tfexec.NewTerraform(component.Path, e.config.Command)
	if err != nil {
		return err
	}

	tf.SetSkipProviderVerify(true)
	tf.SetStdout(os.Stdout)

	err = tf.Init(context.Background(), tfexec.Reconfigure(true))
	if err != nil {
		return err
	}

	return tf.Apply(context.Background(), tfexec.VarFile(varsfile))
}

func (e *executor) Destroy(component *schema.Component) error {
	log.Debug("destroy", "component", component.Name)

	varsfile, err := prepareProvision(component, e.config.GenerateBackend)
	if err != nil {
		return err
	}

	tf, err := tfexec.NewTerraform(component.Path, e.config.Command)
	if err != nil {
		return err
	}

	tf.SetSkipProviderVerify(true)
	tf.SetStdout(os.Stdout)

	err = tf.Init(context.Background(), tfexec.Reconfigure(true))
	if err != nil {
		return err
	}

	return tf.Destroy(context.Background(), tfexec.VarFile(varsfile))
}

func (e *executor) Output(component *schema.Component) (map[string]*schema.OutputMeta, error) {
	log.Debug("output", "component", component.Name)

	tf, err := tfexec.NewTerraform(component.Path, e.config.Command)
	if err != nil {
		return nil, err
	}

	tfoutput, err := tf.Output(context.Background())
	if err != nil {
		return nil, err
	}

	if len(tfoutput) == 0 {
		return nil, fmt.Errorf(errEmptyState, component.Name)
	}

	output := make(map[string]*schema.OutputMeta, len(tfoutput))
	for k, v := range tfoutput {
		output[k] = &schema.OutputMeta{
			Sensitive: v.Sensitive,
			Type:      v.Type,
			Value:     v.Value,
		}

		fmt.Println(k, "=", output[k].String())
	}

	return output, nil
}

func prepareProvision(component *schema.Component, generateBackend bool) (string, error) {
	varsfile := fmt.Sprintf(varsFileFmt, component.Stack.Name, component.Name)
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

	err := writeJSON(backend, component.Path, backendFile)
	if err != nil {
		return "", err
	}

	return backendFile, nil
}

func writeJSON(v any, dir string, filename string) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(dir, filename), b, 0644)
}
