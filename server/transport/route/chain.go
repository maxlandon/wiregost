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

// LastNode returns the last node of the node list.
// If the chain is empty, an empty node will be returned.
// If the last node is a node group, the first node in the group will be returned.
func (c *Chain) LastNode() Node {
	group := c.LastNodeGroup()
	return group.Nodes[0]
}

// LastNodeGroup returns the last group of the group list.
func (c *Chain) LastNodeGroup() *NodeGroup {
	if c.IsEmpty() {
		return nil
	}
	return c.NodeGroups[len(c.NodeGroups)-1]
}

// AddNode appends the node(s) to the chain.
func (c *Chain) AddNode(nodes ...Node) {
	if c == nil {
		return
	}
	for _, node := range nodes {
		c.NodeGroups = append(c.NodeGroups, NewNodeGroup(node))
	}
}

// AddNodeGroup appends the group(s) to the chain.
func (c *Chain) AddNodeGroup(groups ...*NodeGroup) {
	if c == nil {
		return
	}
	for _, group := range groups {
		c.NodeGroups = append(c.NodeGroups, group)
	}
}

// IsEmpty checks if the chain is empty.
// An empty chain means that there is no proxy node or node group in the chain.
func (c *Chain) IsEmpty() bool {
	return c == nil || len(c.NodeGroups) == 0
}
