package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

var (
	outputJSON bool

	outputCmd = &cobra.Command{
		Use:   "output <stack> [component] [key]",
		Short: "Show output values from components",
		Long: `Show output values from components.

If only stack is provided, shows outputs from all components.
If stack and component are provided, shows outputs from that component.
If stack, component, and key are provided, shows only that specific output value.

The --json flag formats the output as JSON, which can be piped to tools like jq.`,
		Run:  output,
		Args: cobra.RangeArgs(1, 3),
	}
)

func init() {
	outputCmd.Flags().BoolVar(&outputJSON, "json", false, "Output in JSON format")
	rootCmd.AddCommand(outputCmd)
}

func output(cmd *cobra.Command, args []string) {
	// Extract the key filter if provided (3rd argument)
	var keyFilter string
	if len(args) == 3 {
		keyFilter = args[2]
		// Temporarily reduce args to pass correct stack/component to run()
		args = args[:2]
	}

	run(args, false, func(component *schema.Component, executor schema.Executor) {
		out, err := executor.Output(component)
		if err != nil {
			log.Fatal(err)
		}

		// JSON output mode
		if outputJSON {
			// If a specific key is requested, output only that value as JSON
			if keyFilter != "" {
				if v, ok := out[keyFilter]; ok {
					var rawValue interface{}
					if err := json.Unmarshal(v.Value, &rawValue); err == nil {
						jsonBytes, _ := json.MarshalIndent(rawValue, "", "  ")
						fmt.Println(string(jsonBytes))
					} else {
						// Fallback to string
						jsonBytes, _ := json.MarshalIndent(v.String(), "", "  ")
						fmt.Println(string(jsonBytes))
					}
				} else {
					log.Fatal(fmt.Errorf("output key '%s' not found in component '%s'", keyFilter, component.Name))
				}
				return
			}

			// Output all values as JSON object
			result := make(map[string]interface{})
			for k, v := range out {
				var rawValue interface{}
				if err := json.Unmarshal(v.Value, &rawValue); err == nil {
					result[k] = rawValue
				} else {
					result[k] = v.String()
				}
			}
			jsonBytes, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonBytes))
			return
		}

		// Plain text output mode
		// If a specific key is requested, only show that
		if keyFilter != "" {
			if v, ok := out[keyFilter]; ok {
				var rawValue interface{}
				if err := json.Unmarshal(v.Value, &rawValue); err == nil {
					switch val := rawValue.(type) {
					case string:
						fmt.Println(val)
					case []interface{}, map[string]interface{}:
						jsonBytes, _ := json.Marshal(val)
						fmt.Println(string(jsonBytes))
					default:
						fmt.Printf("%v\n", val)
					}
				} else {
					fmt.Println(v.String())
				}
			} else {
				log.Fatal(fmt.Errorf("output key '%s' not found in component '%s'", keyFilter, component.Name))
			}
			return
		}

		// Show all outputs in human-readable format
		for k, v := range out {
			// Try to unmarshal to detect the actual type
			var rawValue interface{}
			if err := json.Unmarshal(v.Value, &rawValue); err == nil {
				// Format based on type
				switch val := rawValue.(type) {
				case string:
					fmt.Printf("%s = \"%s\"\n", k, val)
				case []interface{}:
					// Format array as JSON array
					jsonBytes, _ := json.Marshal(val)
					fmt.Printf("%s = %s\n", k, string(jsonBytes))
				case map[string]interface{}:
					// Format object as JSON
					jsonBytes, _ := json.Marshal(val)
					fmt.Printf("%s = %s\n", k, string(jsonBytes))
				default:
					// Numbers, booleans, etc
					fmt.Printf("%s = %v\n", k, val)
				}
			} else {
				// Fallback to string representation
				fmt.Printf("%s = \"%s\"\n", k, v.String())
			}
		}
	})
}
