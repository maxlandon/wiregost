package context

import (
	"github.com/google/uuid"

	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

var (
	// ContextRPC - The context object used with gRPC
	ContextRPC RPCContext
)

const (
	// MetadataKey - Used to reference the Data struct contained in the context
	MetadataKey = "wiregost"
)

// RPCContext - Holds all context metadata used in Wiregost, passed for each request made by a client console
type RPCContext struct {
	ClientID  *uuid.UUID      // Unique number per console instance (for running modules, etc)
	Workspace *dbpb.Workspace // Current workspace
	User      *serverpb.User  // User owning the process context
	Menu      *string         // Current shell menu
	Ghost     *ghostpb.Ghost  // Current implant
}

// SetContextRPC - Set the context used by gRPC calls
func SetContextRPC() {

}
