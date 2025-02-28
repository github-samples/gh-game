package cmd

import (
	"fmt"
	"os"

	"github.com/chrisreddington/gh-game/internal/cointoss"
	userPrompt "github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/spf13/cobra"
)

var cointossCmd = &cobra.Command{
	Use:   "cointoss [guess]",
	Short: "Toss a coin",
	Long:  `Toss a virtual coin and get heads or tails as the result.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires exactly 1 argument (guess)")
		}
		return cointoss.ValidateGuess(args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		input := userPrompt.New(os.Stdin, os.Stdout, os.Stderr)
		cointoss.PlayGame(input, args[0])
	},
}

func init() {
	rootCmd.AddCommand(cointossCmd)
}
