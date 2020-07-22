package https

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// Session - Holds data related to a ghost implant HTTP session
type Session struct{}

// Sessions - All implants currently connected through HTTP
type Sessions struct{}

// Add - Add an HTTP session
func (s *Sessions) Add(session *Session) {
}

// Get - Get an HTTP session
func (s *Sessions) Get(sessionID string) (sess *Session) {
}

// Remove - Remove an HTTP session
func (s *Sessions) Remove(sessionID string) {
}

// LoggingMiddleware - Logs the content of an implant HTTP request/response.
func LoggingMiddleware(next http.Handler) (logged http.Handler) {
	return
}

func newSession() (sess *Session) {
	// return &Session{
	//         ID:      newHTTPSessionID(),
	//         Started: time.Now(),
	//         replay:  map[string]bool{},
	// }
}

// newSessionID - Get a 128bit session ID
func newSessionID() string {
	buf := make([]byte, 16)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
