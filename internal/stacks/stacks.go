package stacks

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/moonwalker/comet/internal/loaders"
	"github.com/moonwalker/comet/internal/loaders/js"
	"github.com/moonwalker/comet/internal/schema"
)

const (
	errNoLoader = "unsupported extension: '%s', no loader found"
)

var (
	jsextensions = []string{
		".js", ".ts",
	}
	extensions  = slices.Concat(jsextensions)
	globpattern = "**/*{" + strings.Join(extensions, ",") + "}"
)

func LoadStacks(dir string) (*schema.Stacks, error) {
	stacks := &schema.Stacks{}

	err := doublestar.GlobWalk(os.DirFS(dir), globpattern, func(p string, d fs.DirEntry) error {
		path := filepath.Join(dir, p)

		loader, err := getLoader(path)
		if err != nil {
			return err
		}

		stack, err := loader.Load(path)
		if err != nil {
			return err
		}

		if stack.Valid() {
			return stacks.AddStack(stack)
		}

		return nil
	})

	return stacks, err
}

func getLoader(path string) (loaders.Loader, error) {
	ext := filepath.Ext(path)

	switch {
	case slices.Contains(jsextensions, ext):
		return js.NewInterpreter()
	}

	return nil, fmt.Errorf(errNoLoader, ext)
}
