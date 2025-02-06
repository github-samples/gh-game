package cmd

import (
	"fmt"
	"strings"

	"github.com/chrisreddington/gh-game/internal/cointoss" // adjust import path as needed

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
		guess := strings.ToLower(strings.TrimSpace(args[0]))
		streak := 0
		keepPlaying := true

		for keepPlaying {
			result := cointoss.TossCoin()
			fmt.Printf("The coin shows: %s!\n", strings.Title(result))

			if guess == result {
				streak++
				fmt.Printf("Correct! Streak: %d\n", streak)
				var continuePlay bool
				guess, continuePlay = cointoss.GetNextGuess()
				keepPlaying = continuePlay
			} else {
				fmt.Printf("Game Over! Final streak: %d\n", streak)
				keepPlaying = false
			}
		}
	},
}
