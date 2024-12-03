package cli

import (
	"os"

	"github.com/arsham/figurine/figurine"
	"github.com/jwalton/go-supportscolor"
	"github.com/olekukonko/tablewriter"

	"github.com/moonwalker/comet/internal/schema"
)

func PrintStyledText(text string) error {
	if supportscolor.Stdout().SupportsColor {
		return figurine.Write(os.Stdout, text, "Sub-Zero.flf")
	}
	return nil
}

func PrintStacksList(stacks *schema.Stacks) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "type", "path"})
	for _, s := range stacks.OrderByName() {
		table.Append([]string{s.Name, s.Type, s.Path})
	}

	table.Render()
}
