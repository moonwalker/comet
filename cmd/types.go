package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/types"
)

var (
	typesCmd = &cobra.Command{
		Use:   "types",
		Short: "Generate TypeScript definitions for IDE support",
		Long: `Generate TypeScript definitions (index.d.ts) in the stacks directory.
This provides autocomplete and type hints when editing stack files in your IDE.`,
		RunE: generateTypes,
	}
)

func init() {
	rootCmd.AddCommand(typesCmd)
}

func generateTypes(cmd *cobra.Command, args []string) error {
	typesPath := filepath.Join(config.StacksDir, "index.d.ts")

	err := os.WriteFile(typesPath, []byte(types.TypeScriptDefinitions), 0644)
	if err != nil {
		return fmt.Errorf("failed to write types file: %w", err)
	}

	log.Info(fmt.Sprintf("Generated TypeScript definitions at %s", typesPath))
	return nil
}
