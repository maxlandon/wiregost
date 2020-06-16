package context

import (
	"context"

	"github.com/google/uuid"

	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
)

var (
	base = context.Background()
)

const (
	// MetadataKey - Used to reference the Data struct contained in the context
	MetadataKey = "wiregost"
)

// RPCContext - Holds all context metadata used in Wiregost, passed for each request made by a client console
type RPCContext struct {
	ClientID  *uuid.UUID      // Unique number per console instance (for running modules, etc)
	Workspace *dbpb.Workspace // Current workspace
	User      *dbpb.User      // User owning the process context
	Menu      *string         // Current shell menu
	Ghost     *ghostpb.Ghost  // Current implant
}

// NewContextRPC - Set the context used by gRPC calls
func NewContextRPC() (ctx context.Context) {

	new := RPCContext{
		ClientID:  &Context.ClientID,
		Workspace: Context.Workspace,
		User:      Context.User,
		Menu:      Context.Menu,
		Ghost:     Context.Ghost,
	}

	ctx = context.WithValue(base, MetadataKey, new)

	return
}
