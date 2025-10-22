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

	// Determine which columns to show based on data
	hasOwner := false
	hasCustom := false
	for _, s := range stacks.OrderByName() {
		if s.Metadata != nil {
			if s.Metadata.Owner != "" {
				hasOwner = true
			}
			if s.Metadata.Custom != nil {
				// Check if it's an ordered slice or a map with content
				if orderedSlice, ok := s.Metadata.Custom.([]interface{}); ok && len(orderedSlice) > 0 {
					hasCustom = true
				} else if customMap, ok := s.Metadata.Custom.(map[string]any); ok && len(customMap) > 0 {
					hasCustom = true
				}
			}
		}
	}

	// Build header dynamically
	headers := []string{"stack", "description"}
	if hasOwner && details {
		headers = append(headers, "owner")
	}
	headers = append(headers, "tags")
	if hasCustom && details {
		headers = append(headers, "custom")
	}
	headers = append(headers, "path")

	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetColWidth(20)
	table.SetAutoMergeCells(false)

	for _, s := range stacks.OrderByName() {
		var desc, owner, tags, custom string
		path := s.Path

		if s.Metadata != nil {
			// Description - no truncation, let table wrap
			desc = s.Metadata.Description

			// Owner
			owner = s.Metadata.Owner

			// Tags - show all in details mode
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

			// Custom fields - show as key=value pairs, one per line
			if s.Metadata.Custom != nil {
				customPairs := []string{}

				// Handle ordered slice format (key, value, key, value, ...)
				if orderedSlice, ok := s.Metadata.Custom.([]interface{}); ok {
					for i := 0; i < len(orderedSlice); i += 2 {
						if i+1 < len(orderedSlice) {
							k := fmt.Sprintf("%v", orderedSlice[i])
							v := orderedSlice[i+1]
							customPairs = append(customPairs, fmt.Sprintf("%s=%v", k, v))
						}
					}
				} else if customMap, ok := s.Metadata.Custom.(map[string]any); ok {
					// Fallback to map (unordered)
					for k, v := range customMap {
						customPairs = append(customPairs, fmt.Sprintf("%s=%v", k, v))
					}
				}

				if len(customPairs) > 0 {
					// Join without manual padding - let tablewriter handle it
					custom = strings.Join(customPairs, "\n")
				}
			}
		}

		// Shorten path for display
		if strings.HasPrefix(path, "stacks/") {
			path = strings.TrimPrefix(path, "stacks/")
		}

		// Build row dynamically
		row := []string{s.Name, desc}
		if hasOwner && details {
			row = append(row, owner)
		}
		row = append(row, tags)
		if hasCustom && details {
			row = append(row, custom)
		}
		row = append(row, path)

		table.Append(row)
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
