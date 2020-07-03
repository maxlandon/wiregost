package route

import "net"

// Handler - A proxy server handler. Configures and runs the connection.
type Handler interface {
	Init(options ...HandlerOption)
	Handle(net.Conn)
}

// HandlerOptions - All options available for a connection
type HandlerOptions struct{}

// HandlerOption - A way to set handler options
type HandlerOption func(opts *HandlerOptions)
