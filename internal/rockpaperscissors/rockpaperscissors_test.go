package rockpaperscissors

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewGame(t *testing.T) {
	tests := []struct {
		name            string
		bestOf          int
		secretMode      bool
		wantBestOf      int
		wantGameOver    bool
		wantGamesPlayed int
		wantSecretMode  bool
	}{
		{
			name:            "New game with odd number of rounds",
			bestOf:          3,
			secretMode:      false,
			wantBestOf:      3,
			wantGameOver:    false,
			wantGamesPlayed: 0,
			wantSecretMode:  false,
		},
		{
			name:            "New game with even number of rounds should increment to odd",
			bestOf:          4,
			secretMode:      false,
			wantBestOf:      5,
			wantGameOver:    false,
			wantGamesPlayed: 0,
			wantSecretMode:  false,
		},
		{
			name:            "New game with secret mode enabled",
			bestOf:          3,
			secretMode:      true,
			wantBestOf:      3,
			wantGameOver:    false,
			wantGamesPlayed: 0,
			wantSecretMode:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(tt.bestOf, tt.secretMode)
			if game.BestOf != tt.wantBestOf {
				t.Errorf("NewGame() BestOf = %v, want %v", game.BestOf, tt.wantBestOf)
			}
			if game.GameOver != tt.wantGameOver {
				t.Errorf("NewGame() GameOver = %v, want %v", game.GameOver, tt.wantGameOver)
			}
			if game.GamesPlayed != tt.wantGamesPlayed {
				t.Errorf("NewGame() GamesPlayed = %v, want %v", game.GamesPlayed, tt.wantGamesPlayed)
				if game.SecretMode != tt.wantSecretMode {
					t.Errorf("NewGame() SecretMode = %v, want %v", game.SecretMode, tt.wantSecretMode)
				}
			}
		})
	}
}

// gameCombination represents a specific game combination for testing.
// It defines a test case with player and computer choices, game mode,
// and the expected winner for that combination.
type gameCombination struct {
	name           string // Name of the test case
	playerChoice   string // Player's choice (rock, paper, scissors, etc.)
	computerChoice string // Computer's choice (rock, paper, scissors, etc.)
	secretMode     bool   // Whether secret mode is enabled
	winner         string // Expected winner (player, computer, or tie)
}

// getGameCombinations returns a comprehensive list of game combinations for testing
// different Rock Paper Scissors outcomes in both normal and secret modes.
// It covers all winning, losing, and tie scenarios.
func getGameCombinations() []gameCombination {
	return []gameCombination{
		// Normal mode combinations
		{
			name:           "Normal mode - Player wins - rock beats scissors",
			playerChoice:   "rock",
			computerChoice: "scissors",
			secretMode:     false,
			winner:         "player",
		},
		{
			name:           "Normal mode - Player wins - paper beats rock",
			playerChoice:   "paper",
			computerChoice: "rock",
			secretMode:     false,
			winner:         "player",
		},
		{
			name:           "Normal mode - Player wins - scissors beats paper",
			playerChoice:   "scissors",
			computerChoice: "paper",
			secretMode:     false,
			winner:         "player",
		},
		{
			name:           "Normal mode - Player loses - paper beats rock",
			playerChoice:   "rock",
			computerChoice: "paper",
			secretMode:     false,
			winner:         "computer",
		},
		{
			name:           "Normal mode - Player loses - scissors beats paper",
			playerChoice:   "paper",
			computerChoice: "scissors",
			secretMode:     false,
			winner:         "computer",
		},
		{
			name:           "Normal mode - Player loses - rock beats scissors",
			playerChoice:   "scissors",
			computerChoice: "rock",
			secretMode:     false,
			winner:         "computer",
		},
		{
			name:           "Normal mode - Draw with same choices (rock)",
			playerChoice:   "rock",
			computerChoice: "rock",
			secretMode:     false,
			winner:         "draw",
		},
		{
			name:           "Normal mode - Draw with same choices (paper)",
			playerChoice:   "paper",
			computerChoice: "paper",
			secretMode:     false,
			winner:         "draw",
		},
		{
			name:           "Normal mode - Draw with same choices (scissors)",
			playerChoice:   "scissors",
			computerChoice: "scissors",
			secretMode:     false,
			winner:         "draw",
		},
		// Secret mode combinations
		{
			name:           "Secret mode - Player wins - rock beats lizard",
			playerChoice:   "rock",
			computerChoice: "lizard",
			secretMode:     true,
			winner:         "player",
		},
		{
			name:           "Secret mode - Player wins - lizard beats spock",
			playerChoice:   "lizard",
			computerChoice: "spock",
			secretMode:     true,
			winner:         "player",
		},
		{
			name:           "Secret mode - Player wins - spock beats scissors",
			playerChoice:   "spock",
			computerChoice: "scissors",
			secretMode:     true,
			winner:         "player",
		},
		{
			name:           "Secret mode - Player wins - scissors beats lizard",
			playerChoice:   "scissors",
			computerChoice: "lizard",
			secretMode:     true,
			winner:         "player",
		},
		{
			name:           "Secret mode - Player wins - lizard beats paper",
			playerChoice:   "lizard",
			computerChoice: "paper",
			secretMode:     true,
			winner:         "player",
		},
		{
			name:           "Secret mode - Player wins - paper beats spock",
			playerChoice:   "paper",
			computerChoice: "spock",
			secretMode:     true,
			winner:         "player",
		},
		{
			name:           "Secret mode - Player wins - spock beats rock",
			playerChoice:   "spock",
			computerChoice: "rock",
			secretMode:     true,
			winner:         "player",
		},
		{
			name:           "Secret mode - Draw - lizard vs lizard",
			playerChoice:   "lizard",
			computerChoice: "lizard",
			secretMode:     true,
			winner:         "draw",
		},
		{
			name:           "Secret mode - Draw - spock vs spock",
			playerChoice:   "spock",
			computerChoice: "spock",
			secretMode:     true,
			winner:         "draw",
		},
		{
			name:           "Secret mode - Player loses - lizard beats rock",
			playerChoice:   "lizard",
			computerChoice: "rock",
			secretMode:     true,
			winner:         "computer",
		},
		{
			name:           "Secret mode - Player loses - lizard beats spock",
			playerChoice:   "spock",
			computerChoice: "lizard",
			secretMode:     true,
			winner:         "computer",
		},
		{
			name:           "Secret mode - Player loses - spock beats scissors",
			playerChoice:   "scissors",
			computerChoice: "spock",
			secretMode:     true,
			winner:         "computer",
		},
		{
			name:           "Secret mode - Player loses - scissors beats lizard",
			playerChoice:   "lizard",
			computerChoice: "scissors",
			secretMode:     true,
			winner:         "computer",
		},
		{
			name:           "Secret mode - Player loses - lizard beats paper",
			playerChoice:   "paper",
			computerChoice: "lizard",
			secretMode:     true,
			winner:         "computer",
		},
		{
			name:           "Secret mode - Player loses - paper beats spock",
			playerChoice:   "spock",
			computerChoice: "paper",
			secretMode:     true,
			winner:         "computer",
		},
		{
			name:           "Secret mode - Player loses - spock beats rock",
			playerChoice:   "rock",
			computerChoice: "spock",
			secretMode:     true,
			winner:         "computer",
		},
	}
}

func TestGame_getWinner(t *testing.T) {
	combinations := getGameCombinations()

	for _, tt := range combinations {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				PlayerChoice:   tt.playerChoice,
				ComputerChoice: tt.computerChoice,
				SecretMode:     tt.secretMode,
			}
			if got := g.getWinner(); got != tt.winner {
				t.Errorf("Game.getWinner() = %v, want %v", got, tt.winner)
			}
		})
	}
}

func TestGame_Play(t *testing.T) {
	tests := []struct {
		name              string
		bestOf            int
		moves             []string
		secretMode        bool
		wantGameOver      bool
		wantPlayerScore   int
		wantComputerScore int
	}{
		{
			name:              "Game ends when player chooses exit",
			bestOf:            3,
			moves:             []string{"exit"},
			secretMode:        false,
			wantGameOver:      true,
			wantPlayerScore:   0,
			wantComputerScore: 0,
		},
		{
			name:              "Game continues for valid moves",
			bestOf:            3,
			moves:             []string{"rock"},
			secretMode:        false,
			wantGameOver:      false,
			wantPlayerScore:   0, // Score will depend on random computer choice
			wantComputerScore: 0, // Score will depend on random computer choice
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame(tt.bestOf, tt.secretMode)
			for _, move := range tt.moves {
				g.Play(move)
			}
			if g.GameOver != tt.wantGameOver {
				t.Errorf("Game.Play() GameOver = %v, want %v", g.GameOver, tt.wantGameOver)
			}
			if tt.moves[0] == "exit" {
				if g.GameOverMessage != "Game ended by player" {
					t.Errorf("Game.Play() GameOverMessage = %v, want 'Game ended by player'", g.GameOverMessage)
				}
			}
		})
	}
}

func TestGame_isGameOver(t *testing.T) {
	tests := []struct {
		name          string
		bestOf        int
		playerScore   int
		computerScore int
		gamesPlayed   int
		want          bool
	}{
		{
			name:          "Game not over - early game",
			bestOf:        3,
			playerScore:   0,
			computerScore: 0,
			gamesPlayed:   1,
			want:          false,
		},
		{
			name:          "Game over - player wins best of 3",
			bestOf:        3,
			playerScore:   2,
			computerScore: 0,
			gamesPlayed:   2,
			want:          true,
		},
		{
			name:          "Game over - computer wins best of 3",
			bestOf:        3,
			playerScore:   0,
			computerScore: 2,
			gamesPlayed:   2,
			want:          true,
		},
		{
			name:          "Game over - all games played",
			bestOf:        3,
			playerScore:   1,
			computerScore: 1,
			gamesPlayed:   3,
			want:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				BestOf:        tt.bestOf,
				PlayerScore:   tt.playerScore,
				ComputerScore: tt.computerScore,
				GamesPlayed:   tt.gamesPlayed,
			}
			if got := g.isGameOver(); got != tt.want {
				t.Errorf("Game.isGameOver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGame_updateScore(t *testing.T) {
	tests := []struct {
		name              string
		winner            string
		wantPlayerScore   int
		wantComputerScore int
	}{
		{
			name:              "Player wins",
			winner:            "player",
			wantPlayerScore:   1,
			wantComputerScore: 0,
		},
		{
			name:              "Computer wins",
			winner:            "computer",
			wantPlayerScore:   0,
			wantComputerScore: 1,
		},
		{
			name:              "Draw",
			winner:            "draw",
			wantPlayerScore:   0,
			wantComputerScore: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Winner: tt.winner,
			}
			g.updateScore()
			if g.PlayerScore != tt.wantPlayerScore {
				t.Errorf("updateScore() PlayerScore = %v, want %v", g.PlayerScore, tt.wantPlayerScore)
			}
			if g.ComputerScore != tt.wantComputerScore {
				t.Errorf("updateScore() ComputerScore = %v, want %v", g.ComputerScore, tt.wantComputerScore)
			}
		})
	}
}

func TestGame_getGameOverMessage(t *testing.T) {
	tests := []struct {
		name            string
		playerScore     int
		computerScore   int
		wantMsgContains string
	}{
		{
			name:            "Player wins",
			playerScore:     2,
			computerScore:   1,
			wantMsgContains: "GAME OVER: Player WINS",
		},
		{
			name:            "Computer wins",
			playerScore:     1,
			computerScore:   2,
			wantMsgContains: "GAME OVER: Player LOSES",
		},
		{
			name:            "Draw",
			playerScore:     1,
			computerScore:   1,
			wantMsgContains: "GAME OVER: DRAW",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				PlayerScore:   tt.playerScore,
				ComputerScore: tt.computerScore,
			}
			got := g.getGameOverMessage()
			if got == "" {
				t.Error("getGameOverMessage() returned empty string")
			}
			if !strings.Contains(got, tt.wantMsgContains) {
				t.Errorf("getGameOverMessage() = %v, want it to contain %v", got, tt.wantMsgContains)
			}
		})
	}
}

func TestGame_getRoundResultMessage(t *testing.T) {
	combinations := getGameCombinations()

	for _, tt := range combinations {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				PlayerChoice:   tt.playerChoice,
				ComputerChoice: tt.computerChoice,
				Winner:         tt.winner,
			}

			got := g.getRoundResultMessage()

			// Verify the message contains expected substrings based on the winner
			wantContains := []string{tt.playerChoice, tt.computerChoice}

			switch tt.winner {
			case "player":
				wantContains = append(wantContains, "Player", "beats")
			case "computer":
				wantContains = append(wantContains, "loses to")
			case "draw":
				wantContains = append(wantContains, "Draw", "Player", "CPU")
			}

			for _, want := range wantContains {
				if !strings.Contains(strings.ToLower(got), strings.ToLower(want)) {
					t.Errorf("getRoundResultMessage() = %v, want it to contain %v", got, want)
				}
			}
		})
	}
}

// MockPrompter implements the Prompter interface for testing
type MockPrompter struct {
	selectReturn int
	selectError  error
}

func (m *MockPrompter) Select(prompt, defaultValue string, options []string) (int, error) {
	return m.selectReturn, m.selectError
}

type mockPromptSequence struct {
	returns []int
	errors  []error
	index   int
}

func (m *mockPromptSequence) Select(prompt, defaultValue string, options []string) (int, error) {
	if m.index >= len(m.returns) {
		return 0, nil
	}
	ret := m.returns[m.index]
	err := m.errors[m.index]
	m.index++
	return ret, err
}

func TestPlayGame(t *testing.T) {
	tests := []struct {
		name       string
		prompter   Prompter
		secretMode bool
	}{
		{
			name: "Complete game sequence - standard mode",
			prompter: &mockPromptSequence{
				returns: []int{1, 0, 0, 0}, // Select 5 rounds, then rock three times
				errors:  []error{nil, nil, nil, nil},
			},
			secretMode: false,
		},
		{
			name: "Complete game sequence - secret mode",
			prompter: &mockPromptSequence{
				returns: []int{1, 0, 0, 0}, // Select 5 rounds, then rock three times
				errors:  []error{nil, nil, nil, nil},
			},
			secretMode: true,
		},
		{
			name: "Error on rounds selection",
			prompter: &mockPromptSequence{
				returns: []int{0},
				errors:  []error{fmt.Errorf("mock error")},
			},
			secretMode: false,
		},
		{
			name: "Error on move selection",
			prompter: &mockPromptSequence{
				returns: []int{0, 0},
				errors:  []error{nil, fmt.Errorf("mock error")},
			},
			secretMode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PlayGame(tt.prompter, tt.secretMode)
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "Valid number",
			input: "5",
			want:  5,
		},
		{
			name:  "Invalid input returns default",
			input: "invalid",
			want:  3,
		},
		{
			name:  "Empty input returns default",
			input: "",
			want:  3,
		},
		{
			name:  "Negative number returns default",
			input: "-1",
			want:  3,
		},
		{
			name:  "Zero returns default",
			input: "0",
			want:  3,
		},
		{
			name:  "Even number returns next odd number",
			input: "4",
			want:  5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInt(tt.input); got != tt.want {
				t.Errorf("parseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGame_getComputerChoice(t *testing.T) {
	// Test standard mode
	g1 := &Game{
		SecretMode: false,
	}

	standardChoices := map[string]bool{
		"rock":     false,
		"paper":    false,
		"scissors": false,
	}

	// Run multiple times to ensure we get a variety of responses
	for i := 0; i < 100; i++ {
		choice := g1.getComputerChoice()
		if _, ok := standardChoices[choice]; !ok {
			t.Errorf("getComputerChoice() in standard mode returned an invalid choice: %v", choice)
		}
		standardChoices[choice] = true
	}

	// Ensure we got all possible standard choices
	for choice, found := range standardChoices {
		if !found {
			t.Errorf("getComputerChoice() in standard mode never returned %v", choice)
		}
	}

	// Test secret mode
	g2 := &Game{
		SecretMode: true,
	}

	secretChoices := map[string]bool{
		"rock":     false,
		"paper":    false,
		"scissors": false,
		"lizard":   false,
		"spock":    false,
	}

	// Run multiple times to ensure we get a variety of responses
	for i := 0; i < 200; i++ {
		choice := g2.getComputerChoice()
		if _, ok := secretChoices[choice]; !ok {
			t.Errorf("getComputerChoice() in secret mode returned an invalid choice: %v", choice)
		}
		secretChoices[choice] = true
	}
}

func TestPlayGame_EnhancedCoverage(t *testing.T) {
	tests := []struct {
		name       string
		returns    []int
		errors     []error
		secretMode bool
	}{
		{
			name:       "Exit early",
			returns:    []int{0, 3}, // Select 3 rounds, then exit
			errors:     []error{nil, nil},
			secretMode: false,
		},
		{
			name:       "Invalid round index",
			returns:    []int{99, 0, 0}, // Invalid round index should use default
			errors:     []error{nil, nil, nil},
			secretMode: false,
		},
		{
			name:       "Secret mode with exit",
			returns:    []int{2, 3}, // Select 7 rounds, then exit
			errors:     []error{nil, nil},
			secretMode: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPrompt := &mockPromptSequence{
				returns: tt.returns,
				errors:  tt.errors,
			}

			PlayGame(mockPrompt, tt.secretMode)
		})
	}
}
