package ghosts

import (
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
)

// Session - Represents a currently connected ghost implant. It can be a ghost running
// on Linux, Windows, BSD, etc, and have a different set of methods, as long as it
// implements the base functions needed for a ghost session to run.
type Session interface {
	ID() (id uint32)                          // Session ID
	OS() (os string)                          // Session Operating system
	Owner() (owner *dbpb.User)                // User owning the session
	Permissions() (perms ghostpb.Permissions) // Who has the right to use implant
	Request()                                 // Function for sending a message to implant (transport-agnostic)
}
