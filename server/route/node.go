package route

import "net/url"

// Node - A proxy node, mainly used to construct a proxy chain
type Node struct {
	ID        uint32
	GhostID   uint32
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
