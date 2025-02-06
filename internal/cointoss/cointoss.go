package cointoss

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	ghprompt "github.com/cli/go-gh/v2/pkg/prompter"
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
	options := []string{"Heads", "Tails", "Quit"}
	p := ghprompt.New(os.Stdin, os.Stdout, os.Stderr)

	answer, err := p.Select("What's your next guess? Heads, Tails or Quit?", "Heads", options)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return "", false
	}

	answerLower := strings.ToLower(strings.TrimSpace(options[answer]))
	if answerLower == "quit" {
		return "", false
	}

	return answerLower, true
}
