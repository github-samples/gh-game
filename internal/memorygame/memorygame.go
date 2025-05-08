package memorygame

import (
	"math/rand"
	"time"
)

// Color represents a color in the sequence
type Color string

const (
	Red    Color = "red"
	Yellow Color = "yellow"
	Green  Color = "green"
	Blue   Color = "blue"
)

// Game represents the memory game state
type Game struct {
	Lives           int
	CurrentRound    int
	MaxRound        int
	Sequence        []Color
	DisplayTime     time.Duration
	AvailableColors []Color
}

// NewGame creates a new memory game instance
func NewGame(lives int) *Game {
	return &Game{
		Lives:           lives,
		CurrentRound:    1,
		MaxRound:        100,
		DisplayTime:     3 * time.Second,
		AvailableColors: []Color{Red, Yellow, Green, Blue},
	}
}

// GenerateSequence creates a new sequence for the current round
func (g *Game) GenerateSequence() {
	sequenceLength := g.CurrentRound + 2 // Start with 3 colors in round 1
	g.Sequence = make([]Color, sequenceLength)

	for i := 0; i < sequenceLength; i++ {
		g.Sequence[i] = g.AvailableColors[rand.Intn(len(g.AvailableColors))]
	}
}

// CheckSequence validates the user's input against the current sequence
func (g *Game) CheckSequence(userSequence []Color) bool {
	if len(userSequence) != len(g.Sequence) {
		return false
	}

	for i := range g.Sequence {
		if userSequence[i] != g.Sequence[i] {
			return false
		}
	}

	return true
}

// DecrementLives reduces the number of lives by 1
func (g *Game) DecrementLives() {
	g.Lives--
}

// NextRound advances to the next round
func (g *Game) NextRound() {
	g.CurrentRound++
}

// IsGameOver checks if the game is over
func (g *Game) IsGameOver() bool {
	return g.Lives <= 0 || g.CurrentRound > g.MaxRound
}
