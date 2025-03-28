package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func authenticateUser(conn net.Conn, reader *bufio.Reader) (string, bool) {
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Send authentication request
	fmt.Fprintf(conn, "%s %s\n", username, password)

	// Read server response
	authResponse, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print(authResponse)

	if !strings.HasPrefix(authResponse, "Authenticated") {
		return "", false
	}

	// Extract session key
	parts := strings.Split(authResponse, ": ")
	if len(parts) != 2 {
		fmt.Println("Invalid authentication response")
		return "", false
	}
	sessionKey := strings.TrimSpace(parts[1])

	return sessionKey, true
}

// startGuessingGame handles the guessing game loop
func startGuessingGame(conn net.Conn, sessionKey string, reader *bufio.Reader) {
	for {
		fmt.Print("Enter your guess (1-100): ")
		guess, _ := reader.ReadString('\n')
		guess = strings.TrimSpace(guess)

		// Send guess to server
		fmt.Fprintf(conn, "%s_GUESS:%s\n", sessionKey, guess)

		// Receive response
		serverResponse, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Server:", serverResponse)

		if strings.Contains(serverResponse, "CORRECT") {
			break
		}
	}
}

// sendMessage sends a regular message to the server
func sendMessage(conn net.Conn, sessionKey string, message string) {
	fmt.Fprintf(conn, "%s_%s\n", sessionKey, message)

	// Receive response
	serverResponse, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Server:", serverResponse)
}

func ReceiveFile(conn net.Conn, savePath string, sessionKey string) {
	reader := bufio.NewReader(conn)

	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	parts := strings.SplitN(response, "_", 2)

	if len(parts) < 2 || parts[0] != sessionKey {
		fmt.Println("Invalid response from server")
		return
	}

	if parts[1] != "READY" {
		fmt.Println("Server error:", parts[1])
		return
	}

	file, err := os.Create(savePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	fmt.Println("File content:")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving file:", err)
			break
		}
		line = strings.TrimSpace(line)

		parts := strings.SplitN(line, "_", 2)
		if len(parts) < 2 || parts[0] != sessionKey {
			fmt.Println("Invalid message format")
			continue
		}

		data := parts[1]
		if data == "EOF" {
			break
		}

		fmt.Println(data) // Print each line to the terminal
		writer.WriteString(data + "\n")
	}
	writer.Flush()
	fmt.Println("File downloaded successfully:", savePath)
}
