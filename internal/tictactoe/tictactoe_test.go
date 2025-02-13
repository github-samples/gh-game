package tictactoe

import (
	"fmt"
	"strings"
	"testing"
)

type mockPrompter struct {
	selectAnswers []int
	selectIndex   int
	selectError   error
}

func (m *mockPrompter) Select(prompt string, defaultValue string, options []string) (int, error) {
	if m.selectError != nil {
		return 0, m.selectError
	}
	if len(m.selectAnswers) == 0 {
		return 0, fmt.Errorf("no answers configured")
	}
	answer := m.selectAnswers[m.selectIndex]
	m.selectIndex++
	return answer, nil
}

// mockGame implements GameInterface for testing
type mockGame struct {
	positions []string
}

func (m *mockGame) GetAvailablePositions() []string {
	return m.positions
}

// Helper functions
func setupGameWithMoves(moves [][2]int) *Game {
	game := NewGame()
	for _, move := range moves {
		game.MakeMove(move[0], move[1])
	}
	return game
}

func createFullBoard() *Game {
	moves := [][2]int{
		{0, 0}, {0, 1}, {0, 2},
		{1, 0}, {1, 1}, {1, 2},
		{2, 0}, {2, 1}, {2, 2},
	}
	return setupGameWithMoves(moves)
}

type gameTestCase struct {
	name          string
	moves         [][2]int
	expectedBoard [][]string
	expectedError string
	currentPlayer string
}

func TestNewGame(t *testing.T) {
	tests := []struct {
		name               string
		expectedPlayer     string
		expectedBoardEmpty bool
	}{
		{
			name:               "new game starts with player X",
			expectedPlayer:     "X",
			expectedBoardEmpty: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame()
			if game.CurrentPlayer != tt.expectedPlayer {
				t.Errorf("Expected first player to be %s, got %s", tt.expectedPlayer, game.CurrentPlayer)
			}
			if tt.expectedBoardEmpty {
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						if game.board[i][j] != "" {
							t.Errorf("Expected empty board at position (%d,%d), got %s", i, j, game.board[i][j])
						}
					}
				}
			}
		})
	}
}

func TestMakeMove(t *testing.T) {
	tests := []gameTestCase{
		{
			name:  "valid move succeeds",
			moves: [][2]int{{0, 0}},
			expectedBoard: [][]string{
				{"X", "", ""},
				{"", "", ""},
				{"", "", ""},
			},
			currentPlayer: "O",
		},
		{
			name:          "move outside board fails",
			moves:         [][2]int{{3, 3}},
			expectedError: "invalid position",
		},
		{
			name: "move on occupied position fails",
			moves: [][2]int{
				{0, 0},
				{0, 0},
			},
			expectedError: "position already taken",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame()
			var lastErr error

			for _, move := range tt.moves {
				lastErr = game.MakeMove(move[0], move[1])
			}

			if tt.expectedError != "" {
				if lastErr == nil || !strings.Contains(lastErr.Error(), tt.expectedError) {
					t.Errorf("Expected error containing %q, got %v", tt.expectedError, lastErr)
				}
				return
			}

			if lastErr != nil {
				t.Errorf("Unexpected error: %v", lastErr)
			}

			if tt.expectedBoard != nil {
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						if game.board[i][j] != tt.expectedBoard[i][j] {
							t.Errorf("Position (%d,%d): expected %q, got %q",
								i, j, tt.expectedBoard[i][j], game.board[i][j])
						}
					}
				}
			}

			if tt.currentPlayer != "" && game.CurrentPlayer != tt.currentPlayer {
				t.Errorf("Expected current player to be %s, got %s", tt.currentPlayer, game.CurrentPlayer)
			}
		})
	}
}

func TestGetWinner(t *testing.T) {
	tests := []struct {
		name     string
		moves    [][2]int
		expected string
	}{
		{
			name: "row win",
			moves: [][2]int{
				{0, 0}, {1, 0},
				{0, 1}, {1, 1},
				{0, 2},
			},
			expected: "X",
		},
		{
			name: "column win",
			moves: [][2]int{
				{0, 0}, {0, 1},
				{1, 0}, {1, 1},
				{2, 0},
			},
			expected: "X",
		},
		{
			name: "diagonal win",
			moves: [][2]int{
				{0, 0}, {0, 1},
				{1, 1}, {1, 0},
				{2, 2},
			},
			expected: "X",
		},
		{
			name: "secondary diagonal win",
			moves: [][2]int{
				{0, 2}, {0, 0},
				{1, 1}, {1, 0},
				{2, 0},
			},
			expected: "X",
		},
		{
			name: "no winner",
			moves: [][2]int{
				{0, 0}, {0, 1},
				{1, 1},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame()
			for _, move := range tt.moves {
				game.MakeMove(move[0], move[1])
			}
			if winner := game.GetWinner(); winner != tt.expected {
				t.Errorf("Expected winner %s, got %s", tt.expected, winner)
			}
		})
	}
}

func TestIsBoardFull(t *testing.T) {
	tests := []struct {
		name         string
		moves        [][2]int
		expectedFull bool
	}{
		{
			name:         "empty board is not full",
			moves:        [][2]int{},
			expectedFull: false,
		},
		{
			name: "partially filled board is not full",
			moves: [][2]int{
				{0, 0}, {0, 1}, {0, 2},
				{1, 0}, {1, 1},
			},
			expectedFull: false,
		},
		{
			name: "completely filled board is full",
			moves: [][2]int{
				{0, 0}, {1, 0}, // X at (0,0), O at (1,0)
				{0, 1}, {1, 1}, // X at (0,1), O at (1,1)
				{0, 2}, {1, 2}, // X at (0,2), O at (1,2)
				{2, 0}, {2, 1}, // X at (2,0), O at (2,1)
				{2, 2}, // X at (2,2)
			},
			expectedFull: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame()
			for _, move := range tt.moves {
				err := game.MakeMove(move[0], move[1])
				if err != nil {
					t.Fatalf("Failed to make move %v: %v", move, err)
				}
			}

			actual := game.IsBoardFull()
			if actual != tt.expectedFull {
				t.Errorf("Expected board full to be %v, got %v\nBoard state:\n%s",
					tt.expectedFull, actual, game.String())
			}
		})
	}
}

func TestGetPlayerMove(t *testing.T) {
	tests := []struct {
		name          string
		boardState    [][2]int
		selectAnswers []int
		selectError   error
		expectedMove  [2]int
		expectedError string
	}{
		{
			name:          "valid center move",
			selectAnswers: []int{4},
			expectedMove:  [2]int{1, 1},
		},
		{
			name:          "valid move after center taken",
			boardState:    [][2]int{{1, 1}},
			selectAnswers: []int{0},
			expectedMove:  [2]int{0, 0},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame()
			for _, move := range tt.boardState {
				game.MakeMove(move[0], move[1])
			}
			mockP := &mockPrompter{
				selectAnswers: tt.selectAnswers,
				selectError:   tt.selectError,
			}
			row, col, err := GetPlayerMove(mockP, game)
			if tt.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing %q, got %v", tt.expectedError, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if row != tt.expectedMove[0] || col != tt.expectedMove[1] {
				t.Errorf("Expected move (%d,%d), got (%d,%d)",
					tt.expectedMove[0], tt.expectedMove[1], row, col)
			}
		})
	}
}

func TestGetPlayerMoveErrors(t *testing.T) {
	tests := []struct {
		name          string
		game          GameInterface
		selectAnswers []int
		selectError   error
		expectedError string
	}{
		{
			name:          "full board returns error",
			game:          createFullBoard(),
			selectAnswers: []int{0},
			expectedError: "no available moves",
		},
		{
			name:          "invalid position selection",
			game:          NewGame(),
			selectAnswers: []int{9},
			expectedError: "invalid position",
		},
		{
			name:          "prompter error is propagated",
			game:          NewGame(),
			selectError:   fmt.Errorf("prompter error"),
			expectedError: "prompter error",
		},
		{
			name:          "invalid position value parsing",
			game:          &mockGame{positions: []string{"invalid"}},
			selectAnswers: []int{0},
			expectedError: "invalid position value",
		},
		{
			name:          "negative position selection",
			game:          NewGame(),
			selectAnswers: []int{-1},
			expectedError: "invalid position",
		},
		{
			name:          "empty position list",
			game:          &mockGame{positions: []string{}},
			selectAnswers: []int{0},
			expectedError: "no available moves",
		},
		{
			name:          "position out of valid range",
			game:          &mockGame{positions: []string{"10"}}, // Position 10 is invalid
			selectAnswers: []int{0},
			expectedError: "invalid position value: 10",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockP := &mockPrompter{
				selectAnswers: tt.selectAnswers,
				selectError:   tt.selectError,
			}

			_, _, err := GetPlayerMove(mockP, tt.game)
			if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing %q, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestPositionToRowCol(t *testing.T) {
	tests := []struct {
		name     string
		position int
		row      int
		col      int
	}{
		{
			name:     "top left corner",
			position: 1,
			row:      0,
			col:      0,
		},
		{
			name:     "top middle position",
			position: 2,
			row:      0,
			col:      1,
		},
		{
			name:     "top right corner",
			position: 3,
			row:      0,
			col:      2,
		},
		{
			name:     "center left position",
			position: 4,
			row:      1,
			col:      0,
		},
		{
			name:     "center middle position",
			position: 5,
			row:      1,
			col:      1,
		},
		{
			name:     "center right position",
			position: 6,
			row:      1,
			col:      2,
		},
		{
			name:     "bottom left corner",
			position: 7,
			row:      2,
			col:      0,
		},
		{
			name:     "bottom middle position",
			position: 8,
			row:      2,
			col:      1,
		},
		{
			name:     "bottom right corner",
			position: 9,
			row:      2,
			col:      2,
		},
		{
			name:     "zero position",
			position: 0,
			row:      -1,
			col:      -1,
		},
		{
			name:     "beyond board position",
			position: 10,
			row:      -1,
			col:      -1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			row, col := positionToRowCol(tt.position)
			if row != tt.row || col != tt.col {
				t.Errorf("positionToRowCol(%d) = (%d,%d); want (%d,%d)",
					tt.position, row, col, tt.row, tt.col)
			}
		})
	}
}

func TestGetAvailablePositions(t *testing.T) {
	tests := []struct {
		name             string
		moves            [][2]int
		expectedCount    int
		excludePositions []string
	}{
		{
			name:          "empty board has all positions",
			moves:         [][2]int{},
			expectedCount: 9,
		},
		{
			name: "three moves taken",
			moves: [][2]int{
				{0, 0}, // Position 1
				{1, 1}, // Position 5
				{2, 2}, // Position 9
			},
			expectedCount:    6,
			excludePositions: []string{"1", "5", "9"},
		},
		{
			name: "corners taken",
			moves: [][2]int{
				{0, 0}, // Position 1
				{0, 2}, // Position 3
				{2, 0}, // Position 7
				{2, 2}, // Position 9
			},
			expectedCount:    5,
			excludePositions: []string{"1", "3", "7", "9"},
		},
		{
			name: "all positions taken",
			moves: [][2]int{
				{0, 0}, {0, 1}, {0, 2},
				{1, 0}, {1, 1}, {1, 2},
				{2, 0}, {2, 1}, {2, 2},
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			game := setupGameWithMoves(tt.moves)
			positions := game.GetAvailablePositions()

			if len(positions) != tt.expectedCount {
				t.Errorf("Expected %d available positions, got %d", tt.expectedCount, len(positions))
			}

			for _, excluded := range tt.excludePositions {
				for _, pos := range positions {
					if pos == excluded {
						t.Errorf("Position %s should not be available", excluded)
					}
				}
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name           string
		moves          [][2]int
		expectedPieces []string
	}{
		{
			name:           "empty board",
			moves:          [][2]int{},
			expectedPieces: []string{" ", " ", " ", " ", " ", " ", " ", " ", " "},
		},
		{
			name:           "X in center",
			moves:          [][2]int{{1, 1}},
			expectedPieces: []string{" ", " ", " ", " ", "X", " ", " ", " ", " "},
		},
		{
			name: "multiple moves",
			moves: [][2]int{
				{0, 0}, // X top-left
				{1, 1}, // O center
			},
			expectedPieces: []string{"X", " ", " ", " ", "O", " ", " ", " ", " "},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			game := setupGameWithMoves(tt.moves)
			board := game.String()

			for i, piece := range tt.expectedPieces {
				if piece != " " && !strings.Contains(board, piece) {
					t.Errorf("Expected piece %s at position %d in board string", piece, i+1)
				}
			}
		})
	}
}

func TestSwitchPlayer(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		expected string
	}{
		{
			name:     "X to O",
			current:  "X",
			expected: "O",
		},
		{
			name:     "O to X",
			current:  "O",
			expected: "X",
		},
		{
			name:     "invalid player",
			current:  "invalid",
			expected: "X", // default to X for invalid input
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if next := switchPlayer(tt.current); next != tt.expected {
				t.Errorf("Expected %s after %s, got %s", tt.expected, tt.current, next)
			}
		})
	}
}
