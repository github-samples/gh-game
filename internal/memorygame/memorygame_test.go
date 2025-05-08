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
