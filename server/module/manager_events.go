package module

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

// Events - "Point de passage obligé" for all module log messages in Wiregost.
func (m *managers) Events(_ *pb.Empty, serv pb.Manager_EventsServer) error {

	// We identify the client, getting a ClientID for the console,

	// Then, for all events coming in:
	for {
		// Each message coming from the ClientConn talking to the stack

		// Log it to all files needing it first

		// Write clientID to log message, and forward to consoles (MAYBE NOT NEEDED)

		// Forward to console (through stream-as-conn)
	}

	return status.Errorf(codes.Unimplemented, "method Events not implemented")
}
