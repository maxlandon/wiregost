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
// This object is used in many cases:
// - When modules request a ghost implant to perform a task
// - When console users request a ghost implant
// - Monitor timeouts for modules, consoles, and others.
type RPCContext struct {
	ClientID uuid.UUID     // Unique number per console instance (for running modules, etc)
	User     dbpb.User     // User owning the process context
	Ghost    ghostpb.Ghost // Current implant
}

// NewContextRPC - Set the context used by gRPC calls
func NewContextRPC() (ctx context.Context) {

	new := RPCContext{
		// ClientID: Context.ClientID,
		// User:     Context.User,
		// Ghost:    Context.Ghost,
	}

	ctx = context.WithValue(base, MetadataKey, new)

	return
}

// GetMetadata - Used by the server and DB to get the context of a RPC call
func GetMetadata(in context.Context) (ctx RPCContext) {
	return in.Value(MetadataKey).(RPCContext)
}
