package main

import (
	"fmt"
	"net"
)

const userFile = "users.json"

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
