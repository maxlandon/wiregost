package context

import (
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
	Context ConsoleContext
)

// ConsoleContext - Stores all variables needed for console context
type ConsoleContext struct {
	ClientID  uuid.UUID               // Unique user ID for module requests
	User      *dbpb.User              // User information sent back after auth
	Shell     *readline.Instance      // Shell object
	Config    *clientpb.ConsoleConfig // Shell configuration
	Menu      *string                 // Current shell menu
	Workspace *dbpb.Workspace         // Current workspace
	Module    *modulepb.Module        // Current module
	Ghost     *ghostpb.Ghost          // Current implant
	Jobs      *int                    // Number of jobs
	Ghosts    *int                    // Number of connected implants
}

// SetConsoleContext - Set the context used by commands
func SetConsoleContext() {

}

// GetContext - Get all information necessary to console upon connection
func GetContext(cli clientpb.ConnectionRPCClient) {

	// Info Request
	info, _ := cli.GetConnectionInfo(ctx, &clientpb.ConnectionInfoRequest{}, grpc.EmptyCallOption{})

	// Set fields
	Context.Workspace = info.Workspace
	*Context.Jobs = int(info.Jobs)
	*Context.Ghosts = int(info.Ghosts)

}

// GetVersion - Get client & server version information upon connection
func GetVersion(cli clientpb.ConnectionRPCClient) (info *clientpb.Version) {
	info, _ = cli.GetVersion(ctx, &clientpb.Empty{}, grpc.EmptyCallOption{})
	return
}
