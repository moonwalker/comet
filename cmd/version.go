package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/cli"
	"github.com/moonwalker/comet/internal/version"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			cli.PrintStyledText(name)
			fmt.Println()
			fmt.Println(name, "version", version.InfoEx())
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Version = version.Info()
}
