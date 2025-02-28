package cointoss

import (
	"errors"
	"strings"
	"testing"
)

// Mock prompter for testing
type mockPrompter struct {
	selectAnswer int
	selectError  error
}

func (m *mockPrompter) Select(prompt string, defaultValue string, options []string) (int, error) {
	return m.selectAnswer, m.selectError
}

func TestValidateGuess(t *testing.T) {
	validGuesses := []string{"heads", "heads ", " HEADS", "tails", "TAILS", " tails "}
	invalidGuesses := []string{"", "foo", "123", "heds", "taol"}

	for _, guess := range validGuesses {
		if err := ValidateGuess(guess); err != nil {
			t.Errorf("ValidateGuess(%q) returned error: %v, expected nil", guess, err)
		}
	}

	for _, guess := range invalidGuesses {
		if err := ValidateGuess(guess); err == nil {
			t.Errorf("ValidateGuess(%q) did not return error, expected error", guess)
		}
	}
}

func TestTossCoin(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := TossCoin()
		if result != "heads" && result != "tails" {
			t.Errorf("TossCoin() returned unexpected value %q", result)
		}
	}
}

func TestGetPlayerGuess(t *testing.T) {
	tests := []struct {
		name          string
		selectAnswer  int
		selectError   error
		expectedGuess string
		expectedCont  bool
	}{
		{
			name:          "select heads",
			selectAnswer:  0,
			selectError:   nil,
			expectedGuess: "heads",
			expectedCont:  true,
		},
		{
			name:          "select tails",
			selectAnswer:  1,
			selectError:   nil,
			expectedGuess: "tails",
			expectedCont:  true,
		},
		{
			name:          "select quit",
			selectAnswer:  2,
			selectError:   nil,
			expectedGuess: "",
			expectedCont:  false,
		},
		{
			name:          "error during selection",
			selectAnswer:  0,
			selectError:   errors.New("mock error"),
			expectedGuess: "",
			expectedCont:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockP := &mockPrompter{
				selectAnswer: tt.selectAnswer,
				selectError:  tt.selectError,
			}
			guess, cont := GetPlayerGuess(mockP)
			if guess != tt.expectedGuess {
				t.Errorf("GetPlayerGuess() guess = %v, want %v", guess, tt.expectedGuess)
			}
			if cont != tt.expectedCont {
				t.Errorf("GetPlayerGuess() cont = %v, want %v", cont, tt.expectedCont)
			}
		})
	}
}

func TestGame_Play(t *testing.T) {
	game := NewGame()

	// Test initial state
	if game.IsOver {
		t.Error("New game should not be over")
	}

	// Play a round
	game.Play("heads")

	// Test that game state is updated
	if !game.IsOver {
		t.Error("Game should be over after playing")
	}
	if game.PlayerGuess != "heads" {
		t.Errorf("PlayerGuess = %v, want heads", game.PlayerGuess)
	}
	if game.Result != "heads" && game.Result != "tails" {
		t.Errorf("Result = %v, want either heads or tails", game.Result)
	}
}

func TestGame_GetResult(t *testing.T) {
	tests := []struct {
		name        string
		playerGuess string
		result      string
		wantWin     bool
	}{
		{
			name:        "player wins with heads",
			playerGuess: "heads",
			result:      "heads",
			wantWin:     true,
		},
		{
			name:        "player wins with tails",
			playerGuess: "tails",
			result:      "tails",
			wantWin:     true,
		},
		{
			name:        "player loses with heads",
			playerGuess: "heads",
			result:      "tails",
			wantWin:     false,
		},
		{
			name:        "player loses with tails",
			playerGuess: "tails",
			result:      "heads",
			wantWin:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &Game{
				PlayerGuess: tt.playerGuess,
				Result:      tt.result,
				IsOver:      true,
			}
			got := game.GetResult()
			if tt.wantWin && !contains(got, "You win!") {
				t.Errorf("GetResult() = %v, want win message", got)
			}
			if !tt.wantWin && !contains(got, "You lose!") {
				t.Errorf("GetResult() = %v, want lose message", got)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestPlayGame(t *testing.T) {
	tests := []struct {
		name         string
		selectAnswer int
		selectError  error
		initialGuess string
		results      []string // sequence of coin flip results to test
	}{
		{
			name:         "win first round then quit",
			selectAnswer: 2, // quit
			initialGuess: "heads",
			results:      []string{"heads"},
		},
		{
			name:         "lose first round",
			initialGuess: "heads",
			results:      []string{"tails"},
		},
		{
			name:         "win twice then lose",
			selectAnswer: 0, // heads
			initialGuess: "heads",
			results:      []string{"heads", "heads", "tails"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockP := &mockPrompter{
				selectAnswer: tt.selectAnswer,
				selectError:  tt.selectError,
			}

			// Override TossCoin for deterministic testing
			resultIndex := 0
			oldTossCoin := TossCoin
			TossCoin = func() string {
				result := tt.results[resultIndex]
				if resultIndex < len(tt.results)-1 {
					resultIndex++
				}
				return result
			}
			defer func() { TossCoin = oldTossCoin }()

			PlayGame(mockP, tt.initialGuess)
		})
	}
}
