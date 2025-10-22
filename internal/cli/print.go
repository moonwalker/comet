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

func PrintStacksList(stacks *schema.Stacks, details bool) {
	table := tablewriter.NewWriter(os.Stdout)
	
	if details {
		table.SetHeader([]string{"stack", "description", "owner", "tags", "path"})
	} else {
		table.SetHeader([]string{"stack", "description", "tags", "path"})
	}
	
	table.SetAutoWrapText(false)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
	})
	
	for _, s := range stacks.OrderByName() {
		var desc, owner, tags string
		
		if s.Metadata != nil {
			// Description - truncate if too long
			desc = s.Metadata.Description
			if !details && len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			
			// Owner
			owner = s.Metadata.Owner
			
			// Tags - show first few
			if len(s.Metadata.Tags) > 0 {
				maxTags := 3
				if details {
					maxTags = len(s.Metadata.Tags)
				}
				
				displayTags := s.Metadata.Tags
				if len(displayTags) > maxTags {
					displayTags = displayTags[:maxTags]
					tags = strings.Join(displayTags, ", ") + "..."
				} else {
					tags = strings.Join(displayTags, ", ")
				}
			}
		}
		
		if details {
			table.Append([]string{s.Name, desc, owner, tags, s.Path})
		} else {
			table.Append([]string{s.Name, desc, tags, s.Path})
		}
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
