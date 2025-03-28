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

	for {
		fmt.Print("Enter command (message/guess/exit/download): ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "exit" {
			break
		} else if command == "guess" {
			startGuessingGame(conn, sessionKey, reader)
		} else if strings.HasPrefix(command, "download ") {
			filename := strings.TrimPrefix(command, "download ")
			formattedRequest := fmt.Sprintf("%s_DOWNLOAD:%s\n", sessionKey, filename)
			fmt.Fprintf(conn, formattedRequest)
			ReceiveFile(conn, "received_"+filename, sessionKey)
		} else {
			sendMessage(conn, sessionKey, command)
		}
	}
}
