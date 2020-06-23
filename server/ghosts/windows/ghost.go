package windows

import "github.com/maxlandon/wiregost/server/ghosts/generic"

// Ghost - A ghost implant running on a Windows host. This implant has many
// fields and methods that are most of the time specific to Windows OS & platforms.
type Ghost struct {
	*generic.Ghost // This is used to fullfil the 'Core' interface of a ghost object. Also provides access to all state

	// Windows-specific objects, such as Registry, .NET, LDAP, etc...
}

// NewGhost - This function initializes all state necessary for a Windows implant to perform its functions.
func NewGhost(new *generic.Ghost) (ghost *Ghost) {
	ghost = &Ghost{
		new, // The state of the implant (owner, permissions, info, etc) and core functions (filesystem, net, etc...)
	}

	return
}
