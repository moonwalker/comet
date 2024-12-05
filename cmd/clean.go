package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
)

var (
	tffiles  = []string{"backend.tf.json", ".terraform", "terraform.tfstate.d", ".terraform.lock.hcl"}
	cleanCmd = &cobra.Command{
		Use:   "clean <stack> [component]",
		Short: "Delete Terraform-related folders and files",
		Long:  "Delete Terraform-related folders and files\n\nincluding:\n\n" + strings.Join(tffiles, "\n"),
		RunE:  clean,
		Args:  cobra.RangeArgs(1, 2),
	}
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func clean(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	globpattern := "**/{" + strings.Join(tffiles, ",") + "}"
	return doublestar.GlobWalk(os.DirFS(dir), globpattern, func(p string, d fs.DirEntry) error {
		path := filepath.Join(dir, p)

		fmt.Println("Deleting", path)

		return nil
	})
}
