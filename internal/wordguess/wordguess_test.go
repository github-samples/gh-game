package wordguess

import (
	"strings"
	"testing"
)

// MockPrompter is a mock implementation of the Prompter interface for testing.
// It provides predefined responses for input, select, and confirm prompts to
// enable deterministic testing of game interaction flows.
type MockPrompter struct {
	InputResponses   []string // Predefined responses for Input calls
	InputIndex       int      // Current index in the InputResponses slice
	SelectResponses  []int    // Predefined responses for Select calls
	SelectIndex      int      // Current index in the SelectResponses slice
	ConfirmResponses []bool   // Predefined responses for Confirm calls
	ConfirmIndex     int      // Current index in the ConfirmResponses slice
}

// Input implements the Prompter interface by returning predefined responses
// from InputResponses. It advances the InputIndex after each call.
func (m *MockPrompter) Input(prompt string, defaultValue string) (string, error) {
	if m.InputIndex < len(m.InputResponses) {
		result := m.InputResponses[m.InputIndex]
		m.InputIndex++
		return result, nil
	}
	return "", nil
}

// Select implements the Prompter interface by returning predefined responses
// from SelectResponses. It advances the SelectIndex after each call.
func (m *MockPrompter) Select(prompt string, defaultValue string, options []string) (int, error) {
	if m.SelectIndex < len(m.SelectResponses) {
		result := m.SelectResponses[m.SelectIndex]
		m.SelectIndex++
		return result, nil
	}
	return 0, nil
}

// Confirm implements the Prompter interface by returning predefined responses
// from ConfirmResponses. It advances the ConfirmIndex after each call.
func (m *MockPrompter) Confirm(prompt string, defaultValue bool) (bool, error) {
	if m.ConfirmIndex < len(m.ConfirmResponses) {
		result := m.ConfirmResponses[m.ConfirmIndex]
		m.ConfirmIndex++
		return result, nil
	}
	return false, nil
}

func TestNewGame(t *testing.T) {
	game := NewGame()

	if game == nil {
		t.Fatal("Expected NewGame() to return a game instance, got nil")
	}

	if game.IsOver {
		t.Error("Expected new game to not be over")
	}

	if game.HasWon {
		t.Error("Expected new game to not be won")
	}

	if game.IncorrectGuesses != 0 {
		t.Errorf("Expected new game to have 0 incorrect guesses, got %d", game.IncorrectGuesses)
	}

	if len(game.Word) == 0 {
		t.Error("Expected game word to not be empty")
	}

	if game.RevealedWord != strings.Repeat("_", len(game.Word)) {
		t.Errorf("Expected revealed word to be %s, got %s",
			strings.Repeat("_", len(game.Word)), game.RevealedWord)
	}
}

func TestGuessLetter_ValidGuess(t *testing.T) {
	game := NewGame()
	originalWord := game.Word

	// Always use a letter we know is in the word for this test
	letter := string(originalWord[0])

	err := game.GuessLetter(letter)

	if err != nil {
		t.Errorf("Expected no error for valid guess, got %v", err)
	}

	// Check that the letter was added to guessed letters
	found := false
	for _, guessed := range game.GuessedLetters {
		if guessed == letter {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected %s to be added to guessed letters", letter)
	}

	// Check that the letter is revealed in the word
	if !strings.Contains(game.RevealedWord, letter) {
		t.Errorf("Expected %s to be revealed in the word", letter)
	}
}

func TestGuessLetter_InvalidGuess(t *testing.T) {
	game := NewGame()

	// Test invalid input (not a letter)
	err := game.GuessLetter("1")
	if err == nil {
		t.Error("Expected error for non-letter guess")
	}

	// Test invalid input (multiple letters)
	err = game.GuessLetter("ab")
	if err == nil {
		t.Error("Expected error for multi-letter guess")
	}

	// Test duplicate guess
	letter := "x"
	game.GuessedLetters = append(game.GuessedLetters, letter)

	err = game.GuessLetter(letter)
	if err == nil {
		t.Error("Expected error for duplicate guess")
	}
}

func TestGameOver_Win(t *testing.T) {
	game := NewGame()
	game.Word = "test"
	// Set the RevealedWord to have only one underscore left
	game.RevealedWord = "tes_"

	// Guess the remaining letter 't'
	err := game.GuessLetter("t")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !game.IsOver {
		t.Error("Expected game to be over after guessing final letter")
	}

	if !game.HasWon {
		t.Error("Expected player to have won after guessing final letter")
	}
}

func TestGameOver_Lose(t *testing.T) {
	game := NewGame()
	game.Word = "test"
	game.IncorrectGuesses = MaxIncorrectGuesses - 1

	// Make an incorrect guess
	err := game.GuessLetter("z")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !game.IsOver {
		t.Error("Expected game to be over after max incorrect guesses")
	}

	if game.HasWon {
		t.Error("Expected player to have lost after max incorrect guesses")
	}
}

func TestGetRemainingLetters(t *testing.T) {
	game := NewGame()
	game.GuessedLetters = []string{"a", "e", "i", "o", "u"}

	remaining := game.GetRemainingLetters()

	// Check that guessed letters are not in remaining
	for _, letter := range game.GuessedLetters {
		if strings.Contains(remaining, letter) {
			t.Errorf("Expected %s to not be in remaining letters", letter)
		}
	}

	// Check that the remaining count is correct
	expectedCount := 26 - len(game.GuessedLetters) // 26 letters in alphabet minus guessed
	if len(remaining) != expectedCount {
		t.Errorf("Expected %d remaining letters, got %d", expectedCount, len(remaining))
	}
}

// Test isLetter function
func TestIsLetter(t *testing.T) {
	tests := []struct {
		input    byte
		expected bool
	}{
		{'a', true},
		{'z', true},
		{'A', true},
		{'Z', true},
		{'m', true},
		{'0', false},
		{'9', false},
		{' ', false},
		{'!', false},
		{'@', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			if got := isLetter(tt.input); got != tt.expected {
				t.Errorf("isLetter(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// Test for edge cases in GuessLetter
func TestGuessLetter_EdgeCases(t *testing.T) {
	tests := []struct {
		name             string
		word             string
		revealed         string
		guessed          []string
		newGuess         string
		expectErr        bool
		expectWin        bool
		expectIncorrect  bool
		expectedRevealed string
	}{
		{
			name:             "Empty guess",
			word:             "test",
			revealed:         "____",
			guessed:          []string{},
			newGuess:         "",
			expectErr:        true,
			expectWin:        false,
			expectIncorrect:  false,
			expectedRevealed: "____", // No change expected
		},
		{
			name:             "Multiple letter guess",
			word:             "test",
			revealed:         "____",
			guessed:          []string{},
			newGuess:         "te",
			expectErr:        true,
			expectWin:        false,
			expectIncorrect:  false,
			expectedRevealed: "____", // No change expected
		},
		{
			name:             "Non-letter guess",
			word:             "test",
			revealed:         "____",
			guessed:          []string{},
			newGuess:         "1",
			expectErr:        true,
			expectWin:        false,
			expectIncorrect:  false,
			expectedRevealed: "____", // No change expected
		},
		{
			name:             "Duplicate guess",
			word:             "test",
			revealed:         "t___",
			guessed:          []string{"t"},
			newGuess:         "t",
			expectErr:        true,
			expectWin:        false,
			expectIncorrect:  false,
			expectedRevealed: "t___", // No change expected
		},
		{
			name:             "Guess reveals all remaining letters (win)",
			word:             "date",
			revealed:         "dat_",
			guessed:          []string{"d", "a", "t"},
			newGuess:         "e",
			expectErr:        false,
			expectWin:        true,
			expectIncorrect:  false,
			expectedRevealed: "date", // Should reveal the 'e'
		},
		{
			name:             "Mixed case guess",
			word:             "test",
			revealed:         "____",
			guessed:          []string{},
			newGuess:         "T",
			expectErr:        false,
			expectWin:        false,
			expectIncorrect:  false,
			expectedRevealed: "t__t", // Should reveal both instances of 't'
		},
		{
			name:             "Correct guess with multiple instances of letter",
			word:             "pepper",
			revealed:         "______",
			guessed:          []string{},
			newGuess:         "p",
			expectErr:        false,
			expectWin:        false,
			expectIncorrect:  false,
			expectedRevealed: "p_pp__", // Should reveal all instances of 'p'
		},
		{
			name:             "Incorrect guess increases incorrect count",
			word:             "github",
			revealed:         "g_____",
			guessed:          []string{"g"},
			newGuess:         "z",
			expectErr:        false,
			expectWin:        false,
			expectIncorrect:  true,
			expectedRevealed: "g_____", // No change in revealed word
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &Game{
				Word:             tt.word,
				RevealedWord:     tt.revealed,
				GuessedLetters:   tt.guessed,
				IncorrectGuesses: 0, // Start with no incorrect guesses
			}

			// Remember initial state
			initialIncorrect := game.IncorrectGuesses

			err := game.GuessLetter(tt.newGuess)

			// Check error expectation
			if (err != nil) != tt.expectErr {
				t.Errorf("GuessLetter() error = %v, expectErr %v", err, tt.expectErr)
			}

			// If we expect a win, check if the game is over and won
			if tt.expectWin {
				if !game.IsOver {
					t.Errorf("Expected game to be over, but it's not")
				}
				if !game.HasWon {
					t.Errorf("Expected player to have won, but hasn't")
				}
			}

			// Check if incorrect guesses increased as expected
			if tt.expectIncorrect && game.IncorrectGuesses <= initialIncorrect {
				t.Errorf("Expected incorrect guesses to increase, but got %d (was %d)",
					game.IncorrectGuesses, initialIncorrect)
			} else if !tt.expectIncorrect && !tt.expectErr && game.IncorrectGuesses > initialIncorrect {
				t.Errorf("Didn't expect incorrect guesses to increase, but got %d (was %d)",
					game.IncorrectGuesses, initialIncorrect)
			}

			// Check if the revealed word matches expectations
			if game.RevealedWord != tt.expectedRevealed {
				t.Errorf("Revealed word = %q, expected %q", game.RevealedWord, tt.expectedRevealed)
			}
		})
	}
}

// Test String method (game display)
func TestString(t *testing.T) {
	tests := []struct {
		name             string
		word             string
		revealedWord     string
		guessedLetters   []string
		incorrectGuesses int
		isOver           bool
		hasWon           bool
		expectContains   []string
	}{
		{
			name:             "Game in progress",
			word:             "github",
			revealedWord:     "g_t___",
			guessedLetters:   []string{"g", "t", "a"},
			incorrectGuesses: 1,
			isOver:           false,
			hasWon:           false,
			expectContains:   []string{"W O R D  G U E S S", "Guesses Remaining: 5/6", "g _ t _ _ _", "Guess a letter"},
		},
		{
			name:             "Game won",
			word:             "github",
			revealedWord:     "github",
			guessedLetters:   []string{"g", "i", "t", "h", "u", "b"},
			incorrectGuesses: 2,
			isOver:           true,
			hasWon:           true,
			expectContains:   []string{"W O R D  G U E S S", "Guesses Remaining: 4/6", "g i t h u b", "Congratulations"},
		},
		{
			name:             "Game lost",
			word:             "github",
			revealedWord:     "g_th__",
			guessedLetters:   []string{"g", "t", "h", "a", "e", "o", "u"},
			incorrectGuesses: MaxIncorrectGuesses,
			isOver:           true,
			hasWon:           false,
			expectContains:   []string{"W O R D  G U E S S", "Guesses Remaining: 0/6", "g _ t h _ _", "Game over", "github"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &Game{
				Word:             tt.word,
				RevealedWord:     tt.revealedWord,
				GuessedLetters:   tt.guessedLetters,
				IncorrectGuesses: tt.incorrectGuesses,
				IsOver:           tt.isOver,
				HasWon:           tt.hasWon,
			}

			result := game.String()

			for _, expected := range tt.expectContains {
				if !strings.Contains(result, expected) {
					t.Errorf("String() output does not contain %q\nGot: %q", expected, result)
				}
			}
		})
	}
}

// Test PlayGame with mock prompter
func TestPlayGame(t *testing.T) {
	// Save the original WordList and restore it after the test
	originalWordList := WordList
	defer func() { WordList = originalWordList }()

	// Use a deterministic word for testing
	WordList = []string{"test"}

	tests := []struct {
		name             string
		inputResponses   []string
		confirmResponses []bool
	}{
		{
			name:             "Win game",
			inputResponses:   []string{"t", "e", "s"},
			confirmResponses: []bool{false},
		},
		{
			name:             "Lose game",
			inputResponses:   []string{"a", "b", "c", "d", "f", "g", "h"},
			confirmResponses: []bool{false},
		},
		{
			name:             "Play again",
			inputResponses:   []string{"t", "e", "s", "a", "b", "c", "d", "f", "g"},
			confirmResponses: []bool{true, false}, // First true (play again), then false (quit)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &MockPrompter{
				InputResponses:   tt.inputResponses,
				ConfirmResponses: tt.confirmResponses,
			}

			// This test checks that the function runs without errors
			// Additional validation is done below with the confirm call count check
			PlayGame(mp)

			// Verify confirm was called the expected number of times
			expectedConfirms := len(tt.confirmResponses)
			if mp.ConfirmIndex != expectedConfirms {
				t.Errorf("Expected confirm to be called %d times, got %d",
					expectedConfirms, mp.ConfirmIndex)
			}
		})
	}
}
