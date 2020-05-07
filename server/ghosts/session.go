package ghosts

import (
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

// Session - Represents a currently connected ghost implant. It can be a ghost running
// on Linux, Windows, BSD, etc, and have a different set of methods, as long as it
// implements the base functions needed for a ghost session to run.
type Session interface {
	ID() (id uint32)                          // Session ID
	OS() (os string)                          // Session Operating system
	Owner() (owner *serverpb.User)            // User owning the session
	Permissions() (perms ghostpb.Permissions) // Who has the right to use implant
	Request()                                 // Function for sending a message to implant (transport-agnostic)
}
