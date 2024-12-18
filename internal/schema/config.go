package schema

type Config struct {
	LogLevel        string `mapstructure:"log_level"`
	Command         string `mapstructure:"tf_command"`
	StacksDir       string `mapstructure:"stacks_dir"`
	WorkDir         string `mapstructure:"work_dir"`
	GenerateBackend bool   `mapstructure:"generate_backend"`
}
