package generic

import (
	"sync"
	"time"

	"github.com/hashicorp/yamux"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
)

// Ghost - The base implementation for all implants in Wiregost.
// It provides only the set of methods necessary to implement the "Session" interface.
// This means its the bare minimum to identify and interact with an implant, and it
// does not include any core capability.
type Ghost struct {
	Proto *ghostpb.Ghost // Protobuf Information
	// Session *c2.Session    // Session is independent from OS/architecture
	C2        *yamux.Stream          // A logical connection, reserved to the ghost's requests/responses
	Send      chan []byte            // Outgoing messages
	Resp      map[uint64]chan []byte // Incoming messages, checked for replay attacks
	respMutex *sync.RWMutex
}

// NewGhost - Returns a ghost object, instantiated after an implant has registered.
func NewGhost(new *ghostpb.Ghost) (ghost *Ghost) {
	ghost = &Ghost{
		Proto: new,
	}

	return
}

// ID - Returns the implant ID
func (g *Ghost) ID() (id uint32) {
	return
}

// Info - Returns all informations for this ghost implant
func (g *Ghost) Info() (info *ghostpb.Ghost) {
	return g.Proto
}

// Request - Send a request to a ghost implant connected through a custom transport (DNS, MTLS, HTTPS)
// This functions used the C2 field of the implant, not the Conn.
func (s *Ghost) Request(msgType uint32, timeout time.Duration, req []byte) (res []byte, err error) {

	return
}
