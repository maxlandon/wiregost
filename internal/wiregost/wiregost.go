package wiregost

import (
	"github.com/maxlandon/wiregost/internal/agents"
	"github.com/maxlandon/wiregost/internal/compiler"
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
	Endpoint           *endpoint.Endpoint
	UserManager        *user.Manager
	WorkspaceManager   *workspace.Manager
	ModuleStackManager *modules.Manager
	ServerManager      *server.Manager
	AgentManager       *agents.Manager
	CompilerManager    *compiler.Manager
}

// NewWiregost instantiates a new Wiregost server system
func NewWiregost() *Wiregost {
	wiregost := &Wiregost{
		Endpoint:           endpoint.NewEndpoint(),
		UserManager:        user.NewManager(),
		WorkspaceManager:   workspace.NewManager(),
		ModuleStackManager: modules.NewManager(),
		ServerManager:      server.NewManager(),
		AgentManager:       agents.NewManager(),
		CompilerManager:    compiler.NewManager(),
	}

	return wiregost
}
