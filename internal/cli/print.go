package cli

import (
	"fmt"
	"os"
	"slices"
	"strings"

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
	table.SetHeader([]string{"stack", "type", "path"})
	for _, s := range stacks.OrderByName() {
		table.Append([]string{s.Name, s.Type, s.Path})
	}

	table.Render()
}

func PrintComponentsList(components []*schema.Component) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)

	table.SetHeader([]string{"component", "path", "vars"})
	slices.SortFunc(components, func(a, b *schema.Component) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	for _, c := range components {
		varsList := []string{}
		for k, v := range c.Inputs {
			varsList = append(varsList, k+"="+fmt.Sprintf("%v", v))
		}

		table.Append([]string{c.Name, c.Path, strings.Join(varsList, "\n")})
	}

	table.Render()
}
