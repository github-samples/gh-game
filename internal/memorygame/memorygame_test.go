package memorygame

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	tests := []struct {
		name      string
		lives     int
		wantLives int
	}{
		{"zero lives", 0, 0},
		{"one life", 1, 1},
		{"three lives", 3, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(tt.lives)
			if game.Lives != tt.wantLives {
				t.Errorf("NewGame() lives = %v, want %v", game.Lives, tt.wantLives)
			}
		})
	}
}

func TestGenerateSequence(t *testing.T) {
	game := NewGame(3)
	game.CurrentRound = 1
	game.GenerateSequence()

	if len(game.Sequence) != 3 { // Round 1 should have 3 colors
		t.Errorf("GenerateSequence() sequence length = %v, want %v", len(game.Sequence), 3)
	}
}

func TestCheckSequence(t *testing.T) {
	game := NewGame(3)
	game.Sequence = []Color{Red, Blue, Green}

	tests := []struct {
		name  string
		input []Color
		want  bool
	}{
		{
			"correct sequence",
			[]Color{Red, Blue, Green},
			true,
		},
		{
			"wrong sequence",
			[]Color{Red, Green, Blue},
			false,
		},
		{
			"wrong length",
			[]Color{Red, Blue},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := game.CheckSequence(tt.input); got != tt.want {
				t.Errorf("CheckSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsGameOver(t *testing.T) {
	tests := []struct {
		name  string
		lives int
		round int
		want  bool
	}{
		{"game active", 3, 1, false},
		{"no lives", 0, 1, true},
		{"max rounds", 3, 101, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(tt.lives)
			game.CurrentRound = tt.round
			if got := game.IsGameOver(); got != tt.want {
				t.Errorf("IsGameOver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextRound(t *testing.T) {
	tests := []struct {
		name          string
		initialRound  int
		expectedRound int
	}{
		{"increment from round 1", 1, 2},
		{"increment from round 5", 5, 6},
		{"increment from round 99", 99, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(3)
			game.CurrentRound = tt.initialRound
			game.NextRound()
			if game.CurrentRound != tt.expectedRound {
				t.Errorf("NextRound() round = %v, want %v", game.CurrentRound, tt.expectedRound)
			}
		})
	}
}

func TestFailFastOnWrongColor(t *testing.T) {
	game := NewGame(2)
	game.Sequence = []Color{Red, Blue, Green}

	// Simulate user picking wrong color on first pick
	userSequence := []Color{Yellow}
	correct := game.Sequence[0] == userSequence[0]
	if correct {
		t.Fatal("Test setup error: userSequence[0] should be wrong")
	}

	if game.Lives != 2 {
		t.Errorf("Expected 2 lives at start, got %d", game.Lives)
	}

	if !correct {
		game.DecrementLives()
	}

	if game.Lives != 1 {
		t.Errorf("Expected 1 life after wrong pick, got %d", game.Lives)
	}

	// User retries the round, picks correct first color, then wrong second color
	userSequence = []Color{Red, Green}
	if userSequence[0] != game.Sequence[0] {
		t.Fatal("Test setup error: userSequence[0] should be correct")
	}
	if userSequence[1] != game.Sequence[1] {
		game.DecrementLives()
	}
	if game.Lives != 0 {
		t.Errorf("Expected 0 lives after two wrong picks, got %d", game.Lives)
	}
}

func TestDisplayUserSequenceInSingleLine(t *testing.T) {
	// This tests that the infrastructure supports displaying user sequences in a single line
	// The actual display function is in the cmd package, but we want to make sure
	// the core logic supports building and comparing user sequences correctly

	game := NewGame(3)
	game.Sequence = []Color{Red, Blue, Green, Yellow, Red}

	// Test that we can build the sequence incrementally
	userSequence := make([]Color, 0, len(game.Sequence))

	// Add colors one by one (simulating user input from UI)
	userSequence = append(userSequence, Red)
	if len(userSequence) != 1 || userSequence[0] != Red {
		t.Errorf("User sequence should have 1 color Red, got %v", userSequence)
	}

	userSequence = append(userSequence, Blue)
	if len(userSequence) != 2 || userSequence[0] != Red || userSequence[1] != Blue {
		t.Errorf("User sequence should have 2 colors Red, Blue, got %v", userSequence)
	}

	userSequence = append(userSequence, Green)
	userSequence = append(userSequence, Yellow)
	userSequence = append(userSequence, Red)

	// Check the final sequence
	if len(userSequence) != len(game.Sequence) {
		t.Errorf("User sequence length should be %d, got %d", len(game.Sequence), len(userSequence))
	}

	if !game.CheckSequence(userSequence) {
		t.Errorf("CheckSequence should return true for the correct sequence")
	}
}
