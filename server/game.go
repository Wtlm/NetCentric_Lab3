package main

import (
	"crypto/rand"
	"math/big"
)

// Game represents a single guessing game instance
type Game struct {
	targetNumber int
	attempts     int
}

// NewGame creates a new game instance
func NewGame() Game {
	// Cryptographically secure random number generation
	n, _ := rand.Int(rand.Reader, big.NewInt(100))
	return Game{
		targetNumber: int(n.Int64()) + 1,
		attempts:     0,
	}
}

// ProcessGuess evaluates a player's guess
func (g *Game) ProcessGuess(guess int) (string, bool) {
	g.attempts++

	if guess < g.targetNumber {
		return "TOO_LOW", false
	} else if guess > g.targetNumber {
		return "TOO_HIGH", false
	} else {
		return "CORRECT", true
	}
}

// GetAttempts returns the number of attempts made
func (g *Game) GetAttempts() int {
	return g.attempts
}

// Reset resets the game to a new state
func (g *Game) Reset() {
	n, _ := rand.Int(rand.Reader, big.NewInt(100))
	g.targetNumber = int(n.Int64()) + 1
	g.attempts = 0
}
