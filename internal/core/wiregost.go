package core

import (
	"github.com/maxlandon/wiregost/internal/db"
	"github.com/maxlandon/wiregost/internal/modules"
	"github.com/maxlandon/wiregost/internal/server"
	"github.com/maxlandon/wiregost/internal/user"
	"github.com/maxlandon/wiregost/internal/workspace"
)

type Wiregost struct {
	// Connections
	Endpoint *server.Endpoint

	// DB Access
	DbManager db.DBManager

	// User
	UserManager *user.UserManager

	// Workspace
	WorkspaceManager *workspace.WorkspaceManager

	// ModuleStackManager
	ModuleStackManager *modules.ModuleStackManager
	// Logger
}

func NewWiregost() *Wiregost {
	wiregost := &Wiregost{
		Endpoint: server.NewEndpoint(),
		// DB
		UserManager:        user.NewUserManager(),
		WorkspaceManager:   workspace.NewWorkspaceManager(),
		ModuleStackManager: modules.NewModuleStackManager(),
	}

	return wiregost
}
