package c2

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/yamux"
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

// Session - The transport layer of a ghost implant, server-side. The connection, whether
// it is through DNS, MTLS, HTTP, KCP or any other, is always stored as a Session object.
// Thus, there is a 1-to-1 mapping between the sessions stored in this package, and Ghost
// objects instantiated in the `ghosts` package.
//
// NOTE: The Conn field might be empty (and therefore the *Session also), because if the implant is pivoted,
// the physical connection will not be between the server and the implant. Therefore, the Conn object is
// used for non-pivoted implants, because any routed traffic will need this physical conn to multiplex.
type Session struct {
	ID        *uuid.UUID             // Identity of the related Ghost object
	C2        *yamux.Stream          // A logical connection, reserved to the ghost's requests/responses
	Send      chan []byte            // Outgoing messages
	Resp      map[uint64]chan []byte // Incoming messages, checked for replay attacks
	respMutex *sync.RWMutex
}

// Request - Send a request to a ghost implant connected through a custom transport (DNS, MTLS, HTTPS)
// This functions used the C2 field of the implant, not the Conn.
func (s *Session) Request(msgType uint32, timeout time.Duration, req []byte) (res []byte, err error) {

	return
}
