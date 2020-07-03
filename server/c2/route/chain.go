package route

import (
	routepb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport/route"
)

// Chain - A proxy chain that holds a list of proxy nodes, or group of nodes
// It has methods to connect to the last node of a selected route.
// The Chain NEVER HOLDS the server as a proxy node.
// Also, the Chain is always stored on the Server, and they are never built by implants
// themselves. This ensures all data necessary to route opening/reopening/closing is always
// on the Server, and only on the Server.
type Chain struct {
	IsRoute    bool
	Retries    int
	NodeGroups []*NodeGroup // Each node group is at least one ghost implant
	Route      []Node       // The list of nodes actually used to route traffic
}

// ToProtobuf - Returns a route object used by DB/Comms/Requests
func (c *Chain) ToProtobuf() (route routepb.Route) {

	// For each node, call method ToProtobuf()

	// Add them to route

	return
}

// NewChain - Creates a proxy chain, empty.
func NewChain() (chain *Chain) {
	return
}

// ParseRoute - Takes a routepb object and stores it in route, with several adjustements.
// This should be used when a user submits an OpenRoute request:
// - Either the route is specified by the user (he knows which route is best)
// - Either the server automatically devises it.
func (c *Chain) ParseRoute(route routepb.Route) (err error) {

	// For each node
	// ParseNode()

	// c.IsRoute = true     // Should notify something like "this chain is locked or ready to be used"
	return
}
