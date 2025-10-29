package schema

import (
	"bytes"
	"strings"
	"testing"
	"text/template"
)

func TestKubeconfigTemplate(t *testing.T) {
	tests := []struct {
		name     string
		config   *Kubeconfig
		wantAuth string // Expected auth section in output
	}{
		{
			name: "token authentication",
			config: &Kubeconfig{
				Current: 0,
				Clusters: []*KubeconfgCluster{
					{
						Context:        "test-cluster",
						Host:           "https://kubernetes.example.com",
						Cert:           "LS0tLS1CRUdJTi1DRVJUSUZJQ0FURS0tLS0t",
						Token:          "test-token-12345",
						ExecApiVersion: "client.authentication.k8s.io/v1beta1",
					},
				},
			},
			wantAuth: "token: test-token-12345",
		},
		{
			name: "exec authentication",
			config: &Kubeconfig{
				Current: 0,
				Clusters: []*KubeconfgCluster{
					{
						Context:        "test-cluster",
						Host:           "https://kubernetes.example.com",
						Cert:           "LS0tLS1CRUdJTi1DRVJUSUZJQ0FURS0tLS0t",
						ExecCommand:    "kubectl",
						ExecArgs:       []string{"get-token"},
						ExecApiVersion: "client.authentication.k8s.io/v1beta1",
					},
				},
			},
			wantAuth: "command: kubectl",
		},
		{
			name: "exec takes priority over token",
			config: &Kubeconfig{
				Current: 0,
				Clusters: []*KubeconfgCluster{
					{
						Context:        "test-cluster",
						Host:           "https://kubernetes.example.com",
						Cert:           "LS0tLS1CRUdJTi1DRVJUSUZJQ0FURS0tLS0t",
						Token:          "test-token-12345",
						ExecCommand:    "kubectl",
						ExecArgs:       []string{"get-token"},
						ExecApiVersion: "client.authentication.k8s.io/v1beta1",
					},
				},
			},
			wantAuth: "command: kubectl",
		},
		{
			name: "no authentication",
			config: &Kubeconfig{
				Current: 0,
				Clusters: []*KubeconfgCluster{
					{
						Context:        "test-cluster",
						Host:           "https://kubernetes.example.com",
						Cert:           "LS0tLS1CRUdJTi1DRVJUSUZJQ0FURS0tLS0t",
						ExecApiVersion: "client.authentication.k8s.io/v1beta1",
					},
				},
			},
			wantAuth: "", // No auth expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Normalize exec args
			for _, c := range tt.config.Clusters {
				c.ExecArgs = normalizeExecArgs(c.ExecArgs)
			}

			// Render template
			tmpl, err := template.New("k").Parse(kubeconfigTemplate)
			if err != nil {
				t.Fatalf("template.Parse() error = %v", err)
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, tt.config)
			if err != nil {
				t.Fatalf("template.Execute() error = %v", err)
			}

			output := buf.String()

			// Check for expected auth if specified
			if tt.wantAuth != "" {
				if !strings.Contains(output, tt.wantAuth) {
					t.Errorf("output missing expected auth: %q\nGot:\n%s", tt.wantAuth, output)
				}
			}

			// Ensure valid YAML structure
			if !strings.Contains(output, "apiVersion: v1") {
				t.Errorf("output missing apiVersion")
			}
			if !strings.Contains(output, "kind: Config") {
				t.Errorf("output missing kind")
			}
			if !strings.Contains(output, "current-context: test-cluster") {
				t.Errorf("output missing current-context")
			}
		})
	}
}

func TestNormalizeExecArgs(t *testing.T) {
	tests := []struct {
		name string
		args interface{}
		want []string
	}{
		{
			name: "nil args",
			args: nil,
			want: nil,
		},
		{
			name: "string slice",
			args: []string{"arg1", "arg2"},
			want: []string{"arg1", "arg2"},
		},
		{
			name: "interface slice",
			args: []interface{}{"arg1", "arg2", 123},
			want: []string{"arg1", "arg2", "123"},
		},
		{
			name: "JSON array string",
			args: `["arg1", "arg2"]`,
			want: []string{"arg1", "arg2"},
		},
		{
			name: "Go slice string",
			args: "[arg1 arg2 arg3]",
			want: []string{"arg1", "arg2", "arg3"},
		},
		{
			name: "single string",
			args: "single-arg",
			want: []string{"single-arg"},
		},
		{
			name: "empty string",
			args: "",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeExecArgs(tt.args)
			if len(got) != len(tt.want) {
				t.Errorf("normalizeExecArgs() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("normalizeExecArgs()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
