package cmd

import (
	"fmt"
	"os"

	"github.com/chrisreddington/gh-game/internal/tictactoe"
	userPrompt "github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/spf13/cobra"
)

var tictactoeCmd = &cobra.Command{
	Use:   "tictactoe",
	Short: "Play Tic-tac-toe",
	Long:  `Start a game of Tic-tac-toe where X and O take turns to play.`,
	Run: func(cmd *cobra.Command, args []string) {
		game := tictactoe.NewGame()
		prompter := userPrompt.New(os.Stdin, os.Stdout, os.Stderr)

		for {
			fmt.Println(game)
			fmt.Printf("Player %s's turn\n", game.CurrentPlayer)

			row, col, err := tictactoe.GetPlayerMove(prompter, game)
			if err != nil {
				fmt.Printf("Error getting move: %v\n", err)
				return
			}

			if err := game.MakeMove(row, col); err != nil {
				fmt.Printf("Invalid move: %v\n", err)
				continue
			}

			if winner := game.GetWinner(); winner != "" {
				fmt.Println(game)
				fmt.Printf("Player %s wins!\n", winner)
				return
			}

			if game.IsBoardFull() {
				fmt.Println(game)
				fmt.Println("It's a draw!")
				return
			}
		}
	},
}
