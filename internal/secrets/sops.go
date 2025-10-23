package secrets

import (
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"strings"

	"github.com/getsops/sops/v3/decrypt"
	"github.com/tidwall/gjson"
	"sigs.k8s.io/yaml"
)

var (
	sopsPrefix = "sops://"
	yamlExts   = []string{".yaml", ".yml"}
)

func sopsData(ref string) (string, error) {
	if !strings.HasPrefix(ref, sopsPrefix) {
		return "", fmt.Errorf(opInvalidRef, sopsPrefix)
	}

	ref = strings.TrimPrefix(ref, sopsPrefix)

	u, err := url.Parse(ref)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(u.Path)
	b, err := decrypt.File(u.Path, ext)
	if err != nil {
		// Detect common SOPS errors and provide helpful messages
		errMsg := err.Error()
		if strings.Contains(errMsg, "no age identity") ||
			strings.Contains(errMsg, "0 successful groups required") ||
			strings.Contains(errMsg, "failed to get the data key") {
			return "", fmt.Errorf("failed to decrypt SOPS file: %w\n\nℹ️  Hint: Age key might be missing. Try running:\n  comet bootstrap\n\nOr set the key manually:\n  export SOPS_AGE_KEY=\"...\"", err)
		}
		return "", fmt.Errorf("failed to decrypt SOPS file: %w", err)
	}

	if slices.Contains(yamlExts, ext) {
		b, err = yaml.YAMLToJSON(b)
		if err != nil {
			return "", err
		}
	}

	frag, _ := strings.CutPrefix(u.Fragment, "/")
	p := strings.ReplaceAll(frag, "/", ".")
	jres := gjson.GetBytes(b, p)

	return jres.String(), nil
}
