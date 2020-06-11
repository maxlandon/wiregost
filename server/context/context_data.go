package context

import (
	"github.com/google/uuid"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

const (
	// MetadataKey - Used to reference the Data struct contained in the context
	MetadataKey = "wiregost"
)

// Data - Holds all context metadata used in Wiregost, passed for each request made by a client console
type Data struct {
	ClientID    uuid.UUID     // Unique number per console instance (for running modules, etc)
	WorkspaceID uint32        // Current workspace in which context was launched
	User        serverpb.User // User owning the process context
}
