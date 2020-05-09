package windows

import "github.com/maxlandon/wiregost/server/ghosts/generic"

// Ghost - A ghost implant running on a Windows host
type Ghost struct {
	Base *generic.Ghost // Ensures this ghost has all cross-platform methods
}
