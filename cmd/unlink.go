package cmd

import (
	"fmt"
	"os"

	"github.com/rivethorn/dotdeck/internal"
	"github.com/rivethorn/dotdeck/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unlinkCmd)
	unlinkCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate actions without making changes")
}

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Remove symlinks and restore backups",
	RunE: func(cmd *cobra.Command, args []string) error {
		internal.LogVerbose(verbose, "Starting unlink command")
		internal.LogVerbose(verbose, "Checking for config file")
		cfg, err := config.Load("config.toml")
		if err != nil {
			return err
		}
		internal.LogVerbose(verbose, "Config file loaded successfully")

		internal.LogVerbose(verbose, "Evaluating files")
		for _, dest := range cfg.Files {
			destPath := expandPath(dest)
			backupPath := destPath + ".deckbak"

			internal.LogVerbose(verbose, "Evaluating %s", destPath)
			info, err := os.Lstat(destPath)
			if err != nil {
				fmt.Printf(" %s missing, skipping...\n", destPath)
				continue
			}

			if info.Mode()&os.ModeSymlink != 0 {
				if dryRun {
					fmt.Printf("Would remove symlink %s\n", destPath)
				} else {
					// Check if backup exists
					if _, err := os.Stat(backupPath); err != nil {
						// No backup, ask for confirmation
						fmt.Printf("No backup found for %s. You'll lose your config. Delete symlink? [y/N]: ", destPath)
						var response string
						_, err := fmt.Scanln(&response)
						if err != nil {
							return err
						}
						if response != "y" && response != "Y" {
							fmt.Printf("Skipping %s\n", destPath)
							continue
						}
						err = os.Remove(destPath)
						if err != nil {
							return err
						}
						internal.LogVerbose(verbose, "Removed %s symlink", destPath)
					}
					return nil
				}
				if dryRun {
					fmt.Printf("Would restore backup %s -> %s\n", backupPath, destPath)
				} else {
					if _, err := os.Stat(backupPath); err == nil {
						if err := os.Rename(backupPath, destPath); err != nil {
							fmt.Printf("❌ Restoring backup for %s failed: %v\n", destPath, err)
							continue
						}
						fmt.Printf("󰑌 Restored backup for %s\n", destPath)
					} else {
						fmt.Printf("󰩹 Removed symlink %s (no backup found)\n", destPath)
					}
					return nil
				}
			} else {
				fmt.Printf(" %s is not a symlink, skipping\n", destPath)
				return nil
			}
		}
		return nil
	},
}
