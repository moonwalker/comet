package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/moonwalker/comet/internal/cli"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
	"github.com/moonwalker/comet/internal/version"
)

const (
	name = "comet"
	desc = "Cosmic tool for provisioning and managing infrastructure"
)

var (
	defenv  = ".env"
	usrenv  = ".env.local"
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
	initConfig()
	log.SetLevel(config.LogLevel)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "config file")
	rootCmd.PersistentFlags().StringVar(&config.StacksDir, "dir", config.StacksDir, "stacks directory")
	rootCmd.ParseFlags(os.Args)
	rootCmd.Version = version.Info()

	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
}

func initConfig() {
	// .env (default)
	godotenv.Load(defenv)

	// .env.local # local user specific (usually git ignored)
	godotenv.Overload(usrenv)

	viper.SetConfigFile(cfgFile)

	viper.SetDefault("log_level", "INFO")
	viper.SetDefault("tf_command", "tofu")
	viper.SetDefault("stacks_dir", "stacks")
	viper.SetDefault("work_dir", "stacks/_components")
	viper.SetDefault("use_work_dir", false)
	viper.SetDefault("generate_backend", true)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatal(err)
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}
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
