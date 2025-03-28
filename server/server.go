package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
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

	// Initialize a new game session for the user
	gameMutex.Lock()
	games[username] = &Game{}
	games[username].Reset()
	gameMutex.Unlock()

	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		msgParts := strings.SplitN(message, "_", 2)

		if len(msgParts) != 2 {
			conn.Write([]byte("Invalid message format\n"))
			continue
		}

		clientKey, msgContent := msgParts[0], msgParts[1]
		sessionKeyInt := 0
		fmt.Sscanf(clientKey, "%d", &sessionKeyInt)

		if !userManager.ValidateSession(username, sessionKeyInt) {
			conn.Write([]byte("Invalid session key\n"))
			continue
		}

		if strings.HasPrefix(msgContent, "DOWNLOAD:") {
			filename := strings.TrimPrefix(msgContent, "DOWNLOAD:")
			fmt.Println("Client requested file:", filename)
			SendFile(conn, filename, sessionKeyInt)
			continue
		}
		if strings.HasPrefix(msgContent, "GUESS:") {
			guessStr := strings.TrimPrefix(msgContent, "GUESS:")
			guess, err := strconv.Atoi(guessStr)
			if err != nil {
				conn.Write([]byte("Invalid guess format\n"))
				continue
			}

			// Process the guess
			gameMutex.Lock()
			game := games[username]
			response, correct := game.ProcessGuess(guess)
			gameMutex.Unlock()

			conn.Write([]byte(fmt.Sprintf("%d_%s\n", sessionKeyInt, response)))

			if correct {
				conn.Write([]byte(fmt.Sprintf("%d_Correct! You guessed in %d attempts.\n", sessionKeyInt, game.GetAttempts())))
				gameMutex.Lock()
				game.Reset()
				gameMutex.Unlock()
			}
		} else {
			fmt.Printf("Received from %s: %s\n", username, msgContent)
			conn.Write([]byte(fmt.Sprintf("%d_Echo: %s\n", sessionKeyInt, msgContent)))
		}
	}
}

func SendFile(conn net.Conn, filename string, sessionKey int) {
	file, err := os.Open(filename)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("%d_ERROR: File not found\n", sessionKey)))
		return
	}
	defer file.Close()

	// Notify the client that the file transfer is starting
	conn.Write([]byte(fmt.Sprintf("%d_READY\n", sessionKey)))

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n') // Read line-by-line
		if err == io.EOF {
			break
		}
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("%d_ERROR: Read error\n", sessionKey)))
			return
		}

		// Send each line prefixed with sessionKey
		conn.Write([]byte(fmt.Sprintf("%d_%s", sessionKey, line)))
	}

	// Send EOF signal to indicate end of file transfer
	conn.Write([]byte(fmt.Sprintf("%d_EOF\n", sessionKey)))
	fmt.Println("File sent successfully:", filename)
}
