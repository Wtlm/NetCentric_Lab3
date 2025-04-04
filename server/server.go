package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"slices"
	"strconv"
	"strings"
	"time"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Receive authentication request
	authData, _ := reader.ReadString('\n')
	authData = strings.TrimSpace(authData)
	parts := strings.Split(authData, " ")
	if len(parts) != 2 {
		conn.Write([]byte("Invalid format\n"))
		return
	}

	username, password := parts[0], parts[1]
	sessionKey, err := userManager.Authenticate(username, password)
	if err != nil {
		conn.Write([]byte("Authentication failed\n"))
		return
	}

	// Send session key to client
	conn.Write([]byte(fmt.Sprintf("Authenticated. Session key: %d\n", sessionKey)))

	gameMutex.Lock()
	var session *GameSession

	// Look for an available session with space
	for _, s := range gameSessions {
		if len(s.players) < 4 { // Allow up to 4 players per session
			s.players = append(s.players, conn)
			session = s
			gameMutex.Unlock()
			broadcastMessage(session, fmt.Sprintf("Player %s joined! Players now: %d\n", username, len(session.players)))
			break
		}
	}
	if session == nil {
		gameID := strconv.Itoa(rand.Intn(10000))
		session = &GameSession{players: []net.Conn{conn}, turn: 0, gameStarted: false}
		gameSessions[gameID] = session
		gameMutex.Unlock()
		conn.Write([]byte(fmt.Sprintf("GameID: %s\n", gameID)))
	}
	if len(session.players) == 1 {
		broadcastMessage(session, "Waiting for more players...\n")
	} else {
		// Ask players whether to start or wait
		for {
			gameMutex.Lock()
			if session.gameStarted {
				gameMutex.Unlock()
				break
			}
			gameMutex.Unlock()

			conn.Write([]byte("Type 'start' to begin the game or 'wait' to wait for more players:\n"))
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(response)
			// parts := strings.Split(response, "_")
			// if len(parts) != 2 {
			// 	conn.Write([]byte("Invalid format\n"))
			// 	return
			// }

			// response = parts[1]

			conn.Write([]byte(fmt.Sprintf("%s\n", response)))
			if response == "start" {
				gameMutex.Lock()

				if !session.gameStarted {
					session.gameStarted = true
					// Load words and start the game
					words, err := loadWords(wordFile)
					if err != nil {
						conn.Write([]byte("Error loading words\n"))
						return
					}
					session.game = createGame(words)
					gameMutex.Unlock()

					broadcastMessage(session, "Game starting!\n")
					broadcastMessage(session, fmt.Sprintf("Game started! Description: %s\nThis word has %d letters: %s\n",
						session.game.Word.Description,
						len(strings.ReplaceAll(session.game.Word.Word, " ", "")),
						string(session.game.Word.HiddenWord)))

					for _, player := range session.players {
						go handleGame(session, player)
					}

				} else {
					gameMutex.Unlock()
				}
				break
			} else if response == "wait" {
				conn.Write([]byte("Waiting for more players...\n"))
				time.Sleep(3 * time.Second)
				break
			} else {
				conn.Write([]byte("Invalid input. Type 'start' or 'wait'.\n"))
			}
		}
	}
	// Keep the connection open
	select {} // Keeps the function running
}
func handleGame(session *GameSession, conn net.Conn) {
	for {
		session.mu.Lock()
		if len(session.players) == 0 {
			session.mu.Unlock()
			break // Exit if no players left
		}

		if session.Completed {
			session.mu.Unlock()
			return
		}

		activePlayer := session.players[session.turn]
		session.mu.Unlock()

		// Notify the current player it's their turn
		if activePlayer == conn {
			conn.Write([]byte("Your turn! Enter a letter:\n"))
			reader := bufio.NewReader(conn)
			message, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println("Client disconnected:", err)
				removePlayerFromSession(session, conn)
				return
			}

			message = strings.TrimSpace(message)
			message = strings.ToUpper(message)

			if message == "quit" {
				conn.Write([]byte("Goodbye!\n"))
				removePlayerFromSession(session, conn)
				return
			}

			if err := processGuess(message, session.game, conn, session); err != nil {
				conn.Write([]byte(err.Error()))
			}
			if session.Completed {
				broadcastMessage(session, fmt.Sprint("Game Over\n"))
				for numPlayer, player := range session.players {
					broadcastMessage(session, fmt.Sprintf("Score of Player %d: %d\n", numPlayer+1, session.game.Score[player]))
				}
				return
			}

			// Broadcast to all players who sent the message
			broadcastMessage(session, fmt.Sprintf("Player %d: %s\n", session.turn+1, message))

			// Move to the next player's turn
			session.mu.Lock()
			session.turn = (session.turn + 1) % len(session.players)
			session.mu.Unlock()
		}
	}
}

func processGuess(guess string, game *Game, conn net.Conn, session *GameSession) error {
	if guess == "" {
		return fmt.Errorf("Invalid guess.\n")
	}
	letter := rune(guess[0])
	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	// Check if the letter was already guessed
	if isLetterAlreadyGuessed(letter, game) {
		conn.Write([]byte(fmt.Sprintf("This letter is already guessed . Try another one.\nWord: %s\n", string(game.Word.HiddenWord))))
		return nil
	}

	found, count := false, 0
	for i, ch := range game.Word.Word {
		if ch == letter && game.Word.HiddenWord[i] == '_' {
			game.Word.HiddenWord[i] = ch
			found = true
			count++
		}
	}

	if found {
		game.Score[conn] += 10 * count
		conn.Write([]byte(fmt.Sprintf("Good guess! +%d points\n", 10*count)))
		broadcastMessage(session, fmt.Sprintf("Word: %s\n", string(game.Word.HiddenWord)))

		if !strings.ContainsRune(string(game.Word.HiddenWord), '_') {
			session.Completed = true
		}
	} else {
		conn.Write([]byte("Wrong guess! Waiting for your next turn...\n"))
		time.Sleep(3 * time.Second)
	}
	return nil
}

func broadcastMessage(session *GameSession, message string) {
	session.mu.Lock()
	defer session.mu.Unlock()
	for _, player := range session.players {
		player.Write([]byte(message))
	}
}

func removePlayerFromSession(session *GameSession, conn net.Conn) {
	session.mu.Lock()
	defer session.mu.Unlock()

	// Remove player from session
	for i, player := range session.players {
		if player == conn {
			session.players = append(session.players[:i], session.players[i+1:]...)
			break
		}
	}

	// If session is empty, remove it
	if len(session.players) == 0 {
		for gameID, s := range gameSessions {
			if s == session {
				delete(gameSessions, gameID)
				break
			}
		}
	} else {
		// If players are left, notify them
		broadcastMessage(session, "A player has left. The game continues.")
	}
}

// Check if the letter has already been guessed
func isLetterAlreadyGuessed(letter rune, game *Game) bool {
	return slices.Contains(game.Word.HiddenWord, letter)
}
