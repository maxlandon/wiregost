package c2

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	// Sessions - Holds all transport-layer ghost connections
	Sessions = &sessions{
		Connected: &map[string]*Session{},
		mutex:     &sync.RWMutex{},
	}
)

type sessions struct {
	Connected *map[string]*Session // All implant transports currently connected
	mutex     *sync.RWMutex
}

// Session - The transport-layer object of a ghost implant, server-side.
// The connection, whether it is through DNS,MTLS,HTTP,KCP or any other, is always stored as
// a Session object.
//
// Thus, there is a 1-to-1 mapping between the sessions stored in this package, and Ghost objects
// instantiated in the `ghosts` package.
type Session struct {
	ID        *uuid.UUID             // Provides identity of the related Ghost object
	Send      chan []byte            // Outgoing messages
	Resp      map[uint64]chan []byte // Incoming messages, checked for replay attacks
	respMutex *sync.RWMutex
}

// Request - Send a request to a ghost implant connected through a custom transport (DNS, MTLS, HTTPS)
func (s *Session) Request(msgType uint32, timeout time.Duration, req []byte) (res []byte, err error) {

	return
}
