package commands

import (
	"github.com/google/uuid"
	"github.com/lmorg/readline"

	"github.com/maxlandon/wiregost/client/assets"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

// Context - The console context object
var Context ConsoleContext

// ConsoleContext - Stores all variables needed for console context
type ConsoleContext struct {
	ClientID  *uuid.UUID            // Unique user ID for module requests
	User      *serverpb.User        // User information sent back after auth
	Shell     *readline.Instance    // Shell object
	Config    *assets.ConsoleConfig // Shell configuration
	Menu      *string               // Current shell menu
	Workspace *dbpb.Workspace       // Current workspace
	Module    *modulepb.Module      // Current module
	Ghost     *ghostpb.Ghost        // Current implant
	Jobs      *int                  // Number of jobs
	Ghosts    *int                  // Number of connected implants
	// DBContext context.Context      // DB queries context
}
