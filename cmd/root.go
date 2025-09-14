// Package cmd contains the functionalities that are needed for the app
// to run.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	dryRun  bool
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "deck",
	Short: "DotDeck - simple, clear dotfile linking",
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"Show detailed output of every step",
	)
}

// Execute the Cobra rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
