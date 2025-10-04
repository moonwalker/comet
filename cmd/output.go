package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

var (
	outputCmd = &cobra.Command{
		Use:   "output <stack> [component]",
		Short: "Show output values from components",
		Run:   output,
		Args:  cobra.RangeArgs(1, 2),
	}
)

func init() {
	rootCmd.AddCommand(outputCmd)
}

func output(cmd *cobra.Command, args []string) {
	run(args, false, func(component *schema.Component, executor schema.Executor) {
		out, err := executor.Output(component)
		if err != nil {
			log.Fatal(err)
		}

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
