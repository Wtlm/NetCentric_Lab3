package main

import (
	"fmt"
	"net"
	"sync"
)

type Word struct {
	Word        string `json:"word"`
	Description string `json:"description"`
	HiddenWord  []rune
}

type Game struct {
	Word  Word
	Score map[net.Conn]int
	Mutex sync.Mutex
}

const userFile = "users.json"
const wordFile = "words.json"

var gameMutex sync.Mutex
var userManager *AuthManager
var gameSessions = make(map[string]*GameSession)

func main() {
	var err error
	userManager, err = NewAuthManager(userFile)
	if err != nil {
		fmt.Println("Failed to load users:", err)
		return
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port 8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}
