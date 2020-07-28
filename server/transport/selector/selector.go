package selector

import "github.com/maxlandon/wiregost/server/c2/route"

// SelectOption - Option used when making a select call
type SelectOption func(*SelectOptions)

// SelectOptions - Options for node selection
type SelectOptions struct {
	Filters  []Filter
	Strategy Strategy
}

// Strategy is a selection strategy e.g random, round-robin.
type Strategy interface {
	Apply([]route.Node) route.Node
	String() string
}

// type NodeSelector interface {
//         Select(nodes []Node, opts ...SelectOption) (route.Node, error)
// }
