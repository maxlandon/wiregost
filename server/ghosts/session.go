package ghosts

import (
	"sync"

	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
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
	Core       Core       // All core information/data about this ghost implant
	FileSystem FileSystem // File system methods
	Net        Net        // Network info
	Windows    Windows    // Windows-specific: only available on Windows implants
	Execute    Execute    // Generic execute methods
}

// NewGhost - A ghost implant has registered/connected: depending on its plaform and various other informations,
// register the underlying ghost struct to all appropriate interfaces.
// We return the object, because maybe in the caller function we want to register other interfaces, such as transport/RPC ones.
func NewGhost(new *ghostpb.Ghost) (g *Ghost) {

	// New generic type
	core := generic.NewGhost(new)

	// Register base interfaces
	g.Core = core
	g.FileSystem = core
	g.Net = core
	g.Execute = core

	// Register OS-specific interfaces
	switch g.Core.OS() {
	case "windows":
		g.Windows = windows.NewGhost(core)
	case "linux":
	case "darwin":
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

// Core - Represents a currently connected ghost implant. It can be a ghost running
// on Linux, Windows, BSD, etc, and have a different set of methods, as long as it
// implements the base functions needed for a ghost session to run.
type Core interface {
	ID() (id uint32)                          // Session ID
	Owner() (owner *dbpb.User)                // User owning the session
	Permissions() (perms ghostpb.Permissions) // Who has the right to use implant
	OS() (os string)                          // Session Operating system
	Info() (info *ghostpb.Ghost)              // Information

	// Functions returning core and networking capabilities
	// Transport()
	// Route()
}

type Net interface{}
type FileSystem interface{}
type Proc interface{}
type Execute interface{}
type Windows interface{}
