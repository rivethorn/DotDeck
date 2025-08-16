package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rivethorn/dotdeck/internal"
	"github.com/rivethorn/dotdeck/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(linkCmd)
	linkCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate actions without making changes")
}

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link all files from config.toml",
	RunE: func(cmd *cobra.Command, args []string) error {
		internal.LogVerbose(verbose, "Checking for config file")
		cfg, err := config.Load("config.toml")
		if err != nil {
			return err
		}
		for src, dest := range cfg.Files {
			absSrc, _ := filepath.Abs(src)
			destPath := expandPath(dest)
			internal.LogVerbose(verbose, "Checking if %s exists", destPath)
			if _, err := os.Lstat(destPath); err == nil {
				backupPath := destPath + ".deckbak"
				if dryRun {
					fmt.Printf("Would backup %s -> %s\n", destPath, backupPath)
				} else {
					os.Rename(destPath, backupPath)
					internal.LogVerbose(verbose, "Backed up %s -> %s", destPath, backupPath)
					os.Remove(destPath)
					internal.LogVerbose(verbose, "Removed %s", destPath)
				}
			}

			if dryRun {
				fmt.Printf("Would symlink %s -> %s\n", absSrc, destPath)
			} else {
				if err := os.Symlink(absSrc, destPath); err != nil {
					fmt.Printf("❌  %s → %s failed: %v\n", src, destPath, err)
				} else {
					fmt.Printf("✅  %s → %s\n", src, destPath)
				}
			}
		}
		return nil
	},
}

func expandPath(p string) string {
	if len(p) >= 2 && p[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, p[2:])
	}
	return p
}
