package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rivethorn/dotdeck/internal"
	"github.com/rivethorn/dotdeck/internal/config"
	"github.com/rivethorn/dotdeck/internal/runner"
	"github.com/spf13/cobra"
)

var (
	pulling bool
	pushing bool
	force   bool
)

// Helper functions

func isGitInstalled() (bool, error) {
	_, err := exec.Command("git", "--version").Output()
	if err != nil {
		return false, err
	}
	return true, nil
}

func isGitRepo() (bool, error) {
	_, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return false, err
	}
	return true, nil
}

func isGitDirty() (bool, error) {
	out, err := exec.Command("git", "diff-files").Output()
	if err != nil {
		return false, err
	}
	return len(out) > 0, nil
}

func runCmd(name string, args ...string) error {
	internal.LogVerbose(verbose, "Running: %s %s", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
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
	fmt.Println("󰓂 Pulling latest changes from remote...")
	return runner.RunInteractive("git", "pull", "--ff-only")
}

// Commands

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync dotfiles form system to repo, commit, pull, and push",
	RunE: func(cmd *cobra.Command, args []string) error {
		internal.LogVerbose(verbose, "Starting sync command")
		internal.LogVerbose(verbose, "Checking for config file")
		_, err := config.Load("config.toml")
		if err != nil {
			fmt.Println("You're not in a valid DotDeck repository")
			return err
		}
		internal.LogVerbose(verbose, "Config loaded successfully")

		internal.LogVerbose(verbose, "Checking if Git is installed")
		gitInstalled, err := isGitInstalled()
		if err != nil {
			return err
		}
		if !gitInstalled {
			return fmt.Errorf("Git is not installed, please install Git and try again")
		}

		internal.LogVerbose(verbose, "Checking if inside a Git repository")
		gitRepo, err := isGitRepo()
		if err != nil {
			return err
		}
		if !gitRepo {
			return fmt.Errorf("Not inside a Git repository")
		}

		internal.LogVerbose(verbose, "Checking for changes to sync")
		dirty, err := isGitDirty()
		if err != nil {
			return err
		}
		if dirty {
			internal.LogVerbose(verbose, "Local changes detected")
		}

		if pulling {
			if dirty {
				fmt.Print("You already have local changes, are you sure you want to continue? (y/N) ")
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					return fmt.Errorf("Aborted")
				}
				fmt.Println(" Force pulling changes...")
				runner.RunInteractive("git", "reset", "--hard", "origin")
				return nil
			}
			if err := pullBeforePush(); err != nil {
				return err
			}
			return nil
		}

		if pushing {
			if !dirty {
				if force {
					fmt.Println(" Force pushing changes...")
					runner.RunInteractive("git", "push", "--force")
					return nil
				}
				fmt.Println(" Nothing to push — branch is up to date.")
				return nil

			}
			if err := stageAndCommit([]string{"."}); err != nil {
				return err
			}
			fmt.Println("⇪ Pushing changes to remote...")
			return runner.RunInteractive("git", "push")
		}

		fmt.Println("Usage: dotdeck sync [options]")
		fmt.Println("Options:")
		fmt.Println("  -d, --pull       Pull from the remote repo")
		fmt.Println("  -u, --push       Push to the remote repo")
		fmt.Println("  -f, --force      Force sync even if Git shows clean status")
		fmt.Println("      --dry-run    Simulate actions without making changes")

		return nil
	},
}

func init() {
	syncCmd.Flags().BoolVarP(&pulling, "pull", "d", false, "Pull from the remote repo")
	syncCmd.Flags().BoolVarP(&pushing, "push", "u", false, "Push to the remote repo")
	syncCmd.Flags().BoolVarP(&force, "force", "f", false, "Force sync even if Git shows clean status")
	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate actions without making changes")
	rootCmd.AddCommand(syncCmd)
}
