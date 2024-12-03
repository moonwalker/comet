package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/cli"
	"github.com/moonwalker/comet/internal/config"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/version"
)

const (
	name = "comet"
	desc = "Tool for provisioning and managing infrastructure"
)

var (
	rootCmd = &cobra.Command{
		Use:               name,
		Short:             desc,
		SilenceErrors:     true,
		SilenceUsage:      true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
)

func init() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Init()

	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	rootCmd.PersistentFlags().StringVar(&config.Filename, "config", config.Filename, "config file")
	rootCmd.PersistentFlags().StringVar(&config.Settings.StacksDir, "dir", config.Settings.StacksDir, "stacks directory")
	rootCmd.ParseFlags(os.Args)

	rootCmd.Version = version.Info()
}

func Execute() {
	if len(os.Args) == 1 {
		cli.PrintStyledText(name)
		fmt.Fprintf(os.Stdout, "\n")
	}

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
