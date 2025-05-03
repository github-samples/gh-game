package cmd

import (
	"os"

	"github.com/chrisreddington/gh-game/internal/higherlower"
	userPrompt "github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/spf13/cobra"
)

var (
	minNumber int
	maxNumber int
)

var higherLowerCmd = &cobra.Command{
	Use:   "higherlower",
	Short: "Play Higher or Lower",
	Long: `Play a Higher or Lower number guessing game.
	
A number will be shown and you need to guess whether the
next number will be higher or lower than the current number.

The game continues until you make an incorrect guess. How long of a streak 
can you get?

Example usage:
  gh game higherlower
  gh game higherlower --min 1 --max 1000`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		input := userPrompt.New(os.Stdin, os.Stdout, os.Stderr)
		higherlower.PlayGame(input, minNumber, maxNumber)
	},
}

func init() {
	higherLowerCmd.Flags().IntVarP(&minNumber, "min", "m", 1, "Minimum possible number")
	higherLowerCmd.Flags().IntVarP(&maxNumber, "max", "M", 100, "Maximum possible number")

	rootCmd.AddCommand(higherLowerCmd)
}
