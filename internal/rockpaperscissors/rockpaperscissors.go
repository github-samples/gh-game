package rockpaperscissors

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Game options
var (
	options = []string{"rock", "paper", "scissors", "exit"}
)

// Game represents a single game of Rock Paper Scissors.
type Game struct {
	// PlayerChoice is the choice made by the player.
	PlayerChoice string
	// ComputerChoice is the choice made by the computer.
	ComputerChoice string
	// Winner is the winner of the current round.
	Winner string
	// PlayerScore tracks the player's wins
	PlayerScore int
	// ComputerScore tracks the computer's wins
	ComputerScore int
	// BestOf determines how many games to play (e.g., best of 3, 5, 7)
	BestOf int
	// GameOver is a flag to indicate if the game is over.
	GameOver bool
	// GameOverMessage is the message to display when the game is over.
	GameOverMessage string
	// GamesPlayed tracks the number of games played
	GamesPlayed int
}

// Prompter defines an interface for getting user input
type Prompter interface {
	Select(prompt, defaultValue string, options []string) (int, error)
}

func NewGame(bestOf int) *Game {
	if bestOf%2 == 0 {
		bestOf++ // Ensure we have an odd number for "best of"
	}
	return &Game{
		PlayerChoice:    "",
		ComputerChoice:  "",
		Winner:          "",
		GamesPlayed:     0,
		PlayerScore:     0,
		ComputerScore:   0,
		BestOf:          bestOf,
		GameOver:        false,
		GameOverMessage: "",
	}
}

// Play plays a single round of Rock Paper Scissors.
func (g *Game) Play(playerChoice string) {
	if playerChoice == "exit" {
		g.GameOver = true
		g.GameOverMessage = "Game ended by player"
		return
	}

	g.PlayerChoice = playerChoice
	g.ComputerChoice = g.getComputerChoice()
	g.Winner = g.getWinner()
	g.updateScore()
	g.GamesPlayed++
	g.GameOver = g.isGameOver()
	if g.GameOver {
		g.GameOverMessage = g.getGameOverMessage()
	}
}

// getComputerChoice returns the choice made by the computer.
func (g *Game) getComputerChoice() string {
	rand.Seed(time.Now().UnixNano())
	// Only use the game options excluding "exit"
	choices := options[:len(options)-1]
	return choices[rand.Intn(len(choices))]
}

// getWinner returns the winner of the current round.
func (g *Game) getWinner() string {
	if g.PlayerChoice == g.ComputerChoice {
		return "draw"
	}
	if (g.PlayerChoice == "rock" && g.ComputerChoice == "scissors") ||
		(g.PlayerChoice == "paper" && g.ComputerChoice == "rock") ||
		(g.PlayerChoice == "scissors" && g.ComputerChoice == "paper") {
		return "player"
	}
	return "computer"
}

// updateScore updates the score based on the round winner
func (g *Game) updateScore() {
	if g.Winner == "player" {
		g.PlayerScore++
	} else if g.Winner == "computer" {
		g.ComputerScore++
	}
}

// isGameOver returns true if the game is over.
func (g *Game) isGameOver() bool {
	winsNeeded := (g.BestOf / 2) + 1

	// Game is over if:
	// 1. Either player has reached the required wins, or
	// 2. We've played all games in the series
	return g.PlayerScore >= winsNeeded ||
		g.ComputerScore >= winsNeeded ||
		g.GamesPlayed >= g.BestOf
}

// getGameOverMessage returns the message to display when the game is over.
func (g *Game) getGameOverMessage() string {
	if g.PlayerScore > g.ComputerScore {
		return fmt.Sprintf("GAME OVER: Player WINS (%d - %d)", g.PlayerScore, g.ComputerScore)
	} else if g.ComputerScore > g.PlayerScore {
		return fmt.Sprintf("GAME OVER: Player LOSES (%d - %d)", g.PlayerScore, g.ComputerScore)
	}
	return fmt.Sprintf("GAME OVER: DRAW (%d - %d)", g.PlayerScore, g.ComputerScore)
}

// getRoundResultMessage returns a concise message about the round result
func (g *Game) getRoundResultMessage() string {
	// Capitalize the first letter of choices for better display
	playerChoice := strings.Title(g.PlayerChoice)
	computerChoice := strings.Title(g.ComputerChoice)

	if g.Winner == "draw" {
		return fmt.Sprintf("Draw! Player (%s) - CPU (%s)", playerChoice, computerChoice)
	} else if g.Winner == "player" {
		return fmt.Sprintf("Player (%s) beats CPU (%s)", playerChoice, computerChoice)
	} else {
		return fmt.Sprintf("Player (%s) loses to CPU (%s)", playerChoice, computerChoice)
	}
}

// PlayGame plays a game of Rock Paper Scissors.
func PlayGame(prompter Prompter) {
	// Get the number of rounds from the user
	roundOptions := []string{"3", "5", "7", "9"}
	roundIndex, err := prompter.Select("How many rounds would you like to play (best of)?", "3", roundOptions)
	if err != nil {
		fmt.Printf("Error getting number of rounds: %v\n", err)
		return
	}
	bestOf := 3 // Default value
	if roundIndex >= 0 && roundIndex < len(roundOptions) {
		bestOf = parseInt(roundOptions[roundIndex])
	}

	game := NewGame(bestOf)
	fmt.Printf("Playing best of %d games\n", bestOf)

	for !game.GameOver {
		fmt.Printf("\nCurrent score - Player: %d, Computer: %d\n", game.PlayerScore, game.ComputerScore)

		// Get player choice using prompter
		playerChoiceIndex, err := prompter.Select("Choose your move", "rock", options)
		if err != nil {
			fmt.Printf("Error getting player choice: %v\n", err)
			return
		}
		playerChoice := options[playerChoiceIndex]

		game.Play(playerChoice)
		if playerChoice == "exit" {
			break
		}

		// Display a more concise round result
		fmt.Println(game.getRoundResultMessage())
	}
	fmt.Println(game.GameOverMessage)
}

// parseInt safely converts a string to an integer
func parseInt(s string) int {
	val := 3 // Default value
	fmt.Sscanf(s, "%d", &val)
	if val <= 0 {
		return 3
	}
	if val%2 == 0 {
		val++ // Convert even numbers to next odd number
	}
	return val
}
