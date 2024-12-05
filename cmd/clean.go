package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"

	"github.com/moonwalker/comet/internal/log"
)

const (
	cleanShort = "Delete Terraform-related folders and files"
)

var (
	tffiles  = []string{"backend.tf.json", ".terraform", "terraform.tfstate.d", ".terraform.lock.hcl"}
	cleanCmd = &cobra.Command{
		Use:   "clean <stack> [component]",
		Short: cleanShort,
		Long:  cleanShort + "\n\nincluding:\n\n" + strings.Join(tffiles, "\n"),
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
		log.Info("removing", "file", path)
		return os.RemoveAll(path)
	})
}
