package main

import (
	"bufio"
	"fmt"
	"net"
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


