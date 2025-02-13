package tictactoe

import (
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
	answer := m.selectAnswers[m.selectIndex]
	m.selectIndex++
	return answer, nil
}

func TestNewGame(t *testing.T) {
	game := NewGame()
	if game.CurrentPlayer != "X" {
		t.Errorf("Expected first player to be X, got %s", game.CurrentPlayer)
	}
}

func TestMakeMove(t *testing.T) {
	game := NewGame()

	// Test valid move
	err := game.MakeMove(0, 0)
	if err != nil {
		t.Errorf("Expected valid move, got error: %v", err)
	}
	if game.board[0][0] != "X" {
		t.Errorf("Expected X at position (0,0), got %s", game.board[0][0])
	}

	// Test invalid position
	err = game.MakeMove(3, 3)
	if err == nil {
		t.Error("Expected error for invalid position, got nil")
	}

	// Test already taken position
	err = game.MakeMove(0, 0)
	if err == nil {
		t.Error("Expected error for taken position, got nil")
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
	game := NewGame()
	if game.IsBoardFull() {
		t.Error("Expected empty board not to be full")
	}

	moves := [][2]int{
		{0, 0}, {0, 1}, {0, 2},
		{1, 0}, {1, 1}, {1, 2},
		{2, 0}, {2, 1}, {2, 2},
	}

	for _, move := range moves {
		game.MakeMove(move[0], move[1])
	}

	if !game.IsBoardFull() {
		t.Error("Expected board to be full")
	}
}

func TestGetPlayerMove(t *testing.T) {
	game := NewGame()
	mockP := &mockPrompter{
		selectAnswers: []int{4}, // Select position 5 (center)
	}

	row, col, err := GetPlayerMove(mockP, game)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if row != 1 || col != 1 {
		t.Errorf("Expected move (1,1), got (%d,%d)", row, col)
	}

	// Test when position is taken
	game.MakeMove(1, 1) // Take center position
	mockP = &mockPrompter{
		selectAnswers: []int{0}, // Select first available position
	}

	row, col, err = GetPlayerMove(mockP, game)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if row != 0 || col != 0 {
		t.Errorf("Expected move (0,0), got (%d,%d)", row, col)
	}
}

func TestPositionToRowCol(t *testing.T) {
	tests := []struct {
		position int
		row      int
		col      int
	}{
		{1, 0, 0},
		{2, 0, 1},
		{3, 0, 2},
		{4, 1, 0},
		{5, 1, 1},
		{6, 1, 2},
		{7, 2, 0},
		{8, 2, 1},
		{9, 2, 2},
	}

	for _, tt := range tests {
		row, col := positionToRowCol(tt.position)
		if row != tt.row || col != tt.col {
			t.Errorf("positionToRowCol(%d) = (%d,%d); want (%d,%d)",
				tt.position, row, col, tt.row, tt.col)
		}
	}
}

func TestGetAvailablePositions(t *testing.T) {
	game := NewGame()
	positions := game.getAvailablePositions()
	if len(positions) != 9 {
		t.Errorf("Expected 9 available positions in empty board, got %d", len(positions))
	}

	// Make some moves
	game.MakeMove(0, 0) // Position 1
	game.MakeMove(1, 1) // Position 5
	game.MakeMove(2, 2) // Position 9

	positions = game.getAvailablePositions()
	if len(positions) != 6 {
		t.Errorf("Expected 6 available positions after three moves, got %d", len(positions))
	}

	// Verify positions 1, 5, and 9 are not in the available positions
	for _, pos := range positions {
		if pos == "1" || pos == "5" || pos == "9" {
			t.Errorf("Position %s should not be available", pos)
		}
	}
}
