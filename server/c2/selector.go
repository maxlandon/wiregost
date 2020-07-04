package route

import (
	"github.com/maxlandon/wiregost/server/c2/route"
	"github.com/maxlandon/wiregost/server/c2/selector"
)

// NodeSelector - The node selector is used to perform filtering, strategy building and finally node selection,
// between a list of Nodes. Various types of filters will be used on these NodeGroups, like implant ownership,
// implant communication capacities, versions, stealth networks, etc...
//
// Additionally: This NodeSelector performs as a pluggable strategy and filtering tool, and therefore many tools
// and controls can be implemented in order to refine the routing process in Wiregost.
type NodeSelector interface {
	Select(nodes []route.Node, opts ...selector.SelectOption) (route.Node, error)
}
