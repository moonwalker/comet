package exec

import (
	"fmt"
	"slices"

	"github.com/moonwalker/comet/internal/exec/tf"
	"github.com/moonwalker/comet/internal/schema"
)

const (
	CmdTofu             = "tofu"
	CmdTerraform        = "terraform"
	errExecutorNotFound = "executor not found for command: %s"
)

var (
	tfCommands = []string{CmdTofu, CmdTerraform}
)

func GetExecutor(config *schema.Config) (schema.Executor, error) {
	switch {
	case slices.Contains(tfCommands, config.Command):
		return tf.NewExecutor(config)
	}

	return nil, fmt.Errorf(errExecutorNotFound, config.Command)
}
