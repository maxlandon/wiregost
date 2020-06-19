package context

import (
	"context"

	"github.com/google/uuid"
	"github.com/lmorg/readline"
	"google.golang.org/grpc"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

var (
	// Context - The console context object
	Context = newContext()
)

// Menu Contexts
const (
	// MainMenu - Available only in main menu
	MainMenu = "main"
	// ModuleMenu - Available only when a module is loaded
	ModuleMenu = "module"
	// GhostMenu - Available only when interacting with a ghost implant
	GhostMenu = "ghost"
)

// ConsoleContext - Stores all variables needed for console context
type ConsoleContext struct {
	ClientID  uuid.UUID              // Unique user ID for module requests
	User      dbpb.User              // User information sent back after auth
	Shell     *readline.Instance     // Shell object
	Config    clientpb.ConsoleConfig // Shell configuration
	Menu      string                 // Current shell menu
	Workspace dbpb.Workspace         // Current workspace
	Module    modulepb.Module        // Current module
	Ghost     ghostpb.Ghost          // Current implant
	Jobs      int                    // Number of jobs
	Ghosts    int                    // Number of connected implants
}

func newContext() (ctx *ConsoleContext) {

	ctx = &ConsoleContext{}
	ctx.User = dbpb.User{}
	ctx.Config = clientpb.ConsoleConfig{}
	ctx.Workspace = dbpb.Workspace{}
	ctx.Module = modulepb.Module{}
	ctx.Ghost = ghostpb.Ghost{}

	return
}

// GetConnectionInfo - Set the context used by commands & shell
func GetConnectionInfo(cli clientpb.ConnectionRPCClient) (info *clientpb.ConnectionInfo, config *clientpb.ConsoleConfig) {

	// Info Request
	info, _ = cli.GetConnectionInfo(context.Background(), &clientpb.ConnectionInfoRequest{}, grpc.EmptyCallOption{})

	// Set fields (beware of nil fields in pb message)
	// Context.Workspace = (*info.Workspace)
	Context.Jobs = int(info.Jobs)
	Context.Ghosts = int(info.Ghosts)
	Context.Menu = MainMenu

	// Get and use Console configuration
	config = info.ConsoleConfig

	return
}

// GetVersion - Get client & server version information upon connection
func GetVersion(cli clientpb.ConnectionRPCClient) (info *clientpb.Version) {
	info, _ = cli.GetVersion(context.Background(), &clientpb.Empty{}, grpc.EmptyCallOption{})
	return
}
