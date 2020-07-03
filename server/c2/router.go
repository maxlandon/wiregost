package route

import "github.com/maxlandon/wiregost/server/c2/route"

// Router - Serves, initiates and handle connections all along a provided chain.
// It is somehow the center of the routing mechanism in Wiregost, because it holds state of
// pretty much everything needed to initiate and route communications.
type Router struct {
	Node    route.Node
	Server  *route.Server // Should be used only when proxy module is on (MAYBE NOT, MAYBE NEEDED ANYWAY)
	Handler route.Handler // Handles all connection details to nodes
	Chain   *route.Chain  // Stores routes and dials nodes
	// resolver gost.Resolver
	// hosts    *gost.Hosts
}
