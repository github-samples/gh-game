package cointoss

import (
	"testing"
)

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
