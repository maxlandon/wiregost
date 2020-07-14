package route

import "net"

// Accepter represents a network endpoint that can accept connection from peer.
// This is, a priori, where will direct our ghost connections, which then will be routed.
type Accepter interface {
	Accept() (net.Conn, error)
}

// Server - A proxy server
type Server struct {
	Listener Listener
	Handler  Handler
	Options  *ServerOptions
}

// Init intializes server with given options.
func (s *Server) Init(opts ...ServerOption) {
	if s.Options == nil {
		s.Options = &ServerOptions{}
	}
	for _, opt := range opts {
		opt(s.Options)
	}
}

// Close closes the server
func (s *Server) Close() error {
	return s.Listener.Close()
}

// Serve serves as a proxy server.
func (s *Server) Serve(h Handler, opts ...ServerOption) error {
	return nil
}

// Run starts to serve.
func (s *Server) Run() error {
	return s.Serve(s.Handler)
}

// ServerOptions - Holds the options for Server.
type ServerOptions struct {
}

// ServerOption - Used to set server options.
type ServerOption func(opts *ServerOptions)
