package tictactoe

import (
	"fmt"
	"strings"
	"testing"
)

// mockPrompter implements the Prompter interface for tictactoe game testing.
// It provides predefined select responses and can be configured to return errors.
type mockPrompter struct {
	selectAnswers []int // Predefined responses for Select calls
	selectIndex   int   // Current index in selectAnswers
	selectError   error // Error to be returned by Select
}

// Select implements the Prompter interface by returning either the configured error
// or a predefined answer from selectAnswers. It advances selectIndex after each call.
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

// mockGame implements GameInterface for testing purposes.
// It allows controlling the available positions returned for tests.
type mockGame struct {
	positions []string // The positions to be returned by GetAvailablePositions
}

// GetAvailablePositions implements GameInterface by returning the
// preconfigured positions for testing.
func (m *mockGame) GetAvailablePositions() []string {
	return m.positions
}

// Helper functions

// setupGameWithMoves creates a new game and plays a sequence of moves on it.
// This is used to setup test scenarios with specific board states.
// It fails the test if any of the moves are invalid.
func setupGameWithMoves(moves [][2]int, t *testing.T) *Game {
	game := NewGame(LocalGame)
	for _, move := range moves {
		if err := game.MakeMove(move[0], move[1]); err != nil {
			// In test setup we expect all moves to be valid
			t.Fatalf("setupGameWithMoves: failed to make move %v: %v", move, err)
		}
	}
	return game
}

// createFullBoard creates a game with a completely filled board.
// This is useful for testing tie scenarios and board validation.
func createFullBoard(t *testing.T) *Game {
	moves := [][2]int{
		{0, 0}, {0, 1}, {0, 2},
		{1, 0}, {1, 1}, {1, 2},
		{2, 0}, {2, 1}, {2, 2},
	}
	return setupGameWithMoves(moves, t)
}

// gameTestCase defines a test case structure for game-related tests.
// It contains the test name, a sequence of moves to execute, the expected
// resulting board state, any expected error, and the expected current player.
type gameTestCase struct {
	name          string     // Name of the test case
	moves         [][2]int   // Sequence of row,column moves to make
	expectedBoard [][]string // Expected board state after moves
	expectedError string     // Expected error message, empty if no error expected
	currentPlayer string     // Expected current player after moves
}

func TestNewGame(t *testing.T) {
	tests := []struct {
		name           string
		mode           GameMode
		wantPlayer     string
		wantMode       GameMode
		wantCompMark   string
		expectComputer bool
	}{
		{
			name:           "local multiplayer game",
			mode:           LocalGame,
			wantPlayer:     "X",
			wantMode:       LocalGame,
			wantCompMark:   "",
			expectComputer: false,
		},
		{
			name:           "computer game",
			mode:           ComputerGame,
			wantPlayer:     "X",
			wantMode:       ComputerGame,
			wantCompMark:   "O",
			expectComputer: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(tt.mode)

			if game.CurrentPlayer != tt.wantPlayer {
				t.Errorf("NewGame() CurrentPlayer = %v, want %v", game.CurrentPlayer, tt.wantPlayer)
			}

			if game.Mode != tt.wantMode {
				t.Errorf("NewGame() Mode = %v, want %v", game.Mode, tt.wantMode)
			}

			if game.ComputerMark != tt.wantCompMark {
				t.Errorf("NewGame() ComputerMark = %v, want %v", game.ComputerMark, tt.wantCompMark)
			}

			// Verify board is empty
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if game.board[i][j] != "" {
						t.Errorf("Expected empty board at position (%d,%d), got %s", i, j, game.board[i][j])
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
			game := NewGame(LocalGame)
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
			game := NewGame(LocalGame)
			for _, move := range tt.moves {
				if err := game.MakeMove(move[0], move[1]); err != nil {
					t.Fatalf("Failed to make move %v: %v", move, err)
				}
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
			game := NewGame(LocalGame)
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
			game := NewGame(LocalGame)
			for _, move := range tt.boardState {
				if err := game.MakeMove(move[0], move[1]); err != nil {
					t.Fatalf("Failed to make move %v: %v", move, err)
				}
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
			game:          createFullBoard(t),
			selectAnswers: []int{0},
			expectedError: "no available moves",
		},
		{
			name:          "invalid position selection",
			game:          NewGame(LocalGame),
			selectAnswers: []int{9},
			expectedError: "invalid position",
		},
		{
			name:          "prompter error is propagated",
			game:          NewGame(LocalGame),
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
			game:          NewGame(LocalGame),
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
			game := setupGameWithMoves(tt.moves, t)
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
		name              string
		moves             [][2]int
		expectedPositions []string
		expectedMarks     map[string]bool // Check for presence of colored X and O
	}{
		{
			name:              "empty board shows all positions",
			moves:             [][2]int{},
			expectedPositions: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
			expectedMarks:     map[string]bool{},
		},
		{
			name:              "X in center shows colored X",
			moves:             [][2]int{{1, 1}},
			expectedPositions: []string{"1", "2", "3", "4", "6", "7", "8", "9"},
			expectedMarks: map[string]bool{
				xStyle.Render("X"): true,
			},
		},
		{
			name: "multiple moves show colored X and O",
			moves: [][2]int{
				{0, 0}, // X top-left
				{1, 1}, // O center
			},
			expectedPositions: []string{"2", "3", "4", "6", "7", "8", "9"},
			expectedMarks: map[string]bool{
				xStyle.Render("X"): true,
				oStyle.Render("O"): true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := setupGameWithMoves(tt.moves, t)
			board := game.String()

			// Check that expected position numbers are present
			for _, pos := range tt.expectedPositions {
				if !strings.Contains(board, pos) {
					t.Errorf("Expected position %s to be present in board", pos)
				}
			}

			// Check that expected colored marks are present
			for mark := range tt.expectedMarks {
				if !strings.Contains(board, mark) {
					t.Errorf("Expected colored mark %q to be present in board", mark)
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

func TestIsComputerTurn(t *testing.T) {
	tests := []struct {
		name          string
		mode          GameMode
		currentPlayer string
		computerMark  string
		want          bool
	}{
		{
			name:          "computer's turn in computer game",
			mode:          ComputerGame,
			currentPlayer: "O",
			computerMark:  "O",
			want:          true,
		},
		{
			name:          "player's turn in computer game",
			mode:          ComputerGame,
			currentPlayer: "X",
			computerMark:  "O",
			want:          false,
		},
		{
			name:          "O's turn in local game",
			mode:          LocalGame,
			currentPlayer: "O",
			computerMark:  "",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &Game{
				Mode:          tt.mode,
				CurrentPlayer: tt.currentPlayer,
				ComputerMark:  tt.computerMark,
			}
			if got := game.IsComputerTurn(); got != tt.want {
				t.Errorf("IsComputerTurn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetComputerMove(t *testing.T) {
	tests := []struct {
		name         string
		board        Board
		computerMark string
		wantWin      bool         // If true, computer should find a winning move
		wantBlock    bool         // If true, computer should find a blocking move
		wantCenter   bool         // If true, computer should take center
		wantCorner   bool         // If true, computer should take a corner
		notWantPos   [][2]int     // Positions that should not be selected
		allowedPos   map[int]bool // Map of allowed row*3+col positions
	}{
		{
			name: "computer finds winning move",
			board: Board{
				{"O", "O", ""},
				{"", "", ""},
				{"", "", ""},
			},
			computerMark: "O",
			wantWin:      true,
			allowedPos: map[int]bool{
				2: true, // position (0,2)
			},
		},
		{
			name: "computer blocks opponent win",
			board: Board{
				{"X", "X", ""},
				{"", "", ""},
				{"", "", ""},
			},
			computerMark: "O",
			wantBlock:    true,
			allowedPos: map[int]bool{
				2: true, // position (0,2)
			},
		},
		{
			name: "computer takes center when available",
			board: Board{
				{"X", "", ""},
				{"", "", ""},
				{"", "", ""},
			},
			computerMark: "O",
			wantCenter:   true,
			allowedPos: map[int]bool{
				4: true, // position (1,1)
			},
		},
		{
			name: "computer takes corner when center taken",
			board: Board{
				{"", "", ""},
				{"", "X", ""},
				{"", "", ""},
			},
			computerMark: "O",
			wantCorner:   true,
			notWantPos:   [][2]int{{1, 1}}, // Don't take center
		},
		{
			name: "computer takes any available space when no better options",
			board: Board{
				{"X", "O", "X"},
				{"X", "O", ""},
				{"O", "X", "X"},
			},
			computerMark: "O",
			allowedPos: map[int]bool{
				5: true, // only position (1,2) is available
			},
		},
		{
			name: "computer takes any available middle edge when no better options",
			board: Board{
				{"X", "O", "X"},
				{"", "X", "O"},
				{"O", "X", "X"},
			},
			computerMark: "O",
			allowedPos: map[int]bool{
				3: true, // only position (1,0) is available
			},
		},
		{
			name: "computer takes last available space",
			board: Board{
				{"X", "O", "X"},
				{"O", "X", "O"},
				{"X", "", "X"},
			},
			computerMark: "O",
			allowedPos: map[int]bool{
				7: true, // only position (2,1) is available
			},
		},
		{
			name: "computer tries all strategies before taking edge",
			board: Board{
				{"X", "", "X"},
				{"", "O", ""},
				{"X", "", "O"},
			},
			computerMark: "O",
			allowedPos: map[int]bool{
				1: true, // Position (0,1) - edge position
			},
			notWantPos: [][2]int{
				{1, 1},                         // center already taken
				{0, 0}, {0, 2}, {2, 0}, {2, 2}, // corners already taken or blocked
			},
		},
		{
			name: "computer systematically checks all positions",
			board: Board{
				{"X", "O", "X"},
				{"X", "O", "X"},
				{"", "X", "O"},
			},
			computerMark: "O",
			allowedPos: map[int]bool{
				6: true, // Position (2,0) should be the first empty space found
			},
		},
		{
			name: "computer takes middle edge when all other strategies fail",
			board: Board{
				{"X", "", "O"},
				{"", "X", "X"},
				{"O", "X", "O"},
			},
			computerMark: "O",
			allowedPos: map[int]bool{
				1: true, // Position (0,1) is the only strategic move left
			},
		},
		{
			name: "computer takes first available space when all other options exhausted",
			board: Board{
				{"X", "", "O"},
				{"X", "O", "X"},
				{"O", "X", ""},
			},
			computerMark: "O",
			allowedPos: map[int]bool{
				1: true, // Position (0,1) should be taken as it's the first available space
			},
		},
		{
			name: "computer takes last remaining position",
			board: Board{
				{"X", "O", "X"},
				{"O", "X", "O"},
				{"O", "X", ""},
			},
			computerMark: "O",
			allowedPos: map[int]bool{
				8: true, // Position (2,2) is the only remaining space
			},
		},
		{
			name: "computer takes first empty space in first row",
			board: Board{
				{"", "O", "X"},
				{"X", "O", "X"},
				{"O", "X", "O"},
			},
			computerMark: "O",
			allowedPos: map[int]bool{
				0: true, // Position (0,0)
			},
		},
		{
			name: "computer takes first empty space in second row",
			board: Board{
				{"X", "O", "X"},
				{"", "O", "X"},
				{"O", "X", "O"},
			},
			computerMark: "O",
			allowedPos: map[int]bool{
				3: true, // Position (1,0)
			},
		},
		{
			name: "computer must take edge when no other options exist",
			board: Board{
				{"X", "", "O"}, // corners taken
				{"", "X", ""},  // center taken
				{"O", "", "X"}, // corners taken
			},
			computerMark: "O",
			notWantPos: [][2]int{
				{0, 0}, {0, 2}, // top corners blocked
				{1, 1},         // center blocked
				{2, 0}, {2, 2}, // bottom corners blocked
			},
			allowedPos: map[int]bool{
				1: true, // (0,1)
				3: true, // (1,0)
				5: true, // (1,2)
				7: true, // (2,1)
			},
		},
		{
			name: "computer returns (-1,-1) when board is full",
			board: Board{
				{"X", "O", "X"},
				{"O", "X", "O"},
				{"O", "X", "O"},
			},
			computerMark: "O",
			// No allowed positions as board is full
			// Should return (-1, -1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &Game{
				board:        tt.board,
				ComputerMark: tt.computerMark,
			}

			row, col := game.GetComputerMove()

			// Special case for full board
			if tt.name == "computer returns (-1,-1) when board is full" {
				if row != -1 || col != -1 {
					t.Errorf("GetComputerMove() with full board returned (%d,%d), want (-1,-1)", row, col)
				}
				return
			}

			// Check if position is in allowed positions when specified
			if len(tt.allowedPos) > 0 {
				pos := row*3 + col
				if !tt.allowedPos[pos] {
					t.Errorf("GetComputerMove() returned (%d,%d), position not in allowed set", row, col)
				}
			}

			// Check that move is not in not-wanted positions
			for _, notWant := range tt.notWantPos {
				if row == notWant[0] && col == notWant[1] {
					t.Errorf("GetComputerMove() returned (%d,%d), position should be avoided", row, col)
				}
			}

			// Verify move is valid
			if row < 0 || row > 2 || col < 0 || col > 2 {
				t.Errorf("GetComputerMove() returned invalid position (%d,%d)", row, col)
			}

			// Verify position is empty
			if game.board[row][col] != "" {
				t.Errorf("GetComputerMove() returned occupied position (%d,%d)", row, col)
			}
		})
	}
}
