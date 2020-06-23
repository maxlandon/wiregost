package ghosts

import (
	"sync"

	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	"github.com/maxlandon/wiregost/server/ghosts/generic"
	"github.com/maxlandon/wiregost/server/ghosts/windows"
)

var (
	// Ghosts - All currently connected ghost implants/sessions.
	Ghosts = ghosts{
		Connected: &map[uint32]*Ghost{},
		mutex:     &sync.Mutex{},
	}
)

// Ghost - A connected ghost implant. Provides access to all functions available, cross-platform and transport-agnostic.
// The various interfaces, such as FileSystem, are sometimes only meant to provide an easier, categorized access to implant methods.
// Because objects may implement several interfaces at once, the code in windows/, linux/ and darwin/ may use these same functions
// directly, without passing through this interface.
type Ghost struct {
	Core       Core               // All core information/data about this ghost implant
	FileSystem generic.FileSystem // File system methods
	Net        generic.Net        // Network info
	Windows    windows.Windows    // Windows-specific: only available on Windows implants
	Execute    generic.Execute    // Generic execute methods
}

// NewGhost - A ghost implant has registered/connected: depending on its plaform and various other informations,
// register the underlying ghost struct to all appropriate interfaces.
// We return the object, because maybe in the caller function we want to register other interfaces, such as transport/RPC ones.
func NewGhost(new *ghostpb.Ghost) (g *Ghost) {

	// New generic type
	core := generic.NewGhost(new)

	// For each Operating System, we carefully register all interfaces, so that we avoid
	// potential clashes and mistakes at compile-time.
	switch g.Core.Info().OS {
	case "windows":
		g.Core = core
		g.FileSystem = core
		g.Net = core
		g.Execute = core
		g.Windows = windows.NewGhost(core)
	case "linux":
		g.Core = core
		g.FileSystem = core
		g.Net = core
	case "darwin":
		g.Core = core
		g.FileSystem = core
		g.Net = core
	}

	// Add to connected ghosts
	Ghosts.Add(g)

	return
}

type ghosts struct {
	Connected *map[uint32]*Ghost
	mutex     *sync.Mutex
}

// Add - A new ghost is registered, authenticated and connected
func (g *ghosts) Add(new *Ghost) {
	g.mutex.Lock()
	(*g.Connected)[new.Core.ID()] = new
	g.mutex.Unlock()
}

func (g *ghosts) Get(id uint32) (ghost *Ghost) {
	return (*g.Connected)[id]
}

// Core - Represents a currently connected ghost implant. It can be a ghost running
// on Linux, Windows, BSD, etc, and have a different set of methods, as long as it
// implements the base functions needed for a ghost session to run.
type Core interface {
	ID() (id uint32)             // Session ID
	Info() (info *ghostpb.Ghost) // Information

	// Functions returning core and networking capabilities
	// Transport()
	// Route()
}
