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

type Wiregost struct {
	// Connections
	Endpoint *endpoint.Endpoint

	// DB Access
	DbManager db.DBManager

	// User
	UserManager *user.UserManager

	// Workspace
	WorkspaceManager *workspace.WorkspaceManager

	// ModuleStackManager
	ModuleStackManager *modules.ModuleStackManager

	// Logger

	// Server
	ServerManager *server.ServerManager

	// Compiler
	CompilerManager *compiler.Manager
}

func NewWiregost() *Wiregost {
	wiregost := &Wiregost{
		Endpoint: endpoint.NewEndpoint(),
		// DB
		UserManager:        user.NewUserManager(),
		WorkspaceManager:   workspace.NewWorkspaceManager(),
		ModuleStackManager: modules.NewModuleStackManager(),
		ServerManager:      server.NewServerManager(),
		CompilerManager:    compiler.NewManager(),
	}

	return wiregost
}
