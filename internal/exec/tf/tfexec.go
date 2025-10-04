package tf

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"

	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

var (
	errCmdNotFound = "command not found: %s"
	errEmptyState  = "empty state for: %s"
	backendFile    = "backend.tf.json"
	providersFile  = "providers_gen.tf"
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
	tf.SetStderr(os.Stderr)

	err = tf.Init(context.Background(), tfexec.Reconfigure(true))
	if err != nil {
		return false, err
	}

	planfile := fmt.Sprintf(planFileFmt, component.Stack, component.Name)
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
	tf.SetStderr(os.Stderr)

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
	tf.SetStderr(os.Stderr)

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
	}

	return output, nil
}

// utils

func prepareProvision(component *schema.Component, generateBackend bool) (string, error) {
	varsfile := fmt.Sprintf(varsFileFmt, component.Stack, component.Name)
	err := writeJSON(component.Inputs, component.Path, varsfile)
	if err != nil {
		return "", err
	}

	if generateBackend {
		err := writeBackendJSON(component)
		if err != nil {
			return "", err
		}
	}

	err = writeProvidersTF(component)
	if err != nil {
		return "", err
	}

	return varsfile, nil
}

func writeBackendJSON(component *schema.Component) error {
	backend := map[string]interface{}{
		"terraform": map[string]interface{}{
			"backend": map[string]interface{}{
				component.Backend.Type: component.Backend.Config,
			},
		},
	}

	err := writeJSON(backend, component.Path, backendFile)
	if err != nil {
		return err
	}

	return nil
}

func writeProviderConfig(i int, pc map[string]interface{}) string {
	sb := strings.Builder{}

	for k, v := range pc {
		if k == "alias" && i > 2 {
			continue
		}
		if m, ok := v.(map[string]interface{}); ok {
			sb.WriteString(strings.Repeat(" ", i))
			sb.WriteString(fmt.Sprintf("%s {", k))
			sb.WriteString("\n")
			pc := writeProviderConfig(i+2, m)
			sb.WriteString(pc)
			sb.WriteString(strings.Repeat(" ", i))
			sb.WriteString("}\n")
			continue
		}
		sb.WriteString(strings.Repeat(" ", i))
		vs := fmt.Sprintf("%s", v)
		if strings.HasPrefix(vs, "data.") ||
			strings.HasPrefix(vs, "module.") ||
			strings.HasPrefix(vs, "local.") ||
			strings.HasPrefix(vs, "var.") {
			sb.WriteString(fmt.Sprintf("%s = %s", k, v))
		} else {
			_, err := base64.StdEncoding.DecodeString(vs)
			if err == nil {
				sb.WriteString(fmt.Sprintf(`%s = base64decode("%s")`, k, v))
			} else {
				sb.WriteString(fmt.Sprintf(`%s = "%s"`, k, v))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func writeProvidersTF(component *schema.Component) error {
	if len(component.Providers) == 0 {
		return nil
	}

	sb := strings.Builder{}

	// Use captured provider dependencies from component
	if len(component.ProviderDependencies) > 0 {
		// Generate enhanced provider configuration with remote state fallbacks
		generateRemoteStateDataSources(&sb, component.ProviderDependencies, component)
		generateLocalFallbacks(&sb, component.ProviderDependencies)
		generateVariableOverrides(&sb, component.ProviderDependencies)

		// Generate providers with enhanced configurations using locals
		for k, v := range component.Providers {
			sb.WriteString(fmt.Sprintf(`provider "%s" {`, k))
			if m, ok := v.(map[string]interface{}); ok {
				pc := writeProviderConfigWithLocals(2, m, component.ProviderDependencies)
				if len(pc) > 0 {
					sb.WriteString("\n")
					sb.WriteString(pc)
				}
			}
			sb.WriteString("}\n\n")
		}
	} else {
		// Use standard provider generation
		for k, v := range component.Providers {
			sb.WriteString(fmt.Sprintf(`provider "%s" {`, k))
			if m, ok := v.(map[string]interface{}); ok {
				pc := writeProviderConfig(2, m)
				if len(pc) > 0 {
					sb.WriteString("\n")
					sb.WriteString(pc)
				}
			}
			sb.WriteString("}\n\n")
		}
	}

	for _, line := range component.Appends["providers"] {
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	s := strings.TrimSpace(sb.String()) + "\n"
	return os.WriteFile(path.Join(component.Path, providersFile), []byte(s), 0644)
}

// Generate remote state data sources
func generateRemoteStateDataSources(sb *strings.Builder, deps map[string]string, component *schema.Component) {
	sb.WriteString("# Auto-generated remote state data sources for component dependencies\n")
	for comp := range deps {
		sb.WriteString(fmt.Sprintf(`data "terraform_remote_state" "%s" {
  backend = "%s"
  config = {`, comp, component.Backend.Type))

		sb.WriteString("\n")
		for k, v := range component.Backend.Config {
			configValue := fmt.Sprintf("%v", v)
			// Replace current component name with dependency component name in the path
			if strings.Contains(configValue, component.Name) {
				configValue = strings.ReplaceAll(configValue, component.Name, comp)
			}
			sb.WriteString(fmt.Sprintf(`    %s = "%s"`, k, configValue))
			sb.WriteString("\n")
		}
		sb.WriteString("  }\n}\n\n")
	}
}

// Generate local variables with try() fallbacks
func generateLocalFallbacks(sb *strings.Builder, deps map[string]string) {
	sb.WriteString("# Locals with safe fallbacks for component dependencies\n")
	sb.WriteString("locals {\n")
	for comp := range deps {
		sb.WriteString(fmt.Sprintf(`  %s_kube_host = try(
    data.terraform_remote_state.%s.outputs.kube_host,
    var.%s_kube_host,
    "https://127.0.0.1"
  )
  %s_kube_cert = try(
    data.terraform_remote_state.%s.outputs.kube_cert,
    var.%s_kube_cert,
    ""
  )
`, comp, comp, comp, comp, comp, comp))
	}
	sb.WriteString("}\n\n")
}

// Generate variable overrides for manual configuration
func generateVariableOverrides(sb *strings.Builder, deps map[string]string) {
	for comp := range deps {
		sb.WriteString(fmt.Sprintf(`# Variables for manual override of %s outputs (optional)
variable "%s_kube_host" {
  description = "Kubernetes cluster host from %s component (auto-detected from remote state)"
  type        = string
  default     = null
}

variable "%s_kube_cert" {
  description = "Kubernetes cluster CA certificate from %s component (auto-detected from remote state)"
  type        = string
  default     = null
}

`, comp, comp, comp, comp, comp))
	}
}

// Enhanced provider config writer that replaces <no value> with local references
func writeProviderConfigWithLocals(i int, pc map[string]interface{}, deps map[string]string) string {
	sb := strings.Builder{}

	for k, v := range pc {
		if k == "alias" && i > 2 {
			continue
		}
		if m, ok := v.(map[string]interface{}); ok {
			sb.WriteString(strings.Repeat(" ", i))
			sb.WriteString(fmt.Sprintf("%s {", k))
			sb.WriteString("\n")
			pc := writeProviderConfigWithLocals(i+2, m, deps)
			sb.WriteString(pc)
			sb.WriteString(strings.Repeat(" ", i))
			sb.WriteString("}\n")
			continue
		}
		sb.WriteString(strings.Repeat(" ", i))
		vs := fmt.Sprintf("%s", v)

		// Replace <no value> with appropriate local references
		if vs == "<no value>" {
			if replacement := getLocalReferenceForProperty(k, deps); replacement != "" {
				sb.WriteString(fmt.Sprintf("%s = %s", k, replacement))
			} else {
				// Fallback to null for unknown properties
				sb.WriteString(fmt.Sprintf("%s = null", k))
			}
		} else if strings.HasPrefix(vs, "data.") ||
			strings.HasPrefix(vs, "module.") ||
			strings.HasPrefix(vs, "local.") ||
			strings.HasPrefix(vs, "var.") {
			sb.WriteString(fmt.Sprintf("%s = %s", k, v))
		} else {
			_, err := base64.StdEncoding.DecodeString(vs)
			if err == nil {
				sb.WriteString(fmt.Sprintf(`%s = base64decode("%s")`, k, v))
			} else {
				sb.WriteString(fmt.Sprintf(`%s = "%s"`, k, v))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// Get appropriate local reference for common Kubernetes provider properties
func getLocalReferenceForProperty(propertyName string, deps map[string]string) string {
	// For now, assume the dependency is 'gke' (most common case)
	// In a more complete implementation, this would be smarter about detecting which dependency provides which properties
	for component := range deps {
		switch propertyName {
		case "host":
			return fmt.Sprintf("local.%s_kube_host", component)
		case "cluster_ca_certificate":
			return fmt.Sprintf("local.%s_kube_cert", component)
		}
	}
	return ""
}

func writeJSON(v any, dir string, filename string) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(dir, filename), b, 0644)
}
