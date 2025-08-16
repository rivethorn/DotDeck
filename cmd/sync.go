package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/rivethorn/dotdeck/internal"
	"github.com/rivethorn/dotdeck/internal/config"
	"github.com/rivethorn/dotdeck/internal/runner"
	"github.com/spf13/cobra"
)

// Helper functions

func isGitDirty(path string) (bool, error) {
	out, err := exec.Command("git", "status", "--porcelain", path).Output()
	if err != nil {
		return false, err
	}
	return len(out) > 0, nil
}

func filesDiffer(a, b string) (bool, error) {
	f1, err := os.ReadFile(a)
	if err != nil {
		return false, err
	}
	f2, err := os.ReadFile(b)
	if err != nil {
		return false, err
	}
	return !bytes.Equal(f1, f2), nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func runCmd(name string, args ...string) error {
	internal.LogVerbose(verbose, "Running: %s %s", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runCmdOutput(name string, args ...string) (string, error) {
	internal.LogVerbose(verbose, "Running (capture): %s %s", name, strings.Join(args, " "))
	out, err := exec.Command(name, args...).Output()
	return strings.TrimSpace(string(out)), err
}

func hostname() string {
	h, _ := os.Hostname()
	return h
}

func stageAndCommit(files []string) error {
	args := append([]string{"add"}, files...)
	internal.LogVerbose(verbose, "Staging files: %v", files)
	if err := runCmd("git", args...); err != nil {
		return err
	}
	msg := fmt.Sprintf("Sync from %s at %s", hostname(), time.Now().Format("2006-01-02 15:04:05"))
	internal.LogVerbose(verbose, "Committing with message: %s", msg)
	return runCmd("git", "commit", "-m", msg)
}

func pullBeforePush() error {
	fmt.Println("🔄  Pulling latest changes from remote...")
	return runner.RunInteractive("git", "pull", "--ff-only")
}

func pushIfAhead() error {
	if err := runner.RunInteractive("git", "fetch"); err != nil {
		return err
	}
	status, err := runCmdOutput("git", "status", "-sb")
	if err != nil {
		return err
	}
	if !strings.Contains(status, "ahead") {
		fmt.Println("✓ Nothing to push — branch is up to date.")
		return nil
	}
	fmt.Println("⇪ Pushing changes to remote...")
	return runner.RunInteractive("git", "push")
}

// Commands

var force bool

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync dotfiles form system to repo, commit, merge, and push",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("config.toml")
		if err != nil {
			return err
		}

		changedFiles := []string{}

		for src, dest := range cfg.Files {
			absSrc, _ := filepath.Abs(src)
			destPath := expandPath(dest)

			if _, err := os.Stat(destPath); os.IsNotExist(err) {
				continue
			}

			dirty, err := isGitDirty(absSrc)
			if err != nil {
				return err
			}

			changed, err := filesDiffer(destPath, absSrc)
			if err != nil {
				return err
			}

			if changed && (dirty || force) {
				fmt.Printf("↩ syncing %s → %s\n", destPath, absSrc)
				if !dryRun {
					if err := copyFile(destPath, absSrc); err != nil {
						return err
					}
					changedFiles = append(changedFiles, absSrc)
				}
			} else {
				internal.LogVerbose(verbose, "Skipping %s - no changes or clean", src)
			}
		}

		if len(changedFiles) > 0 && !dryRun {
			if err := stageAndCommit(changedFiles); err != nil {
				return err
			}
			if err := pullBeforePush(); err != nil {
				return err
			}
			if err := pushIfAhead(); err != nil {
				return err
			}
		} else if len(changedFiles) == 0 {
			fmt.Println("✓ No changes to sync")
		}

		return nil
	},
}

func init() {
	syncCmd.Flags().BoolVarP(&force, "force", "f", false, "Force sync even if Git shows clean status")
	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate actions without making changes")
	rootCmd.AddCommand(syncCmd)
}
