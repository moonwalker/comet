package exec

import (
	"fmt"
	"slices"

	"github.com/moonwalker/comet/internal/exec/execintf"
	"github.com/moonwalker/comet/internal/exec/tfexecutor"
)

const (
	CmdTofu             = "tofu"
	CmdTerraform        = "terraform"
	errExecutorNotFound = "executor not found for command: %s"
)

var (
	tfCommands = []string{CmdTofu, CmdTerraform}
)

func GetExecutor(cmd string) (execintf.Executor, error) {
	switch {
	case slices.Contains(tfCommands, cmd):
		return tfexecutor.NewExecutor(cmd)
	}

	return nil, fmt.Errorf(errExecutorNotFound, cmd)
}
