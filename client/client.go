package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type GameClient struct {
	conn     net.Conn
	token    string
	username string
	scanner  *bufio.Scanner
}

func NewGameClient() *GameClient {
	return &GameClient{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (c *GameClient) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.conn = conn

	// Authenticate the user
	if err := c.authenticate(); err != nil {
		return err
	}

	return nil
}

func (c *GameClient) authenticate() error {
	// Handle username request
	prompt, err := c.readServerResponse()
	if err != nil || prompt != "USERNAME_REQ" {
		return fmt.Errorf("unexpected server response during username request")
	}
	fmt.Print("Enter username: ")
	username := c.readConsoleInput()
	c.conn.Write([]byte(username + "\n"))

	// Handle password request
	prompt, err = c.readServerResponse()
	if err != nil || prompt != "PASSWORD_REQ" {
		return fmt.Errorf("unexpected server response during password request")
	}
	fmt.Print("Enter password: ")
	password := c.readConsoleInput()
	c.conn.Write([]byte(password + "\n"))

	// Read authentication result
	response, err := c.readServerResponse()
	if err != nil {
		return err
	}

	// Check for token or error
	if strings.HasPrefix(response, "TOKEN:") {
		c.token = strings.TrimPrefix(response, "TOKEN:")
		c.username = username
		fmt.Println("Authentication successful!")
		return nil
	} else if strings.HasPrefix(response, "ERR_AUTH:") {
		return fmt.Errorf("authentication failed: %s", strings.TrimPrefix(response, "ERR_AUTH:"))
	}

	return fmt.Errorf("unexpected server response")
}

func (c *GameClient) StartGame() {
	fmt.Println("Welcome to the Guessing Game!")
	fmt.Println("Try to guess a number between 1 and 100")

	for {
		fmt.Print("Enter your guess (or 'quit' to exit): ")
		input := c.readConsoleInput()

		if strings.ToLower(input) == "quit" {
			c.sendMessage("QUIT")
			break
		}

		// Send guess with token prefix
		c.sendMessage(fmt.Sprintf("GUESS:%s", input))

		// Read and process server response
		c.processGameResponse()
	}
}

func (c *GameClient) sendMessage(message string) {
	// Prefix message with token
	fullMessage := fmt.Sprintf("%s_%s\n", c.token, message)
	c.conn.Write([]byte(fullMessage))
}

func (c *GameClient) processGameResponse() {
	response, err := c.readServerResponse()
	if err != nil {
		fmt.Println("Error reading server response:", err)
		return
	}

	// Remove token prefix
	parts := strings.SplitN(response, "_", 2)
	if len(parts) < 2 {
		fmt.Println("Invalid server response")
		return
	}
	message := parts[1]

	switch {
	case message == "TOO_LOW_GUESS":
		fmt.Println("Too low! Try a higher number.")
	case message == "TOO_HIGH_GUESS":
		fmt.Println("Too high! Try a lower number.")
	case strings.HasPrefix(message, "CORRECT_GUESS_ATTEMPTS:"):
		attempts := strings.TrimPrefix(message, "CORRECT_GUESS_ATTEMPTS:")
		fmt.Printf("Congratulations! You guessed the correct number in %s attempts!\n", attempts)
	case message == "NEW_GAME_STARTED":
		fmt.Println("A new game has started.")
	case message == "GOODBYE":
		fmt.Println("Goodbye! Thank you for playing.")
		c.Close()
		os.Exit(0) // Safe exit after closing connection
	case strings.HasPrefix(message, "ERR_"):
		fmt.Println("Error:", message)
	default:
		fmt.Println("Server response:", message)
	}
}

func (c *GameClient) readConsoleInput() string {
	c.scanner.Scan()
	return strings.TrimSpace(c.scanner.Text())
}

func (c *GameClient) readServerResponse() (string, error) {
	buffer := make([]byte, 1024)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buffer[:n])), nil
}

func (c *GameClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
