package route

// Chain - A proxy chain that holds a list of proxy nodes, or group of nodes
type Chain struct {
	IsRoute    bool
	Retries    int
	NodeGroups []*NodeGroup // A list of group of nodes, from which to choose and use as route
	Route      []Node       // The list of nodes actually used to route traffic
}
