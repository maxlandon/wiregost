package route

import (
	"crypto/tls"
	"net"
	"net/url"
	"time"
)

// Handler - A proxy server handler. Configures and runs the connection.
type Handler interface {
	Init(options ...HandlerOption)
	Handle(net.Conn)
}

// HandlerOptions - All options available for a connection
type HandlerOptions struct {
	Addr  string
	Chain *Chain
	Users []*url.Userinfo
	// Authenticator Authenticator
	TLSConfig *tls.Config
	// Whitelist     *Permissions
	// Blacklist     *Permissions
	// Strategy      Strategy
	MaxFails    int
	FailTimeout time.Duration
	// Bypass        *Bypass
	Retries int
	Timeout time.Duration
	// Resolver      Resolver
	// Hosts         *Hosts
	ProbeResist  string
	KnockingHost string
	// Node          Node
	Host    string
	IPs     []string
	TCPMode bool
	// IPRoutes      []IPRoute
}

// HandlerOption - A way to set handler options
type HandlerOption func(opts *HandlerOptions)
