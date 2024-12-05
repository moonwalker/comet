package exec

import (
	"fmt"
	"slices"

	"github.com/moonwalker/comet/internal/exec/execintf"
	"github.com/moonwalker/comet/internal/exec/tfexecutor"
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

func GetExecutor(config *schema.Config) (execintf.Executor, error) {
	switch {
	case slices.Contains(tfCommands, config.Command):
		return tfexecutor.NewExecutor(config)
	}

	return nil, fmt.Errorf(errExecutorNotFound, config.Command)
}
