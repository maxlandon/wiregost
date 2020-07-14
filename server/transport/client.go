package transport

import (
	"net"

	"github.com/maxlandon/wiregost/server/c2/connector"
	"github.com/maxlandon/wiregost/server/c2/transporter"
)

var (
// DefaultClient is a standard HTTP proxy client.
// DefaultClient = &Client{Connector: HTTPConnector(nil), Transporter: TCPTransporter()}
)

// Client - A proxy client, divided between two layers: connector and transporter.
type Client struct {
	Connector   Connector
	Transporter Transporter
}

// Connector is responsible for connecting to the destination address through this proxy.
type Connector interface {
	Connect(conn net.Conn, addr string, options ...connector.ConnectOption) (net.Conn, error)
}

// Transporter performs a handshake with this proxy.
type Transporter interface {
	Dial(addr string, options ...transporter.DialOption) (net.Conn, error)
	Handshake(conn net.Conn, options ...transporter.HandshakeOption) (net.Conn, error)
	Multiplex() bool // Indicate that the Transporter supports multiplex
}

// Dial - Connects to the target address.
func (c *Client) Dial(addr string, options ...transporter.DialOption) (net.Conn, error) {
	return c.Transporter.Dial(addr, options...)
}

// Handshake - Performs a handshake with the proxy over connection conn.
func (c *Client) Handshake(conn net.Conn, options ...transporter.HandshakeOption) (net.Conn, error) {
	return c.Transporter.Handshake(conn, options...)
}

// Connect - Connects to the address via the proxy over connection conn
func (c *Client) Connect(conn net.Conn, addr string, options ...connector.ConnectOption) (net.Conn, error) {
	return c.Connector.Connect(conn, addr, options...)
}
