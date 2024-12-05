package secrets

import (
	"fmt"
	"strings"
)

var (
	errNoHandler = "unsupported prefix: '%s', no handler found"
	opInvalidRef = "invalid reference, must start with '%s'"
)

func Get(ref string) (string, error) {
	if strings.HasPrefix(ref, opPrefix) {
		return opResolve(ref)
	}

	if strings.HasPrefix(ref, sopsPrefix) {
		return sopsData(ref)
	}

	return "", fmt.Errorf(errNoHandler, ref)
}
