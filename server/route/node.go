package route

import (
	"net/url"

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

// ToProtobuf - Helper function used to pack Node information and use it in DB/Comms/Requests
func (n *Node) ToProtobuf() (proto routepb.Node) {
	return
}

// ParseNode - Parses a Protobuf node object and returns a node usable by the routing system.
func ParseNode(proto routepb.Node) (node Node, err error) {
	return
}
