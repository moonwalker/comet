package schema

import (
	"bytes"
	"text/template"
)

type (
	Kubeconfig struct {
		Current  int                 `json:"current"`
		Clusters []*KubeconfgCluster `json:"clusters"`
	}

	KubeconfgCluster struct {
		Context        string   `json:"context"`
		Host           string   `json:"host"`
		Cert           string   `json:"cert"`
		ExecApiVersion string   `json:"exec_apiversion"`
		ExecCommand    string   `json:"exec_command"`
		ExecArgs       []string `json:"exec_args"`
	}
)

const (
	KubeconfigDefaultApiVersion = "client.authentication.k8s.io/v1beta1"
)

func (k *Kubeconfig) Render(config *Config, stacks *Stacks, executor Executor, stackName string) (string, error) {
	if len(k.Clusters) == 0 {
		return "", nil
	}

	if k.Current < 0 || k.Current >= len(k.Clusters) {
		k.Current = 0
	}

	for _, c := range k.Clusters {
		if len(c.ExecApiVersion) == 0 {
			c.ExecApiVersion = KubeconfigDefaultApiVersion
		}
	}

	t, err := NewTemplater(config, stacks, executor, stackName)
	if err != nil {
		return "", err
	}

	err = t.Any(k, nil)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("k").Parse(kubeconfigTemplate)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, k)
	if err != nil {
		return "", err
	}

	return b.String(), nil
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
