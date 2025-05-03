package higherlower

import (
	"errors"
	"strings"
	"testing"
)

// mockPrompter implements the Prompter interface for higher/lower game testing.
// It can be configured with either a single response or a sequence of responses.
type mockPrompter struct {
	selectAnswer  int   // Single response for Select calls
	selectError   error // Error to be returned by Select
	selectAnswers []int // Sequence of responses for Select calls
	selectIndex   int   // Current index in selectAnswers
}

// Select implements the Prompter interface by returning predefined responses.
// It first checks for sequence-based responses in selectAnswers, falling back
// to the single selectAnswer if no sequence is available.
func (m *mockPrompter) Select(prompt string, defaultValue string, options []string) (int, error) {
	// If we have a sequence of answers, use those
	if m.selectAnswers != nil && m.selectIndex < len(m.selectAnswers) {
		answer := m.selectAnswers[m.selectIndex]
		m.selectIndex++
		return answer, nil
	}

	// Otherwise fall back to the single answer
	return m.selectAnswer, m.selectError
}

func TestValidateGuess(t *testing.T) {
	validGuesses := []string{"higher", "higher ", " HIGHER", "lower", "LOWER", " lower "}
	invalidGuesses := []string{"", "foo", "123", "hi", "low"}

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

func TestGetPlayerGuess(t *testing.T) {
	tests := []struct {
		name          string
		selectAnswer  int
		selectError   error
		expectedGuess string
		expectedCont  bool
	}{
		{
			name:          "Select Higher",
			selectAnswer:  0,
			selectError:   nil,
			expectedGuess: "higher",
			expectedCont:  true,
		},
		{
			name:          "Select Lower",
			selectAnswer:  1,
			selectError:   nil,
			expectedGuess: "lower",
			expectedCont:  true,
		},
		{
			name:          "Select Quit",
			selectAnswer:  2,
			selectError:   nil,
			expectedGuess: "",
			expectedCont:  false,
		},
		{
			name:          "Prompter Error",
			selectAnswer:  0,
			selectError:   errors.New("prompt error"),
			expectedGuess: "",
			expectedCont:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mockPrompter{
				selectAnswer: tt.selectAnswer,
				selectError:  tt.selectError,
			}

			gotGuess, gotCont := GetPlayerGuess(p, 50)

			if gotGuess != tt.expectedGuess {
				t.Errorf("GetPlayerGuess() guess = %q, want %q", gotGuess, tt.expectedGuess)
			}

			if gotCont != tt.expectedCont {
				t.Errorf("GetPlayerGuess() continue = %v, want %v", gotCont, tt.expectedCont)
			}
		})
	}
}

func TestGamePlay(t *testing.T) {
	// Save original random generator function to restore it later
	originalGenerateNumber := DefaultGenerateNumber
	defer func() { DefaultGenerateNumber = originalGenerateNumber }()

	tests := []struct {
		name           string
		currentNumber  int
		nextNumber     int
		guess          string
		expectedResult bool
		expectedIsOver bool
	}{
		{
			name:           "Higher Guess, Next Number Higher",
			currentNumber:  50,
			nextNumber:     75,
			guess:          "higher",
			expectedResult: true,
			expectedIsOver: false,
		},
		{
			name:           "Higher Guess, Next Number Lower",
			currentNumber:  50,
			nextNumber:     25,
			guess:          "higher",
			expectedResult: false,
			expectedIsOver: true,
		},
		{
			name:           "Lower Guess, Next Number Lower",
			currentNumber:  50,
			nextNumber:     25,
			guess:          "lower",
			expectedResult: true,
			expectedIsOver: false,
		},
		{
			name:           "Lower Guess, Next Number Higher",
			currentNumber:  50,
			nextNumber:     75,
			guess:          "lower",
			expectedResult: false,
			expectedIsOver: true,
		},
		{
			name:           "Equal Numbers",
			currentNumber:  50,
			nextNumber:     50,
			guess:          "higher",
			expectedResult: false,
			expectedIsOver: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up deterministic number generation
			nextNumber := tt.nextNumber
			DefaultGenerateNumber = func(min, max int) int {
				return nextNumber
			}

			game := &Game{
				CurrentNumber: tt.currentNumber,
				MinNumber:     1,
				MaxNumber:     100,
			}

			game.Play(tt.guess)

			if game.IsCorrect != tt.expectedResult {
				t.Errorf("Play() IsCorrect = %v, want %v", game.IsCorrect, tt.expectedResult)
			}

			if game.IsOver != tt.expectedIsOver {
				t.Errorf("Play() IsOver = %v, want %v", game.IsOver, tt.expectedIsOver)
			}
		})
	}
}

func TestNewGame(t *testing.T) {
	minNumber := 1
	maxNumber := 100

	// Test multiple games to ensure current number is within range
	for i := 0; i < 100; i++ {
		game := NewGame(minNumber, maxNumber)

		if game.CurrentNumber < minNumber || game.CurrentNumber > maxNumber {
			t.Errorf("NewGame() CurrentNumber = %d, want between %d and %d",
				game.CurrentNumber, minNumber, maxNumber)
		}

		if game.IsOver {
			t.Errorf("NewGame() IsOver = %v, want %v", game.IsOver, false)
		}
	}
}

func TestUpdateForNextRound(t *testing.T) {
	game := &Game{
		CurrentNumber: 50,
		NextNumber:    75,
	}

	game.UpdateForNextRound()

	if game.CurrentNumber != 75 {
		t.Errorf("UpdateForNextRound() CurrentNumber = %d, want %d", game.CurrentNumber, 75)
	}
}

// Test for GetResult method
func TestGetResult(t *testing.T) {
	tests := []struct {
		name           string
		currentNumber  int
		nextNumber     int
		playerGuess    string
		isCorrect      bool
		expectContains []string
	}{
		{
			name:           "Correct Higher Guess",
			currentNumber:  50,
			nextNumber:     75,
			playerGuess:    "higher",
			isCorrect:      true,
			expectContains: []string{"50", "75", "higher", "Correct"},
		},
		{
			name:           "Incorrect Higher Guess",
			currentNumber:  50,
			nextNumber:     25,
			playerGuess:    "higher",
			isCorrect:      false,
			expectContains: []string{"50", "25", "higher", "Incorrect"},
		},
		{
			name:           "Correct Lower Guess",
			currentNumber:  50,
			nextNumber:     25,
			playerGuess:    "lower",
			isCorrect:      true,
			expectContains: []string{"50", "25", "lower", "Correct"},
		},
		{
			name:           "Incorrect Lower Guess",
			currentNumber:  50,
			nextNumber:     75,
			playerGuess:    "lower",
			isCorrect:      false,
			expectContains: []string{"50", "75", "lower", "Incorrect"},
		},
		{
			name:           "Same Number Case",
			currentNumber:  50,
			nextNumber:     50,
			playerGuess:    "higher",
			isCorrect:      false,
			expectContains: []string{"50", "50", "same"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &Game{
				CurrentNumber: tt.currentNumber,
				NextNumber:    tt.nextNumber,
				PlayerGuess:   tt.playerGuess,
				IsCorrect:     tt.isCorrect,
			}

			result := game.GetResult()

			for _, expected := range tt.expectContains {
				if !contains(result, expected) {
					t.Errorf("GetResult() result does not contain %q\nGot: %q", expected, result)
				}
			}
		})
	}
}

// Helper function for checking if a string contains another string (case insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Test for GenerateNextNumber method
func TestGenerateNextNumber(t *testing.T) {
	// Save original random generator function to restore it later
	originalGenerateNumber := DefaultGenerateNumber
	defer func() { DefaultGenerateNumber = originalGenerateNumber }()

	// Set up deterministic number generator
	expectedNumber := 75
	DefaultGenerateNumber = func(min, max int) int {
		return expectedNumber
	}

	game := &Game{
		MinNumber: 1,
		MaxNumber: 100,
	}

	game.GenerateNextNumber()

	if game.NextNumber != expectedNumber {
		t.Errorf("GenerateNextNumber() set NextNumber = %d, want %d", game.NextNumber, expectedNumber)
	}
}

// Test for PlayGame function
func TestPlayGame(t *testing.T) {
	// Save original random generator function to restore it later
	originalGenerateNumber := DefaultGenerateNumber
	defer func() { DefaultGenerateNumber = originalGenerateNumber }()
	// We're using the mockPrompter that already has selectAnswers capability

	// Test scenarios
	tests := []struct {
		name          string
		selectAnswers []int // Sequence of select answers (0=Higher, 1=Lower, 2=Quit)
		expectedCalls int   // Expected number of times select is called
		numbers       []int // Sequence of generated numbers for testing
	}{
		{
			name:          "Win one round then quit",
			selectAnswers: []int{0, 2}, // Guess Higher, then Quit
			expectedCalls: 2,
			numbers:       []int{50, 75}, // Start with 50, next is 75 (correct higher guess)
		},
		{
			name:          "Win two rounds then quit",
			selectAnswers: []int{0, 1, 2}, // Higher, Lower, Quit
			expectedCalls: 3,
			numbers:       []int{50, 75, 25}, // Start with 50, next 75, then 25
		},
		{
			name:          "Lose immediately",
			selectAnswers: []int{0}, // Guess Higher and lose
			expectedCalls: 1,
			numbers:       []int{50, 25}, // Start with 50, next is 25 (incorrect higher guess)
		},
		{
			name:          "Quit immediately",
			selectAnswers: []int{2}, // Quit
			expectedCalls: 1,
			numbers:       []int{50},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock prompter with sequence of responses
			mp := &mockPrompter{
				selectAnswers: tt.selectAnswers,
			}

			// Setup deterministic number generation
			numIndex := 0
			DefaultGenerateNumber = func(min, max int) int {
				if numIndex < len(tt.numbers) {
					result := tt.numbers[numIndex]
					numIndex++
					return result
				}
				return 50 // Default fallback
			}

			// This test validates that PlayGame executes with the configured mock prompter
			// The deterministic number generation allows us to control the game flow
			PlayGame(mp, 1, 100)
		})
	}
}
