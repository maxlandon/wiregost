package ghosts

import (
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

// Ghost - The base implementation for all implants in Wiregost.
// It provides only the set of methods necessary to implement the "Session" interface.
// This means its the bare minimum to identify and interact with an implant, and it
// does not include any core capability.
type Ghost struct {
	Base *ghostpb.Ghost
}

// NewGhost - Returns a ghost object, instantiated after an implant has registered.
func NewGhost() (ghost *Ghost) {
	return
}

// ID - Returns the implant ID
func (g *Ghost) ID() (id uint32) {
	return
}

// OS - Returns the ghost target operating system details
func (g *Ghost) OS() (os string) {
	return
}

// Owner - Returns the Wiregost user owning the implant
func (g *Ghost) Owner() (owner *serverpb.User) {
	return
}

// Permissions - Returns who has the right to control this ghost
func (g *Ghost) Permissions() (perms ghostpb.Permissions) {
	return
}

// Request - Sends a message to the ghost implant
func (g *Ghost) Request() {

}
