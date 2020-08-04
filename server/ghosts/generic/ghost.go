package generic

import (
	"sync"
	"time"

	"github.com/hashicorp/yamux"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
)

// Ghost - The base implementation for all implants in Wiregost.
// It provides only the set of methods necessary to implement the "Session" interface.
// This means its the bare minimum to identify and interact with an implant, and it
// does not include any core capability.
type Ghost struct {
	Info      *ghostpb.Ghost         // Protobuf Information
	Host      *dbpb.Host             // Host associated with this implant (on which it runs)
	C2        *yamux.Stream          // A logical connection, reserved to the ghost's requests/responses
	Send      chan []byte            // Outgoing messages
	Resp      map[uint64]chan []byte // Incoming messages, checked for replay attacks
	respMutex *sync.RWMutex
}

// NewGhost - Returns a ghost object, instantiated after an implant has registered.
func NewGhost(new *ghostpb.Ghost) (ghost *Ghost) {
	ghost = &Ghost{
		Info: new,
	}
	return
}

// LogFile - Returns the path to the log file to use for this implant.
func (g *Ghost) LogFile() (path string) {
	return
}

// Cleanup - Handles cleanup of all data related to this session. If lost is true, we will go further and
// cleanup various elements in the database, like route chains containing this implant, while if false
//the implant is just down, so we do a lighter cleanup.
func (g *Ghost) Cleanup(lost string) (err error) {
	return
}

// ID - Returns the implant ID
func (g *Ghost) ID() (id uint32) {
	return
}

// Request - Send a request to a ghost implant connected through a custom transport (DNS, MTLS, HTTPS)
// This functions used the C2 field of the implant, not the Conn.
func (s *Ghost) Request(msgType uint32, timeout time.Duration, req []byte) (res []byte, err error) {

	return
}
