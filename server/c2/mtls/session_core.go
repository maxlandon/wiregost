package mtls

import (
	"net"

	"github.com/hashicorp/yamux"
)

// Session - Handles all transport communications for a single Ghost implant.
// Transport communications include:
// - Ghost requests for core functionality
// - Non-ghost traffic to be routed to either a pivoted ghost, or to a target host
//
// The Session should seamlessly integrate both traffics. Because it might have to
// handle many connections to various destinations, it is best that this object
// remains free of any routing logic.
type Session struct {
	*yamux.Session
}

// HandleSession - The implant is registered: handle all possible requests/responses going from/to it.
func HandleSession(conn net.Conn) {
}
