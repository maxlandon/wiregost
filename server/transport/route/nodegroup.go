package route

import (
	"sync"

	routepb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport/route"
)

// NodeGroup - A group of proxy nodes.
// These groups might have several usages:
// - Gather several implants (in the same network, with some common attributes),
//   that cannot be used at the same time in a given chain.
// - Use various selectors among these nodes. See Strategies & Filters
type NodeGroup struct {
	ID          uint32
	Name        string
	Description string
	Nodes       []Node
	// SelectorOptions []selector.SelectOption
	// Selector        selector.NodeSelector
	mux sync.RWMutex
}

// NewNodeGroup - Creates a node group
func NewNodeGroup(nodes ...Node) (group *NodeGroup) {
	return
}

// AddNodes - Add one or more nodes to a group
func (g *NodeGroup) AddNodes(nodes ...Node) {
}

// SetNodes - Replace all nodes in a group, and return the old ones
func (g *NodeGroup) SetNodes(new ...Node) (old []Node) {
	return
}

// SetSelector sets node selector with options for the group.
// func (group *NodeGroup) SetSelector(selector NodeSelector, opts ...SelectOption) {
// }

// Next selects a node from group.
func (g *NodeGroup) Next() (node Node, err error) {
	return
}

// ToProtobuf - Export all nodes in this group
func (g *NodeGroup) ToProtobuf() (nodes []routepb.Node) {

	for _, node := range g.Nodes {
		nodes = append(nodes, node.ToProtobuf())
	}
	return
}
