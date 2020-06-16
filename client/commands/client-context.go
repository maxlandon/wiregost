package commands

import (
	"github.com/google/uuid"
	"github.com/lmorg/readline"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

var (
	// Context - The console context object
	Context ConsoleContext
	// ContextRPC - The context object used with gRPC
	ContextRPC RPCContext
)

const (
	// MetadataKey - Used to reference the Data struct contained in the context
	MetadataKey = "wiregost"
)

// ConsoleContext - Stores all variables needed for console context
type ConsoleContext struct {
	ClientID  *uuid.UUID              // Unique user ID for module requests
	User      *serverpb.User          // User information sent back after auth
	Shell     *readline.Instance      // Shell object
	Config    *clientpb.ConsoleConfig // Shell configuration
	Menu      *string                 // Current shell menu
	Workspace *dbpb.Workspace         // Current workspace
	Module    *modulepb.Module        // Current module
	Ghost     *ghostpb.Ghost          // Current implant
	Jobs      *int                    // Number of jobs
	Ghosts    *int                    // Number of connected implants
}

// RPCContext - Holds all context metadata used in Wiregost, passed for each request made by a client console
type RPCContext struct {
	ClientID  *uuid.UUID      // Unique number per console instance (for running modules, etc)
	Workspace *dbpb.Workspace // Current workspace
	User      *serverpb.User  // User owning the process context
	Menu      *string         // Current shell menu
	Ghost     *ghostpb.Ghost  // Current implant
}

// SetConsoleContext - Set the context used by commands
func SetConsoleContext() {

}

// SetContextRPC - Set the context used by gRPC calls
func SetContextRPC() {

}
