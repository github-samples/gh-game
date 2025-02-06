package cointoss

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func TossCoin() string {
	rand.Seed(time.Now().UnixNano())
	if rand.Float32() < 0.5 {
		return "heads"
	}
	return "tails"
}

func ValidateGuess(guess string) error {
	guess = strings.ToLower(strings.TrimSpace(guess))
	if guess != "heads" && guess != "tails" {
		return fmt.Errorf("guess must be either 'heads' or 'tails'")
	}
	return nil
}

func GetNextGuess() (string, bool) {
	fmt.Print("Play again? Enter 'heads' or 'tails' (or 'quit' to end): ")
	var answer string
	fmt.Scanln(&answer)

	if strings.ToLower(strings.TrimSpace(answer)) == "quit" {
		return "", false
	}

	if err := ValidateGuess(answer); err != nil {
		fmt.Println(err)
		return GetNextGuess()
	}

	return strings.ToLower(strings.TrimSpace(answer)), true
}
