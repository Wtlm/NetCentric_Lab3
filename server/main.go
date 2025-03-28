package main

import "fmt"

func main() {
	server, err := NewServer(8080)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	server.Start()
}
