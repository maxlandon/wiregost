package route

import (
	"time"

	"github.com/maxlandon/wiregost/server/c2/route"
)

// NodeSelector - The node selector is used to perform filtering, strategy building and finally node selection,
// between a list of Nodes. Various types of filters will be used on these NodeGroups, like implant ownership,
// implant communication capacities, versions, stealth networks, etc...
//
// Additionally: This NodeSelector performs as a pluggable strategy and filtering tool, and therefore many tools
// and controls can be implemented in order to refine the routing process in Wiregost.
type NodeSelector interface {
	Select(nodes []route.Node, opts ...SelectOption) (route.Node, error)
}

// SelectOption - Option used when making a select call
type SelectOption func(*SelectOptions)

// SelectOptions - Options for node selection
type SelectOptions struct {
	Filters  []Filter
	Strategy Strategy
}

// Filter is used to filter a node during the selection process
type Filter interface {
	Filter([]route.Node) []route.Node // Filter will be particularly important.
	String() string
}

// Strategy is a selection strategy e.g random, round-robin.
type Strategy interface {
	Apply([]route.Node) route.Node
	String() string
}

// FailFilter - This is an example object that should have a much larger role in Wiregost routing mechanism:
// This object, or any further version of it, should have all access to needed to perform various, often
// iterative/repeated filtering on all ghost implants. This will ensure permissions, network stacks, operational
// requirements will all be respected.
type FailFilter struct {
	MaxFails    int
	FailTimeout time.Duration
}
