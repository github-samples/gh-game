// This file is an entrypoint for the Rock Paper Scissors game. It uses the cobra cmd to execute the game.
// The game logic is in the internal/rockpaperscissors package.

// This file sets up the command line interface for the game. It should call the PlayGame function from the rockpaperscissors package.

package cmd

import (
	"os"

	"github.com/chrisreddington/gh-game/internal/rockpaperscissors"
	userPrompt "github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/spf13/cobra"
)

var secretMode bool

// rootCmd represents the base command when called without any subcommands
var rockPaperScissorsCmd = &cobra.Command{
	Use:   "rockpaperscissors",
	Short: "A simple Rock Paper Scissors game",
	Long: `A simple Rock Paper Scissors game that allows you to play against the computer.
You can choose from rock, paper, or scissors. The computer will randomly choose its move and the winner will be determined based on the rules of the game.`,
	Run: func(cmd *cobra.Command, args []string) {
		input := userPrompt.New(os.Stdin, os.Stdout, os.Stderr)
		rockpaperscissors.PlayGame(input, secretMode)
	},
}

func init() {
	rockPaperScissorsCmd.Flags().BoolVar(&secretMode, "spock", false, "Enable secret game mode")
	rootCmd.AddCommand(rockPaperScissorsCmd)
}
