package cfg

import (
	"errors"
	"io/fs"

	"github.com/spf13/viper"

	"github.com/moonwalker/comet/internal/schema"
)

func Read(cfgFile string, config *schema.Config) error {
	viper.SetConfigFile(cfgFile)

	viper.SetDefault("log_level", "INFO")
	viper.SetDefault("tf_command", "tofu")
	viper.SetDefault("stacks_dir", "stacks")
	viper.SetDefault("generate_backend", true)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	return nil
}
