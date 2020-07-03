package route

import (
	"net/url"
	"time"

	routepb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport/route"
)

// Node - A proxy node, mainly used to construct a proxy chain
type Node struct {
	ID        uint32
	GhostID   uint32 // A node must be a ghost implant that will route traffic
	Addr      string
	Host      string
	Protocol  string
	Transport string
	Remote    string   // Remote address, used by tcp/udp port forwarding
	url       *url.URL // Raw url
	User      *url.Userinfo
	Values    url.Values

	// DialOptions      []DialOption
	// HandshakeOptions []HandshakeOption
	// ConnectOptions   []ConnectOption
	// Client           *Client
	// marker           *failMarker
	// Bypass           *Bypass
}

// MarkDead - Makes node fail status, and an optional error message
func (n *Node) MarkDead() {
}

// ResetDead - Resets the node fail status
func (n *Node) ResetDead() {
}

// Clone - Clones the node, will prevent data race
func (n *Node) Clone() (clone Node) {
	return
}

// Get returns node parameter specified by key.
func (node *Node) Get(key string) string {
	return node.Values.Get(key)
}

// GetBool converts node parameter value to bool.
func (node *Node) GetBool(key string) (b bool) {
	// b, _ := strconv.ParseBool(node.Values.Get(key))
	return b
}

// GetInt converts node parameter value to int.
func (node *Node) GetInt(key string) (n int) {
	// n, _ := strconv.Atoi(node.Values.Get(key))
	return n
}

// GetDuration converts node parameter value to time.Duration.
func (node *Node) GetDuration(key string) (d time.Duration) {
	// d, _ := time.ParseDuration(node.Values.Get(key))
	return
}

// ToProtobuf - Helper function used to pack Node information and use it in DB/Comms/Requests
func (n *Node) ToProtobuf() (proto routepb.Node) {
	return
}

// ParseNode - Parses a Protobuf node object and returns a node usable by the routing system.
func ParseNode(proto routepb.Node) (node Node, err error) {
	return
}
