package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"
)

// User represents a user in the system
type User struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	FullName  string   `json:"full_name"`
	Emails    []string `json:"emails"`
	Addresses []struct {
		Type    string `json:"type"`
		Address string `json:"address"`
	} `json:"addresses"`
}

type AuthManager struct {
	users    []User
	sessions map[string]int
	mu       sync.Mutex
}

// NewAuthManager initializes and loads users
func NewAuthManager(userFile string) (*AuthManager, error) {
	manager := &AuthManager{
		sessions: make(map[string]int),
	}

	data, err := ioutil.ReadFile(userFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &manager.users); err != nil {
		return nil, err
	}

	return manager, nil
}

// Authenticate checks username and password and assigns a session key
func (am *AuthManager) Authenticate(username, password string) (int, error) {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))

	am.mu.Lock()
	defer am.mu.Unlock()

	for _, user := range am.users {
		if user.Username == username && user.Password == encodedPassword {
			// Generate session key
			rand.Seed(time.Now().UnixNano())
			sessionKey := rand.Intn(1000) + 1

			// Store session key
			am.sessions[username] = sessionKey
			return sessionKey, nil
		}
	}
	return 0, errors.New("authentication failed")
}

// ValidateSession checks if a given session key is valid
func (am *AuthManager) ValidateSession(username string, sessionKey int) bool {
	am.mu.Lock()
	defer am.mu.Unlock()

	expectedKey, exists := am.sessions[username]
	return exists && expectedKey == sessionKey
}
