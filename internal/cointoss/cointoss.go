package cointoss

import (
	"fmt"
	"math/rand"
	"strings"
)

// Game represents the state of a coin toss game
type Game struct {
	PlayerGuess string
	Result      string
	IsOver      bool
}

// prompter interface allows us to mock the prompt functionality in tests
type prompter interface {
	Select(prompt string, defaultValue string, options []string) (int, error)
}

func NewGame() *Game {
	return &Game{
		IsOver: false,
	}
}

// TossCoin is a variable so it can be replaced in tests
var TossCoin = func() string {
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

// GetPlayerGuess gets the player's next guess using the provided prompter
func GetPlayerGuess(p prompter) (string, bool) {
	options := []string{"Heads", "Tails", "Quit"}

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

// Play executes a round of the coin toss game
func (g *Game) Play(guess string) {
	g.PlayerGuess = guess
	g.Result = TossCoin()
	g.IsOver = true
}

// GetResult returns the game result message
func (g *Game) GetResult() string {
	if g.PlayerGuess == g.Result {
		return fmt.Sprintf("You guessed %s and the coin landed on %s. You win!", g.PlayerGuess, g.Result)
	}
	return fmt.Sprintf("You guessed %s but the coin landed on %s. You lose!", g.PlayerGuess, g.Result)
}

// PlayGame handles the main game loop
func PlayGame(p prompter, initialGuess string) {
	game := NewGame()
	streak := 0
	keepPlaying := true
	guess := strings.ToLower(strings.TrimSpace(initialGuess))

	for keepPlaying {
		game.Play(guess)
		fmt.Println(game.GetResult())

		if game.PlayerGuess == game.Result {
			streak++
			fmt.Printf("Streak: %d\n", streak)
			guess, keepPlaying = GetPlayerGuess(p)
		} else {
			fmt.Printf("Game Over! Final streak: %d\n", streak)
			keepPlaying = false
		}
	}
}
