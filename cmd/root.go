package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh-game",
	Short: "A GitHub CLI extension for games",
	Long:  `A GitHub CLI extension that allows you to play games through the GitHub CLI.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(cointossCmd)
	rootCmd.AddCommand(tictactoeCmd)
}
