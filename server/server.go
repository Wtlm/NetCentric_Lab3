package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	listener       net.Listener
	userManager    *UserManager
	sessionManager *SessionManager
}

// NewServer initializes a TCP server
func NewServer(port int) (*Server, error) {
	userManager := NewUserManager()
	if err := userManager.LoadUsers("users.json"); err != nil {
		return nil, fmt.Errorf("failed to load users: %v", err)
	}

	// Start TCP listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener:       listener,
		userManager:    userManager,
		sessionManager: NewSessionManager(),
	}

	return server, nil
}

// Start listens for incoming TCP connections
func (s *Server) Start() {
	fmt.Println("Server listening on TCP port 8080...")

	// Session cleanup every 15 minutes
	go func() {
		for {
			time.Sleep(15 * time.Minute)
			s.sessionManager.CleanupExpiredSessions(30 * time.Minute)
		}
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("New client connected:", conn.RemoteAddr())
		go s.handleConnection(conn) // Handle client in a goroutine
	}
}

// Handle client authentication and game session
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	username, err := s.authenticateClient(conn)
	if err != nil {
		conn.Write([]byte("ERR_AUTH: " + err.Error() + "\n"))
		return
	}

	session := s.sessionManager.CreateSession(username)
	conn.Write([]byte(fmt.Sprintf("TOKEN:%s\n", session.Token)))

	s.handleGameCommands(conn, session)
}

// Authenticate user
func (s *Server) authenticateClient(conn net.Conn) (string, error) {
	conn.Write([]byte("USERNAME_REQ\n"))
	username, err := readClientInput(conn)
	if err != nil {
		return "", err
	}

	conn.Write([]byte("PASSWORD_REQ\n"))
	password, err := readClientInput(conn)
	if err != nil {
		return "", err
	}

	if !s.userManager.Authenticate(username, password) {
		return "", fmt.Errorf("invalid credentials")
	}

	return username, nil
}

// Handle incoming game commands
func (s *Server) handleGameCommands(conn net.Conn, session *ClientSession) {
	for {
		input, err := readClientInput(conn)
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			s.sessionManager.RemoveSession(session.Token) // Remove session on disconnect
			return
		}

		parts := strings.SplitN(input, "_", 2)
		if len(parts) != 2 || parts[0] != session.Token {
			conn.Write([]byte("ERR_TOKEN: Invalid token\n"))
			continue
		}

		command := parts[1]
		switch {
		case strings.HasPrefix(command, "GUESS:"):
			s.handleGuess(conn, session, command)
		case command == "NEW_GAME":
			s.startNewGame(conn, session)
		case command == "QUIT":
			conn.Write([]byte(fmt.Sprintf("%s_GOODBYE\n", session.Token)))
			s.sessionManager.RemoveSession(session.Token)
			return
		default:
			conn.Write([]byte(fmt.Sprintf("%s_ERR_UNKNOWN_COMMAND\n", session.Token)))
		}

		session.LastActive = time.Now()
	}
}

// Process a player's guess
func (s *Server) handleGuess(conn net.Conn, session *ClientSession, command string) {
	guessStr := strings.TrimPrefix(command, "GUESS:")
	guess, err := strconv.Atoi(guessStr)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("%s_ERR_INVALID_GUESS\n", session.Token)))
		return
	}

	result, gameOver := session.Game.ProcessGuess(guess)
	var response string

	switch result {
	case "TOO_LOW":
		response = fmt.Sprintf("%s_TOO_LOW_GUESS\n", session.Token)
	case "TOO_HIGH":
		response = fmt.Sprintf("%s_TOO_HIGH_GUESS\n", session.Token)
	case "CORRECT":
		response = fmt.Sprintf("%s_CORRECT_GUESS_ATTEMPTS:%d\n", session.Token, session.Game.GetAttempts())
	}

	conn.Write([]byte(response))

	if gameOver {
		s.startNewGame(conn, session)
	}
}

// Start a new guessing game
func (s *Server) startNewGame(conn net.Conn, session *ClientSession) {
	session.Game.Reset()
	conn.Write([]byte(fmt.Sprintf("%s_NEW_GAME_STARTED\n", session.Token)))
}

// Read client input
func readClientInput(conn net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buffer[:n])), nil
}
