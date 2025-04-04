package main

import (
	"encoding/json"
	"math/rand"
	"net"
	"os"
	"time"
)

// Load words from file
func loadWords(filename string) ([]Word, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []Word
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&words); err != nil {
		return nil, err
	}
	return words, nil
}

// Create a new game session
func createGame(words []Word) *Game {
	rand.Seed(time.Now().UnixNano())
	selectedWord := words[rand.Intn(len(words))]

	// Initialize hidden word with underscores but preserve spaces
	hiddenWord := make([]rune, len(selectedWord.Word))
	for i, char := range selectedWord.Word {
		if char == ' ' {
			hiddenWord[i] = ' '
		} else {
			hiddenWord[i] = '_'
		}
	}

	return &Game{
		Word: Word{
			Word:        selectedWord.Word,
			Description: selectedWord.Description,
			HiddenWord:  hiddenWord,
		},
		Score: make(map[net.Conn]int),
	}
}
