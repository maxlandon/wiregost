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
	Proc       generic.Proc       // Process information/manipulation
	Execute    generic.Execute    // Generic execute methods
	Windows    windows.Windows    // Windows-specific: only available on Windows implants
}

// NewGhost - A ghost implant has registered/connected: depending on its plaform and various
// other informations, register the underlying ghost struct to all appropriate interfaces.
// This function registers the Ghost for usage by modules and console users.
//
// NOTE: This function does not take care of handling the initial registration messages that
// contain all target/implant information. This means:
// - Transport components are all up and running, with according security needs.
// - All server-to-ghost RPC handlers are registered.
func NewGhost(new *ghostpb.Ghost) (g *Ghost) {

	// New ghost base type
	base := generic.NewGhost(new)

	// Register core interfaces, generally working cross-platform
	// through the generic.Ghost base type.
	g.Core = base
	g.FileSystem = base
	g.Net = base
	g.Proc = base

	// OS-specific interfaces/objects
	switch g.Core.Info().OS {
	case "windows":
		g.Windows = windows.NewGhost(base)
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

func (g *ghosts) Get(id uint32) (ghost *Ghost) {
	return (*g.Connected)[id]
}

// Core - Represents a currently connected ghost implant. It can be a ghost running
// on Linux, Windows, BSD, etc, and have a different set of methods, as long as it
// implements the base functions needed for a ghost session to run.
type Core interface {
	ID() (id uint32)             // Session ID
	Info() (info *ghostpb.Ghost) // Information
}
