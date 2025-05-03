package cmd

import (
	"os"

	"github.com/chrisreddington/gh-game/internal/wordguess"
	userPrompt "github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/spf13/cobra"
)

var wordguessCmd = &cobra.Command{
	Use:   "wordguess",
	Short: "Play Word Guess",
	Long: `Start a game of Word Guess where you guess a GitHub-related term one letter at a time.
	
The rules are simple:
1. A random word will be selected
2. Guess one letter at a time
3. If the letter is in the word, it will be revealed
4. If not, you lose one of your available guesses
5. You win by guessing the word before running out of guesses
6. You lose if you make 6 incorrect guesses`,
	Run: func(cmd *cobra.Command, args []string) {
		input := userPrompt.New(os.Stdin, os.Stdout, os.Stderr)
		wordguess.PlayGame(input)
	},
}

func init() {
	rootCmd.AddCommand(wordguessCmd)
}
