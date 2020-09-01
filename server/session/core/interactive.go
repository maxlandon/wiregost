package core

import (
	"bufio"
	"io"
	"time"
)

// Interactive - This type implements the core stubs that provide an interactive session.
// Generally, and given Go's std library straightforward character, this interactive
// session can plug itself on any ReadWriteCloser (thus any network connection), and
// will provide concurrent read/write operations on it.
// However, these stubs will only support primary encoding/serializing, because we still don't
// know how and by what the data we send through the tunnel is processed: therefore, basic stuff.
//
// Latter, some embeders of this Interactive type will implement more elaborated
// means of communication on top of this, like Ghost sessions, which will add
// protobuf serialization and handler matching.
//
// NOTE: Think about adding any basic encryption/security things in this Session if needed.
type Interactive struct {
	// Session - Base session information, route and logging.
	*Session

	// Stream - I/0 Stream for this session. This can be anything: net.Conn or
	// a mux.Stream passed by a handler/route. It can be the OS Stdin/Stdout/Stderr
	// of a host, etc. Reader, Writer and Closer may even point to different streams.
	Stream  io.ReadWriteCloser
	reader  *bufio.Reader // Fine-grained control over how to read the stream.
	timeout time.Duration // Default timeout when not set from user

}

// NewInteractive - Providing an appropriate data stream (which can be a
// net.Conn, a mux.Stream, etc.), instantiate a new interactive session.
// The session is ready to read and write to the stream with its reader.
func NewInteractive(stream io.ReadWriteCloser) (s *Interactive) {

	s = &Interactive{
		New(),            // New base session object
		stream,           // IO stream passed to this session
		nil,              // Instantiate the reader below, in Setup().
		10 * time.Second, // Default timeout
	}
	s.reader = bufio.NewReader(s.Stream) // Initialize conn reader
	s.Interactive = true

	return
}

// Cleanup - Clean any state related to this Interactive Session.
// Should call the *Session implementation at some point.
func (i *Interactive) Cleanup() (err error) {
	return
}

// Kill - Terminate the Interactive session. Cleans up the resources and
// calls the *Session Kill() implementation for deregistering the Session.
func (i *Interactive) Kill() (err error) {

	// This involves handling the way we kill the ReadWriteCloser.
	// The issue here is that we don't know anything about
	return
}

// Suspend - The user has sent a kill signal to the stream. The issue here is that
// this object has no idea of what the underlying communication medium is: it might
// be a muxed conn that is routed through 5 pivots, the last one only handling the
// physical connection.
// We thus don't yet if this function will be used at the Interactive level.
//
// NOTE: There might be a bit of code to do for handling all cases, here and there
// in the framework, and on the ghost implants as well (listeners, etc.). For instance
// boiler plate code could be necessary if we are interacting with a very primitive shell,
// like netcat for instance, which by default does not basic escapes. The Go program has
// access to syscall variables.
func (i *Interactive) Suspend() (err error) {

	return
}
