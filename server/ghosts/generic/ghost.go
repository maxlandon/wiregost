package generic

import (
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	"github.com/maxlandon/wiregost/server/c2"
)

// Ghost - The base implementation for all implants in Wiregost.
// It provides only the set of methods necessary to implement the "Session" interface.
// This means its the bare minimum to identify and interact with an implant, and it
// does not include any core capability.
type Ghost struct {
	Proto   *ghostpb.Ghost
	Session *c2.Session
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
