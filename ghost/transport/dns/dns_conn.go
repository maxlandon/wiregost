package dns

import (
	"net"
	"sync"
	"time"
)

// Conn - A custom DNS Connection that implements net.Conn. However, it has some
// fundamental differences:
// - Unlike a TCP conn that reads bytes until EOF is reached, the dns.Conn object
// reads bytes until it has internally delimited the size of the complete DNS message.
// In fact, this Conn object uses net.LookupTXT twice, which means it makes use of
// net.Conn behind the scenes.
type Conn struct {
	// Fields from Sliver
	Send    chan []byte // Sliver uses pb.Envelope, we go low-level
	Recv    chan []byte
	IsOpen  bool
	ctrl    chan bool
	cleanup func()
	once    *sync.Once
	mutex   *sync.RWMutex
	// tunnels *map[uint64]*Tunnel

	// Wiregost added fields
}

// Implementation of the Conn interface.

// Read implements the Conn Read method.
func (c *Conn) Read(b []byte) (int, error) {
	return 0, nil
}

// Write implements the Conn Write method.
func (c *Conn) Write(b []byte) (int, error) {
	return 0, nil
}

// Close closes the Connection.
func (c *Conn) Close() error {
	return nil
}

// LocalAddr returns the local network address.
// The Addr returned is shared by all invocations of LocalAddr, so
// do not modify it.
func (c *Conn) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr returns the remote network address.
// The Addr returned is shared by all invocations of RemoteAddr, so
// do not modify it.
func (c *Conn) RemoteAddr() net.Addr {
	return nil
}

// SetDeadline implements the Conn SetDeadline method.
func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline implements the Conn SetReadDeadline method.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}
