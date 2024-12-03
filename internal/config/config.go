package config

import (
	"errors"
	"io/fs"

	"github.com/spf13/viper"
)

type config struct {
	LogLevel   string `mapstructure:"log_level"`
	Command    string `mapstructure:"tf_command"`
	StacksDir  string `mapstructure:"stacks_dir"`
	WorkDir    string `mapstructure:"work_dir"`
	UseWorkDir bool   `mapstructure:"use_work_dir"`
}

var (
	Filename = "comet.yaml"
	Settings = &config{}
)

func Init() error {
	viper.SetConfigFile(Filename)

	viper.SetDefault("log_level", "INFO")
	viper.SetDefault("tf_command", "tofu")
	viper.SetDefault("stacks_dir", "stacks")
	viper.SetDefault("use_work_dir", false)

	viper.SetDefault("work_dir", "stacks/_components")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	err = viper.Unmarshal(&Settings)
	if err != nil {
		return err
	}

	return nil
}
