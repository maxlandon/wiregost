package core

import "io"

// Interactive - This type implements the core stubs that provide an interactive session.
type Interactive struct {
	// Session - Base session information, route and logging.
	*Session

	// IO stream for this Session
	stream io.ReadWriteCloser
}

// NewInteractive - Providing an appropriate data stream, instantiate a new interactive session.
func NewInteractive(stream io.ReadWriteCloser) (s *Interactive) {

	s = &Interactive{
		New(),  // New base session object
		stream, // IO stream passed to this session
	}
	s.Interactive = true

	return
}
