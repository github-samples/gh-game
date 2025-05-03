// / Package wordguess implements a text-based Word Guess game where the player guesses
// a word letter by letter while trying to avoid too many incorrect guesses.
package wordguess

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	// MaxIncorrectGuesses is the maximum number of incorrect guesses allowed before losing
	MaxIncorrectGuesses = 6
)

var (
	// Stylized output for the game
	titleStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	wordStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	correctStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // green
	incorrectStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // red
	remainingStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15")) // white
	instructionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
)

// Game represents the state of a word guess game
type Game struct {
	Word             string   // The word to be guessed
	RevealedWord     string   // Current state of the word with guessed letters revealed
	GuessedLetters   []string // Letters that have been guessed
	IncorrectGuesses int      // Number of incorrect guesses
	IsOver           bool     // Whether the game is over
	HasWon           bool     // Whether the player has won
}

// WordList contains a selection of words for the game
var WordList = []string{
	"github", "actions", "workflow", "repository", "branch",
	"commit", "merge", "issues", "pull", "request", "codespace",
	"copilot", "project", "discussion", "milestone", "release",
	"clone", "fork", "gist", "markdown", "license", "readme",
}

// Prompter interface allows us to mock the prompt functionality in tests
type Prompter interface {
	Input(prompt string, defaultValue string) (string, error)
	Select(prompt string, defaultValue string, options []string) (int, error)
	Confirm(prompt string, defaultValue bool) (bool, error)
}

// NewGame creates and initializes a new Word Guess game
func NewGame() *Game {
	word := WordList[rand.Intn(len(WordList))]

	return &Game{
		Word:             strings.ToLower(word),
		RevealedWord:     strings.Repeat("_", len(word)),
		GuessedLetters:   []string{},
		IncorrectGuesses: 0,
		IsOver:           false,
		HasWon:           false,
	}
}

// GuessLetter processes a letter guess and updates the game state
func (g *Game) GuessLetter(letter string) error {
	// Convert to lowercase
	letter = strings.ToLower(letter)

	// Validate input
	if len(letter) != 1 || !isLetter(letter[0]) {
		return fmt.Errorf("please enter a single letter")
	}

	// Check if letter was already guessed
	for _, guessed := range g.GuessedLetters {
		if guessed == letter {
			return fmt.Errorf("you've already guessed '%s'", letter)
		}
	}

	// Add to guessed letters
	g.GuessedLetters = append(g.GuessedLetters, letter)

	// Check if letter is in the word
	if strings.Contains(g.Word, letter) {
		// Update revealed word
		newRevealedWord := []rune(g.RevealedWord)
		for i, char := range g.Word {
			if string(char) == letter {
				newRevealedWord[i] = char
			}
		}
		g.RevealedWord = string(newRevealedWord)

		// Check if the word is completely revealed (win condition)
		if !strings.Contains(g.RevealedWord, "_") {
			g.IsOver = true
			g.HasWon = true
		}
	} else {
		// Incorrect guess
		g.IncorrectGuesses++

		// Check if max incorrect guesses reached (lose condition)
		if g.IncorrectGuesses >= MaxIncorrectGuesses {
			g.IsOver = true
		}
	}

	return nil
}

// isLetter checks if a byte is a letter
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// GetRemainingLetters returns a string of letters that haven't been guessed yet
func (g *Game) GetRemainingLetters() string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	var remaining strings.Builder

	for _, char := range alphabet {
		found := false
		for _, guessed := range g.GuessedLetters {
			if string(char) == guessed {
				found = true
				break
			}
		}

		if !found {
			remaining.WriteRune(char)
		}
	}

	return remaining.String()
}

// String returns a string representation of the current game state
func (g *Game) String() string {
	var sb strings.Builder

	// Display title
	sb.WriteString(titleStyle.Render("W O R D  G U E S S") + "\n\n")

	// Display the remaining guesses prominently
	incorrectLeft := MaxIncorrectGuesses - g.IncorrectGuesses
	guessesDisplay := fmt.Sprintf("Guesses Remaining: %d/%d", incorrectLeft, MaxIncorrectGuesses)

	if incorrectLeft > 3 {
		sb.WriteString(correctStyle.Render(guessesDisplay))
	} else if incorrectLeft > 1 {
		sb.WriteString(instructionStyle.Render(guessesDisplay))
	} else {
		sb.WriteString(incorrectStyle.Render(guessesDisplay))
	}
	sb.WriteString("\n\n")

	// Display the word with guessed letters
	displayWord := ""
	for _, char := range g.RevealedWord {
		if char == '_' {
			displayWord += "_ "
		} else {
			displayWord += string(char) + " "
		}
	}
	sb.WriteString(wordStyle.Render(displayWord) + "\n\n")

	// Display guessed letters
	sb.WriteString("Guessed: ")
	for _, letter := range g.GuessedLetters {
		if strings.Contains(g.Word, letter) {
			sb.WriteString(correctStyle.Render(letter + " "))
		} else {
			sb.WriteString(incorrectStyle.Render(letter + " "))
		}
	}
	sb.WriteString("\n\n")

	// Display remaining letters
	sb.WriteString("Available: ")
	sb.WriteString(remainingStyle.Render(strings.Join(strings.Split(g.GetRemainingLetters(), ""), " ")))
	sb.WriteString("\n\n")

	// Display game status
	if g.IsOver {
		if g.HasWon {
			sb.WriteString(correctStyle.Render("🎉 Congratulations! You guessed the word!"))
		} else {
			sb.WriteString(incorrectStyle.Render("😔 Game over! The word was: ") +
				wordStyle.Render(g.Word))
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString(instructionStyle.Render("Guess a letter to continue.\n"))
	}

	return sb.String()
}

// PlayGame starts a word guessing game session with the provided prompter
func PlayGame(p Prompter) {
	game := NewGame()

	fmt.Println(titleStyle.Render("\nWelcome to Word Guess!"))
	fmt.Println(instructionStyle.Render("Guess the GitHub-related term one letter at a time."))
	fmt.Println()

	// Main game loop
	for !game.IsOver {
		fmt.Println(game)

		// Get player's guess
		guess, err := p.Input("Enter a letter: ", "")
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		err = game.GuessLetter(guess)
		if err != nil {
			fmt.Println(incorrectStyle.Render(err.Error()))
		}
	}

	// Show final state
	fmt.Println(game)

	// Ask to play again
	playAgain, err := p.Confirm("Play again?", true)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	if playAgain {
		PlayGame(p)
	} else {
		fmt.Println(titleStyle.Render("Thanks for playing Word Guess!"))
	}
}
