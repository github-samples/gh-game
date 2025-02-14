// Package tictactoe implements a classic Tic-tac-toe game where two players
// take turns marking spaces on a 3x3 grid. Players alternate placing X's and O's
// on the board until either one player wins by placing three marks in a row,
// column, or diagonal, or the game ends in a draw when the board is full.
package tictactoe

import (
	"errors"
	"fmt"
)

// Prompter defines an interface for getting user input
type Prompter interface {
	Select(prompt, defaultValue string, options []string) (int, error)
}

// GameInterface defines the methods needed for GetPlayerMove
type GameInterface interface {
	GetAvailablePositions() []string
}

// Board represents a 3x3 Tic-tac-toe game board.
// Empty squares are represented by empty strings,
// and played squares contain either "X" or "O".
type Board [3][3]string

// Game represents the current state of a Tic-tac-toe game.
type Game struct {
	board         Board  // The game board storing player moves ("X" or "O", or empty for unplayed)
	CurrentPlayer string // CurrentPlayer indicates whose turn it is ("X" or "O")
}

// NewGame creates and initializes a new Tic-tac-toe game with an empty board.
// Setting as X always plays first.
func NewGame() *Game {
	return &Game{
		board:         Board{},
		CurrentPlayer: "X",
	}
}

// MakeMove attempts to place the current player's mark at the specified position.
// The position is specified using zero-based indices for row and column.
// Returns an error if:
// - The position is out of bounds (not between 0 and 2)
// - The position is already occupied by a player's mark
func (g *Game) MakeMove(rowIndex, columnIndex int) error {
	if rowIndex < 0 || rowIndex > 2 || columnIndex < 0 || columnIndex > 2 {
		return errors.New("invalid position: must be between 0 and 2")
	}
	if g.board[rowIndex][columnIndex] != "" {
		return errors.New("position already taken")
	}
	g.board[rowIndex][columnIndex] = g.CurrentPlayer
	g.CurrentPlayer = switchPlayer(g.CurrentPlayer)
	return nil
}

// GetWinner checks if there is a winner by examining all possible winning combinations:
// - Three rows
// - Three columns
// - Two diagonals
// Returns the winning player's mark ("X" or "O"), or an empty string if there's no winner.
func (g *Game) GetWinner() string {
	// Check rows for a winner
	for rowIndex := 0; rowIndex < 3; rowIndex++ {
		if g.board[rowIndex][0] != "" &&
			g.board[rowIndex][0] == g.board[rowIndex][1] &&
			g.board[rowIndex][1] == g.board[rowIndex][2] {
			return g.board[rowIndex][0]
		}
	}

	// Check columns for a winner
	for columnIndex := 0; columnIndex < 3; columnIndex++ {
		if g.board[0][columnIndex] != "" &&
			g.board[0][columnIndex] == g.board[1][columnIndex] &&
			g.board[1][columnIndex] == g.board[2][columnIndex] {
			return g.board[0][columnIndex]
		}
	}

	// Check main diagonal (top-left to bottom-right)
	if g.board[0][0] != "" &&
		g.board[0][0] == g.board[1][1] &&
		g.board[1][1] == g.board[2][2] {
		return g.board[0][0]
	}

	// Check secondary diagonal (top-right to bottom-left)
	if g.board[0][2] != "" &&
		g.board[0][2] == g.board[1][1] &&
		g.board[1][1] == g.board[2][0] {
		return g.board[0][2]
	}

	return ""
}

// IsBoardFull determines if all positions on the board have been played.
// Returns true if no empty positions remain, false otherwise.
func (g *Game) IsBoardFull() bool {

	for rowIndex := 0; rowIndex < 3; rowIndex++ {
		for columnIndex := 0; columnIndex < 3; columnIndex++ {
			if g.board[rowIndex][columnIndex] == "" {
				return false
			}
		}
	}
	return true
}

// String returns a formatted string representation of the current board state.
// Empty squares are shown as numbers 1-9 for position selection, making it
// easier for players to choose their moves. Played squares show the player's mark.
// The board is displayed with column separators (|) and row separators (-).
func (g *Game) String() string {
	result := "\n"
	position := 1
	for rowIndex := 0; rowIndex < 3; rowIndex++ {
		for columnIndex := 0; columnIndex < 3; columnIndex++ {
			// Add leading space for first column
			if columnIndex == 0 {
				result += " "
			}

			// Show position number for empty squares, otherwise show the player's mark
			if g.board[rowIndex][columnIndex] == "" {
				result += fmt.Sprintf("%d", position)
			} else {
				result += g.board[rowIndex][columnIndex]
			}

			// Add column separators except for the last column
			if columnIndex < 2 {
				result += " | "
			}
			position++
		}
		result += "\n"

		// Add row separators except for the last row
		if rowIndex < 2 {
			result += "---+---+---\n"
		}
	}
	return result
}

// switchPlayer determines the next player's turn.
// Following traditional Tic-tac-toe rules, players alternate between X and O.
func switchPlayer(currentPlayer string) string {
	if currentPlayer == "X" {
		return "O"
	}
	return "X"
}

// positionToRowCol converts a one-based position (1-9) to zero-based row and column indices.
// The board positions are mapped as follows:
// 1 2 3
// 4 5 6
// 7 8 9
// Returns (-1, -1) for invalid positions (0 or >9)
func positionToRowCol(position int) (rowIndex, columnIndex int) {
	if position <= 0 || position > 9 {
		return -1, -1
	}
	position-- // Convert to 0-based index
	return position / 3, position % 3
}

// GetAvailablePositions returns a slice of strings representing unoccupied positions.
// The positions are numbered 1-9 (one-based) to match the display format.
func (g *Game) GetAvailablePositions() []string {
	var availablePositions []string
	for position := 1; position <= 9; position++ {
		rowIndex, columnIndex := positionToRowCol(position)
		if g.board[rowIndex][columnIndex] == "" {
			availablePositions = append(availablePositions, fmt.Sprintf("%d", position))
		}
	}
	return availablePositions
}

// GetPlayerMove prompts the user to select a valid move and returns the chosen
// row and column indices. It uses the provided Prompter interface to get user input.
// Returns an error if:
// - No valid moves are available (board is full)
// - User input is invalid
// - Selected position is invalid
func GetPlayerMove(prompter Prompter, game GameInterface) (rowIndex, columnIndex int, err error) {
	availablePositions := game.GetAvailablePositions()
	if len(availablePositions) == 0 {
		return -1, -1, errors.New("no available moves")
	}

	posIndex, err := prompter.Select("Select position (1-9):", "1", availablePositions)
	if err != nil {
		return -1, -1, err
	}

	if posIndex < 0 || posIndex >= len(availablePositions) {
		return -1, -1, fmt.Errorf("invalid position selection: %d", posIndex)
	}

	var position int
	_, err = fmt.Sscanf(availablePositions[posIndex], "%d", &position)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid position value: %v", err)
	}

	rowIndex, columnIndex = positionToRowCol(position)
	if rowIndex < 0 || columnIndex < 0 {
		return -1, -1, fmt.Errorf("invalid position value: %d", position)
	}

	return rowIndex, columnIndex, nil
}
