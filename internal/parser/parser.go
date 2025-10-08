package parser

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/moonwalker/comet/internal/parser/js"
	"github.com/moonwalker/comet/internal/schema"
	"github.com/moonwalker/comet/internal/types"
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

	// Ensure TypeScript definitions exist once for the entire stacks directory
	ensureTypeScriptDefinitions(dir)

	err := doublestar.GlobWalk(os.DirFS(dir), globpattern, func(p string, d fs.DirEntry) error {
		// Skip TypeScript definition files
		if strings.HasSuffix(p, ".d.ts") {
			return nil
		}

		path := filepath.Join(dir, p)

		parser, err := getParser(path)
		if err != nil {
			return err
		}

		stack, err := parser.Parse(path)
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

func ensureTypeScriptDefinitions(stacksDir string) {
	typesPath := filepath.Join(stacksDir, "index.d.ts")
	if _, err := os.Stat(typesPath); os.IsNotExist(err) {
		// Silently write the types file if it doesn't exist
		_ = os.WriteFile(typesPath, []byte(types.TypeScriptDefinitions), 0644)
	}
}

func getParser(path string) (schema.Parser, error) {
	ext := filepath.Ext(path)

	switch {
	case slices.Contains(jsextensions, ext):
		return js.NewInterpreter()
	}

	return nil, fmt.Errorf(errNoLoader, ext)
}
