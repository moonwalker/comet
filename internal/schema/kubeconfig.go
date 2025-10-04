package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	clientauthentication "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/moonwalker/comet/internal/log"
)

type (
	Kubeconfig struct {
		Current  int                 `json:"current"`
		Clusters []*KubeconfgCluster `json:"clusters"`
	}

	KubeconfgCluster struct {
		Context        string      `json:"context"`
		Host           string      `json:"host"`
		Cert           string      `json:"cert"`
		ExecApiVersion string      `json:"exec_apiversion"`
		ExecCommand    string      `json:"exec_command"`
		ExecArgs       interface{} `json:"exec_args"` // Can be []string or string (template)
	}
)

func (k *Kubeconfig) Save(config *Config, stacks *Stacks, executor Executor, stackName string) error {
	pathOptions := clientcmd.NewDefaultPathOptions()

	kubeconfig, err := pathOptions.GetStartingConfig()
	if err != nil {
		return err
	}

	// write out stack's kubeconfig
	var b bytes.Buffer
	err = k.Write(&b, config, stacks, executor, stackName)
	if err != nil {
		return err
	}

	// load stack's kubeconfig
	stackconfig, err := clientcmd.Load(b.Bytes())
	if err != nil {
		return err
	}

	err = mergeKubeconfig(stackconfig, kubeconfig, true)
	if err != nil {
		return err
	}

	return clientcmd.ModifyConfig(pathOptions, *kubeconfig, false)
}

func (k *Kubeconfig) Write(out io.Writer, config *Config, stacks *Stacks, executor Executor, stackName string) error {
	if len(k.Clusters) == 0 {
		return nil
	}

	if k.Current < 0 || k.Current >= len(k.Clusters) {
		k.Current = 0
	}

	for _, c := range k.Clusters {
		if len(c.ExecApiVersion) == 0 {
			c.ExecApiVersion = clientauthentication.SchemeGroupVersion.String()
		}
	}

	t, err := NewTemplater(config, stacks, executor, stackName)
	if err != nil {
		return err
	}

	err = t.Any(k, nil)
	if err != nil {
		return err
	}

	// After templating, convert ExecArgs to []string if needed
	for _, c := range k.Clusters {
		c.ExecArgs = normalizeExecArgs(c.ExecArgs)
	}

	tmpl, err := template.New("k").Parse(kubeconfigTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(out, k)
}

// normalizeExecArgs converts ExecArgs to []string regardless of input type
func normalizeExecArgs(args interface{}) []string {
	if args == nil {
		return nil
	}

	switch v := args.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, len(v))
		for i, arg := range v {
			result[i] = fmt.Sprintf("%v", arg)
		}
		return result
	case string:
		// If it's a JSON array string, try to unmarshal it
		var arr []string
		if err := json.Unmarshal([]byte(v), &arr); err == nil {
			return arr
		}
		// Otherwise return as single element
		if v != "" {
			return []string{v}
		}
		return nil
	default:
		return []string{fmt.Sprintf("%v", v)}
	}
}

// mergeKubeconfig merges a remote cluster's config file with a local config file,
// assuming that the current context in the remote config file points to the
// cluster details to add to the local config
func mergeKubeconfig(remote, local *clientcmdapi.Config, setCurrentContext bool) error {
	remoteCtx, ok := remote.Contexts[remote.CurrentContext]
	if !ok {
		return fmt.Errorf("config has no context entry named %q", remote.CurrentContext)
	}

	remoteCluster, ok := remote.Clusters[remoteCtx.Cluster]
	if !ok {
		return fmt.Errorf("config has no cluster entry named %q", remoteCtx.Cluster)
	}

	remoteAuthInfo, ok := remote.AuthInfos[remoteCtx.AuthInfo]
	if !ok {
		return fmt.Errorf("config has no auth entry named %q", remoteCtx.AuthInfo)
	}

	local.Contexts[remote.CurrentContext] = remoteCtx
	local.Clusters[remoteCtx.Cluster] = remoteCluster
	local.AuthInfos[remoteCtx.AuthInfo] = remoteAuthInfo

	if setCurrentContext {
		log.Debug("setting current kube context to %q", remote.CurrentContext)
		local.CurrentContext = remote.CurrentContext
	}

	return nil
}

const kubeconfigTemplate = `apiVersion: v1
kind: Config
current-context: {{ (index .Clusters .Current).Context }}
contexts:
{{- range .Clusters }}
  - name: {{ .Context }}
    context:
      cluster: {{ .Context }}
      user: {{ .Context }}
{{- end }}
clusters:
{{- range .Clusters }}
  - name: {{ .Context }}
    cluster:
      server: {{ .Host }}
      certificate-authority-data: {{ .Cert }}
{{- end }}
users:
{{- range .Clusters }}
  - name: {{ .Context }}
    user:
      exec:
        apiVersion: {{ .ExecApiVersion }}
        command: {{ .ExecCommand }}
		{{- if .ExecArgs }}
        args:
		{{- range .ExecArgs }}
        - {{ . }}
        {{-  end }}
        {{- end }}
{{- end }}`
