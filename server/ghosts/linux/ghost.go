package linux

import "github.com/maxlandon/wiregost/server/ghosts/generic"

// Ghost - A ghost implant running on a Linux host
type Ghost struct {
	Base *generic.Ghost
}

// NewGhost - This function initializes all state necessary for a Windows implant to perform its functions.
func NewGhost(new *generic.Ghost) (ghost *Ghost) {
	ghost = &Ghost{
		new, // The state of the implant (owner, permissions, info, etc) and core functions (filesystem, net, etc...)
	}

	return
}
