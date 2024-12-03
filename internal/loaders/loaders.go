package loaders

import (
	"github.com/moonwalker/comet/internal/schema"
)

type Loader interface {
	Load(path string) (*schema.Stack, error)
}
