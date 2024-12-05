package secrets

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/1password/onepassword-sdk-go"
)

const (
	opPrefix = "op://"
)

var (
	opOnce    sync.Once
	opClient  *onepassword.Client
	opOnceErr error
)

func opResolve(ref string) (string, error) {
	if !strings.HasPrefix(ref, opPrefix) {
		return "", fmt.Errorf(opInvalidRef, opPrefix)
	}

	if opClient == nil {
		opOnce.Do(func() {
			opClient, opOnceErr = onepassword.NewClient(
				context.Background(),
				onepassword.WithServiceAccountToken(os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")),
				onepassword.WithIntegrationInfo("comet op integration", "v0.1.0"),
			)
		})
	}
	if opOnceErr != nil {
		return "", opOnceErr
	}

	item, err := opClient.Secrets.Resolve(context.Background(), ref)
	if err != nil {
		return "", err
	}

	return item, nil
}
