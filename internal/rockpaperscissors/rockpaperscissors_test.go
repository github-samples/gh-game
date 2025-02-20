package rockpaperscissors

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	tests := []struct {
		name            string
		bestOf          int
		wantBestOf      int
		wantGameOver    bool
		wantGamesPlayed int
	}{
		{
			name:            "New game with odd number of rounds",
			bestOf:          3,
			wantBestOf:      3,
			wantGameOver:    false,
			wantGamesPlayed: 0,
		},
		{
			name:            "New game with even number of rounds should increment to odd",
			bestOf:          4,
			wantBestOf:      5,
			wantGameOver:    false,
			wantGamesPlayed: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(tt.bestOf)
			if game.BestOf != tt.wantBestOf {
				t.Errorf("NewGame() BestOf = %v, want %v", game.BestOf, tt.wantBestOf)
			}
			if game.GameOver != tt.wantGameOver {
				t.Errorf("NewGame() GameOver = %v, want %v", game.GameOver, tt.wantGameOver)
			}
			if game.GamesPlayed != tt.wantGamesPlayed {
				t.Errorf("NewGame() GamesPlayed = %v, want %v", game.GamesPlayed, tt.wantGamesPlayed)
			}
		})
	}
}

func TestGame_getWinner(t *testing.T) {
	tests := []struct {
		name           string
		playerChoice   string
		computerChoice string
		want           string
	}{
		{
			name:           "Player wins with rock vs scissors",
			playerChoice:   "rock",
			computerChoice: "scissors",
			want:           "player",
		},
		{
			name:           "Player wins with paper vs rock",
			playerChoice:   "paper",
			computerChoice: "rock",
			want:           "player",
		},
		{
			name:           "Player wins with scissors vs paper",
			playerChoice:   "scissors",
			computerChoice: "paper",
			want:           "player",
		},
		{
			name:           "Player loses with rock vs paper",
			playerChoice:   "rock",
			computerChoice: "paper",
			want:           "computer",
		},
		{
			name:           "Player loses with paper vs scissors",
			playerChoice:   "paper",
			computerChoice: "scissors",
			want:           "computer",
		},
		{
			name:           "Player loses with scissors vs rock",
			playerChoice:   "scissors",
			computerChoice: "rock",
			want:           "computer",
		},
		{
			name:           "Draw with same choices (rock)",
			playerChoice:   "rock",
			computerChoice: "rock",
			want:           "draw",
		},
		{
			name:           "Draw with same choices (paper)",
			playerChoice:   "paper",
			computerChoice: "paper",
			want:           "draw",
		},
		{
			name:           "Draw with same choices (scissors)",
			playerChoice:   "scissors",
			computerChoice: "scissors",
			want:           "draw",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				PlayerChoice:   tt.playerChoice,
				ComputerChoice: tt.computerChoice,
			}
			if got := g.getWinner(); got != tt.want {
				t.Errorf("Game.getWinner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGame_Play(t *testing.T) {
	tests := []struct {
		name              string
		bestOf            int
		moves             []string
		wantGameOver      bool
		wantPlayerScore   int
		wantComputerScore int
	}{
		{
			name:              "Game ends when player chooses exit",
			bestOf:            3,
			moves:             []string{"exit"},
			wantGameOver:      true,
			wantPlayerScore:   0,
			wantComputerScore: 0,
		},
		{
			name:              "Game continues for valid moves",
			bestOf:            3,
			moves:             []string{"rock"},
			wantGameOver:      false,
			wantPlayerScore:   0, // Score will depend on random computer choice
			wantComputerScore: 0, // Score will depend on random computer choice
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame(tt.bestOf)
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
