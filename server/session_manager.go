package main

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

// ClientSession represents an individual client's game session
type ClientSession struct {
	Token      string
	Username   string
	Game       Game
	LastActive time.Time
}

// SessionManager handles multiple client sessions
type SessionManager struct {
	sessions map[string]*ClientSession
	mutex    sync.RWMutex
}

// NewSessionManager creates a new SessionManager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*ClientSession),
	}
}

// CreateSession generates a new client session
func (sm *SessionManager) CreateSession(username string) *ClientSession {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Generate unique token
	token := sm.generateUniqueToken()

	// Create new session
	session := &ClientSession{
		Token:      token,
		Username:   username,
		Game:       NewGame(),
		LastActive: time.Now(),
	}

	// Store session
	sm.sessions[token] = session

	return session
}

// GetSession retrieves a session by token
func (sm *SessionManager) GetSession(token string) (*ClientSession, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	session, exists := sm.sessions[token]
	return session, exists
}

// RemoveSession removes a session by token
func (sm *SessionManager) RemoveSession(token string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.sessions, token)
}

// generateUniqueToken creates a cryptographically secure unique token
func (sm *SessionManager) generateUniqueToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// CleanupExpiredSessions removes sessions that have been inactive for too long
func (sm *SessionManager) CleanupExpiredSessions(timeout time.Duration) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	now := time.Now()
	for token, session := range sm.sessions {
		if now.Sub(session.LastActive) > timeout {
			delete(sm.sessions, token)
		}
	}
}
