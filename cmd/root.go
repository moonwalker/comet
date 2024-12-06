package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/cfg"
	"github.com/moonwalker/comet/internal/cli"
	"github.com/moonwalker/comet/internal/env"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

const (
	name = "comet"
	desc = "Cosmic tool for provisioning and managing infrastructure"
)

var (
	cfgFile = "comet.yaml"
	config  = &schema.Config{}
	rootCmd = &cobra.Command{
		Use:               name,
		Short:             desc,
		SilenceErrors:     true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
)

func init() {
	env.Load()
	cfg.Read(cfgFile, config)
	log.SetLevel(config.LogLevel)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "config file")
	rootCmd.PersistentFlags().StringVar(&config.StacksDir, "dir", config.StacksDir, "stacks directory")
	rootCmd.ParseFlags(os.Args)

	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
}

func Execute() {
	if len(os.Args) == 1 {
		cli.PrintStyledText(name)
		fmt.Println()
	}

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
