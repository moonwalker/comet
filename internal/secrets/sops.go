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
		return "", err
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
