package core

import "io"

// BasicIO - A base session type that only implements basic input/output
type BasicIO struct {
	// Session - Base session information, route and logging.
	*Session

	// IO stream for this Session
	IO *io.ReadWriteCloser
}

// NewBasicIO - Instantiates a new basic Session with primitive I/O capabilities.
func NewBasicIO() (s *BasicIO) {
	s = &BasicIO{
		New(), // Session
		nil,   // This IO is instantiated when passed to a net.Conn
	}
	return
}
