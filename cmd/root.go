package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/cfg"
	"github.com/moonwalker/comet/internal/cli"
	"github.com/moonwalker/comet/internal/env"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
	"github.com/moonwalker/comet/internal/secrets"
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
	loadConfigEnv(config)
	log.SetLevel(config.LogLevel)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "config file")
	rootCmd.PersistentFlags().StringVar(&config.StacksDir, "dir", config.StacksDir, "stacks directory")
	rootCmd.ParseFlags(os.Args)

	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
}

func loadConfigEnv(config *schema.Config) {
	for key, value := range config.Env {
		// Skip if already set in shell environment (shell wins)
		if os.Getenv(key) != "" {
			continue
		}

		// Resolve secrets if value starts with op:// or sops://
		if strings.HasPrefix(value, "op://") || strings.HasPrefix(value, "sops://") {
			resolved, err := secrets.Get(value)
			if err != nil {
				log.Error(fmt.Sprintf("failed to resolve env var %s: %v", key, err))
				continue
			}
			os.Setenv(key, resolved)
		} else {
			// Plain value
			os.Setenv(key, value)
		}
	}
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
