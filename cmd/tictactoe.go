package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/chrisreddington/gh-game/internal/tictactoe"
	userPrompt "github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/spf13/cobra"
)

var (
	xStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))  // bright blue
	oStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // bright yellow
)

var tictactoeCmd = &cobra.Command{
	Use:   "tictactoe",
	Short: "Play Tic-tac-toe",
	Long: `Start a game of Tic-tac-toe where you can play against another player locally or against the computer.
Choose between two game modes:
- Local Multiplayer: Play against another player on the same computer
- Play Against Computer: Play against an AI opponent that uses basic strategy`,
	Run: func(cmd *cobra.Command, args []string) {
		prompter := userPrompt.New(os.Stdin, os.Stdout, os.Stderr)

		// Select game mode
		modeIndex, err := prompter.Select(
			"Select game mode:",
			"Local Multiplayer",
			[]string{"Local Multiplayer", "Play Against Computer"},
		)
		if err != nil {
			fmt.Printf("Error selecting game mode: %v\n", err)
			return
		}

		// Initialize game with selected mode
		mode := tictactoe.LocalGame
		if modeIndex == 1 {
			mode = tictactoe.ComputerGame
		}
		game := tictactoe.NewGame(mode)

		// Main game loop
		for {
			fmt.Println(game)
			currentMark := game.CurrentPlayer
			style := xStyle
			if currentMark == "O" {
				style = oStyle
			}
			fmt.Printf("Player %s's turn\n", style.Render(currentMark))

			// Get move from either computer or human player
			var rowIndex, columnIndex int
			if game.IsComputerTurn() {
				rowIndex, columnIndex = game.GetComputerMove()
				position := rowIndex*3 + columnIndex + 1
				fmt.Printf("Computer places %s at position %d\n", oStyle.Render("O"), position)
			} else {
				var err error
				rowIndex, columnIndex, err = tictactoe.GetPlayerMove(prompter, game)
				if err != nil {
					fmt.Printf("Error getting move: %v\n", err)
					return
				}
			}

			// Apply the move
			if err := game.MakeMove(rowIndex, columnIndex); err != nil {
				fmt.Printf("Invalid move: %v\n", err)
				continue
			}

			// Check win condition
			if winner := game.GetWinner(); winner != "" {
				fmt.Println(game)
				style := xStyle
				if winner == "O" {
					style = oStyle
				}
				if game.Mode == tictactoe.ComputerGame && winner == game.ComputerMark {
					fmt.Printf("Computer (%s) wins!\n", style.Render(winner))
				} else {
					fmt.Printf("Player %s wins!\n", style.Render(winner))
				}
				return
			}

			// Check draw condition
			if game.IsBoardFull() {
				fmt.Println(game)
				fmt.Println("It's a draw!")
				return
			}
		}
	},
}
