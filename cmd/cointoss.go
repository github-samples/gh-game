package cmd

import (
	"fmt"
	"os"
	"strings"

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
		game := cointoss.NewGame()
		prompter := userPrompt.New(os.Stdin, os.Stdout, os.Stderr)
		guess := strings.ToLower(strings.TrimSpace(args[0]))
		streak := 0
		keepPlaying := true

		for keepPlaying {
			game.Play(guess)
			fmt.Println(game.GetResult())

			if game.PlayerGuess == game.Result {
				streak++
				fmt.Printf("Streak: %d\n", streak)
				guess, keepPlaying = cointoss.GetPlayerGuess(prompter)
			} else {
				fmt.Printf("Game Over! Final streak: %d\n", streak)
				keepPlaying = false
			}
		}
	},
}
