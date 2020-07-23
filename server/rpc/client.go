package rpc

import "github.com/maxlandon/wiregost/server/rpc/generic"

// GhostClient - The RPC layer of a ghost implant connection. This layer is attached to a Ghost
// type in the `ghosts` package. This makes an easier an more transparent way of switching the
// RPC component, depending on the transport used.
// However, the methods offered by this GhostClient object should always be the same.
// Therefore, the GhostClient must always do the work of using the right RPC stub.
//
// At this level, several things should be managed:
// - Performance (timeouts)
// - Permissions (context.Context objects should be passed around)
type GhostClient struct {
	Type    string
	Generic *generic.Client
}
