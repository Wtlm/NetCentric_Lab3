package main

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"sync"
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

// UserManager handles user authentication operations
type UserManager struct {
	users map[string]User
	mutex sync.RWMutex
}

// NewUserManager creates a new UserManager instance
func NewUserManager() *UserManager {
	return &UserManager{
		users: make(map[string]User),
	}
}

// LoadUsers reads users from a JSON file
func (um *UserManager) LoadUsers(filename string) error {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	// Read existing users
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, &um.users)
}

// SaveUsers writes users to a JSON file
func (um *UserManager) SaveUsers(filename string) error {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	data, err := json.MarshalIndent(um.users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// Authenticate checks user credentials
func (um *UserManager) Authenticate(username, password string) bool {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	user, exists := um.users[username]
	if !exists {
		return false
	}

	// Simple password comparison (use secure hashing in production)
	decodedPassword, _ := base64.StdEncoding.DecodeString(user.Password)
	return string(decodedPassword) == password
}
