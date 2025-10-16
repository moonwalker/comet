package schema

type Config struct {
	LogLevel        string            `mapstructure:"log_level"`
	Command         string            `mapstructure:"tf_command"`
	StacksDir       string            `mapstructure:"stacks_dir"`
	WorkDir         string            `mapstructure:"work_dir"`
	GenerateBackend bool              `mapstructure:"generate_backend"`
	Env             map[string]string `mapstructure:"env"`
	Bootstrap       []*BootstrapStep  `mapstructure:"bootstrap"`
}

// BootstrapStep represents a single bootstrap operation
type BootstrapStep struct {
	Name     string `mapstructure:"name"`
	Type     string `mapstructure:"type"`     // secret, command, check, etc.
	Source   string `mapstructure:"source"`   // Where to fetch from (op://, sops://, etc.)
	Target   string `mapstructure:"target"`   // Where to save (file path)
	Mode     string `mapstructure:"mode"`     // File permissions (e.g., "0600")
	Command  string `mapstructure:"command"`  // For type: command
	Check    string `mapstructure:"check"`    // Command to check if step is needed
	Optional bool   `mapstructure:"optional"` // Skip if fails
}
