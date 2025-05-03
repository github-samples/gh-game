// Package higherlower implements the Higher or Lower game
package higherlower

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Define styles for the game
var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))  // purple
	numberStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))  // blue
	correctStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))             // green
	incorrectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))              // red
	guessStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))            // yellow
	streakStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("208")) // orange
)

// Game represents the state of a Higher or Lower game
type Game struct {
	CurrentNumber int
	NextNumber    int
	PlayerGuess   string
	IsCorrect     bool
	IsOver        bool
	MinNumber     int
	MaxNumber     int
}

// prompter interface allows us to mock the prompt functionality in tests
type prompter interface {
	Select(prompt string, defaultValue string, options []string) (int, error)
}

// NewGame creates a new Higher or Lower game with default settings
func NewGame(minNumber, maxNumber int) *Game {
	// Generate a local random generator with its own source
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate initial number
	currentNumber := rng.Intn(maxNumber-minNumber+1) + minNumber

	return &Game{
		CurrentNumber: currentNumber,
		NextNumber:    0,
		IsOver:        false,
		MinNumber:     minNumber,
		MaxNumber:     maxNumber,
	}
}

// ValidateGuess checks if the guess is valid ("higher" or "lower")
func ValidateGuess(guess string) error {
	guess = strings.ToLower(strings.TrimSpace(guess))
	if guess != "higher" && guess != "lower" {
		return fmt.Errorf("guess must be either 'higher' or 'lower'")
	}
	return nil
}

// GetPlayerGuess gets the player's next guess using the provided prompter
func GetPlayerGuess(p prompter, currentNumber int) (string, bool) {
	options := []string{"Higher", "Lower", "Quit"}
	prompt := fmt.Sprintf("Current number is %d. Will the next number be Higher or Lower?", currentNumber)

	answer, err := p.Select(prompt, "Higher", options)
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

// GenerateNumberFunc is a function type for generating random numbers
type GenerateNumberFunc func(min, max int) int

// DefaultGenerateNumber generates a random number between min and max (inclusive)
var DefaultGenerateNumber = func(min, max int) int {
	// Create a new random source each time for better randomness
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rng.Intn(max-min+1) + min
}

// GenerateNextNumber produces the next random number for the game
func (g *Game) GenerateNextNumber() {
	g.NextNumber = DefaultGenerateNumber(g.MinNumber, g.MaxNumber)
}

// Play executes a round of the Higher or Lower game
func (g *Game) Play(guess string) {
	g.PlayerGuess = strings.ToLower(strings.TrimSpace(guess))
	g.GenerateNextNumber()

	// Handle same number case - counts as incorrect
	if g.NextNumber == g.CurrentNumber {
		g.IsCorrect = false
		g.IsOver = true
		return
	}

	// Determine if player's guess is correct
	if g.PlayerGuess == "higher" {
		g.IsCorrect = g.NextNumber > g.CurrentNumber
	} else { // guess is "lower"
		g.IsCorrect = g.NextNumber < g.CurrentNumber
	}

	// Only set game to over if player is incorrect
	g.IsOver = !g.IsCorrect
}

// GetResult returns the game result message
func (g *Game) GetResult() string {
	var outcomeStyle lipgloss.Style
	outcome := "Correct"
	if !g.IsCorrect {
		outcome = "Incorrect"
		outcomeStyle = incorrectStyle
	} else {
		outcomeStyle = correctStyle
	}

	currentNum := numberStyle.Render(fmt.Sprintf("%d", g.CurrentNumber))
	nextNum := numberStyle.Render(fmt.Sprintf("%d", g.NextNumber))
	guess := guessStyle.Render(g.PlayerGuess)
	outcomeText := outcomeStyle.Render(outcome)

	// Special message for same number case
	if g.NextNumber == g.CurrentNumber {
		return fmt.Sprintf(
			"Current number: %s, Next number: %s\n"+
				"The numbers are the same! Game over!",
			currentNum, nextNum,
		)
	}

	return fmt.Sprintf(
		"Current number: %s, Next number: %s\n"+
			"You guessed the next number would be %s: %s!",
		currentNum, nextNum, guess, outcomeText,
	)
}

// UpdateForNextRound prepares the game for the next round
func (g *Game) UpdateForNextRound() {
	g.CurrentNumber = g.NextNumber
}

// PlayGame handles the main game loop
func PlayGame(p prompter, minNumber, maxNumber int) {
	title := titleStyle.Render("Welcome to Higher or Lower!")
	rangeText := fmt.Sprintf("Numbers range from %s to %s",
		numberStyle.Render(fmt.Sprintf("%d", minNumber)),
		numberStyle.Render(fmt.Sprintf("%d", maxNumber)))

	fmt.Printf("%s %s\n\n", title, rangeText)

	// Display game rules
	rules := []string{
		"Rules:",
		"1. You'll be shown a random number",
		"2. Guess if the next number will be HIGHER or LOWER",
		"3. If you guess correctly, you continue and build your streak",
		"4. If you guess incorrectly, the game ends",
		"5. If the numbers are the same, the game ends",
	}

	for _, rule := range rules {
		fmt.Println(rule)
	}
	fmt.Println()

	game := NewGame(minNumber, maxNumber)
	streak := 0

	startingNumber := numberStyle.Render(fmt.Sprintf("%d", game.CurrentNumber))
	fmt.Printf("Starting number: %s\n", startingNumber)

	// Get initial guess from user
	guess, keepPlaying := GetPlayerGuess(p, game.CurrentNumber)

	for keepPlaying {
		game.Play(guess)
		fmt.Println(game.GetResult())

		if game.IsCorrect {
			streak++
			streakText := streakStyle.Render(fmt.Sprintf("Streak: %d", streak))
			fmt.Printf("%s\n", streakText)
			game.UpdateForNextRound()
			guess, keepPlaying = GetPlayerGuess(p, game.CurrentNumber)
		} else {
			gameOver := incorrectStyle.Render("Game Over!")
			finalStreak := streakStyle.Render(fmt.Sprintf("Final streak: %d", streak))
			fmt.Printf("%s %s\n", gameOver, finalStreak)
			keepPlaying = false
		}
	}
}
