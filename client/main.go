package main

import "fmt"

func main() {
	client := NewGameClient()
	err := client.Connect("127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer client.Close()

	client.StartGame()
}
