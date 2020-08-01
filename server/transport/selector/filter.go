package selector

import (
	"time"

	"github.com/maxlandon/wiregost/server/transport/route"
)

// Filter is used to filter a node during the selection process
type Filter interface {
	Filter([]route.Node) []route.Node // Filter will be particularly important.
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
