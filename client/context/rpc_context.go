package context

import (
	"github.com/maxlandon/wiregost/client/assets"
	contextpb "github.com/maxlandon/wiregost/proto/v1/gen/go/context"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
)

// RPCContext - Holds all context metadata used in Wiregost, passed for each request made by a client console
type RPCContext struct {
	Token     string         // Unique number per console instance (for running modules, etc)
	Workspace dbpb.Workspace // Current workspace
	User      dbpb.User      // User owning the process context
	Menu      string         // Current shell menu
	Ghost     ghostpb.Ghost  // Current implant
}

// SetMetadata - Set the context used by gRPC calls
func SetMetadata() (new *contextpb.RPCContext) {

	new = &contextpb.RPCContext{
		Token:       assets.Token,
		WorkspaceID: Context.Workspace.ID,
		Username:    Context.User.Name,
		Admin:       Context.User.Admin,
		Menu:        Context.Menu,
		GhostID:     Context.Ghost.ID,
	}

	return
}
