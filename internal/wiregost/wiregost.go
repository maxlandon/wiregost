package wiregost

import (
	"github.com/maxlandon/wiregost/internal/compiler"
	"github.com/maxlandon/wiregost/internal/db"
	"github.com/maxlandon/wiregost/internal/endpoint"
	"github.com/maxlandon/wiregost/internal/modules"
	"github.com/maxlandon/wiregost/internal/server"
	"github.com/maxlandon/wiregost/internal/user"
	"github.com/maxlandon/wiregost/internal/workspace"
)

// Wiregost is the central point of the Wiregost server system.
// It is used to instantiate all managers and components.
type Wiregost struct {
	// Connections
	Endpoint *endpoint.Endpoint

	// DB Access
	DbManager db.Manager

	// User
	UserManager *user.Manager

	// Workspace
	WorkspaceManager *workspace.Manager

	// ModuleStackManager
	ModuleStackManager *modules.Manager

	// Logger

	// Server
	ServerManager *server.Manager

	// Compiler
	CompilerManager *compiler.Manager
}

// NewWiregost instantiates a new Wiregost server system
func NewWiregost() *Wiregost {
	wiregost := &Wiregost{
		Endpoint: endpoint.NewEndpoint(),
		// DB
		UserManager:        user.NewManager(),
		WorkspaceManager:   workspace.NewManager(),
		ModuleStackManager: modules.NewManager(),
		ServerManager:      server.NewManager(),
		CompilerManager:    compiler.NewManager(),
	}

	return wiregost
}
