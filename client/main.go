package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	// Authenticate User
	sessionKey, success := authenticateUser(conn, reader)
	if !success {
		return
	}

	go handleIncomingMessages(sessionKey, conn)
	// fmt.Print("Enter message (or 'quit' to exit)\n")
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "quit" {
			break
		}

		fmt.Fprintf(conn, "%s\n", message)
	}
}
