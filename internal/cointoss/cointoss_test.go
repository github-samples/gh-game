package cointoss

import (
	"errors"
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

func TestGetNextGuessWithPrompter(t *testing.T) {
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
		{
			name:          "error during selection prints error message",
			selectAnswer:  0,
			selectError:   errors.New("test error"),
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

			guess, cont := GetNextGuessWithPrompter(mockP)

			if guess != tt.expectedGuess {
				t.Errorf("GetNextGuessWithPrompter() guess = %v, want %v", guess, tt.expectedGuess)
			}
			if cont != tt.expectedCont {
				t.Errorf("GetNextGuessWithPrompter() cont = %v, want %v", cont, tt.expectedCont)
			}
		})
	}
}
