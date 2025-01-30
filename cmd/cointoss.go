package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Helper functions
func tossCoin() string {
	rand.Seed(time.Now().UnixNano())
	if rand.Float32() < 0.5 {
		return "heads"
	}
	return "tails"
}

func validateGuess(guess string) error {
	guess = strings.ToLower(strings.TrimSpace(guess))
	if guess != "heads" && guess != "tails" {
		return fmt.Errorf("guess must be either 'heads' or 'tails'")
	}
	return nil
}

func getNextGuess() (string, bool) {
	fmt.Print("Play again? Enter 'heads' or 'tails' (or 'quit' to end): ")
	var answer string
	fmt.Scanln(&answer)

	if strings.ToLower(strings.TrimSpace(answer)) == "quit" {
		return "", false
	}

	if err := validateGuess(answer); err != nil {
		fmt.Println(err)
		return getNextGuess()
	}

	return strings.ToLower(strings.TrimSpace(answer)), true
}

var cointossCmd = &cobra.Command{
	Use:   "cointoss [guess]",
	Short: "Toss a coin",
	Long:  `Toss a virtual coin and get heads or tails as the result.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires exactly 1 argument (guess)")
		}
		return validateGuess(args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		guess := strings.ToLower(strings.TrimSpace(args[0]))
		streak := 0
		keepPlaying := true

		for keepPlaying {
			result := tossCoin()
			fmt.Printf("The coin shows: %s!\n", strings.Title(result))

			if guess == result {
				streak++
				fmt.Printf("Correct! Streak: %d\n", streak)
				var continuePlay bool
				guess, continuePlay = getNextGuess()
				keepPlaying = continuePlay
			} else {
				fmt.Printf("Game Over! Final streak: %d\n", streak)
				keepPlaying = false
			}
		}
	},
}
