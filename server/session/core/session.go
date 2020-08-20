package core

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/maxlandon/wiregost/ghost/transport"
	"github.com/maxlandon/wiregost/modules/templates/exploit"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/transport/route"
)

// Session - Similarly to Metasploit, a Session object is a general means of
// interacting with various post-exploitation payloads through a common interface
// that is not necessarily tied to a network connection.
//
// For instance, if an exploit spawns a command shell over the network, the read/write
// operations end up reading and writing to that shell. These raw functions will be
// progressively reimplemented by embedders when they need more elaborated logic.
type Session struct {
	// Base information for this Session, no matter its embedders.
	*serverpb.Session

	// Routes - A session might be accessible through many network routes.
	// Each of these routes has been added by a routing module, and they
	// have their own specifics, like reliabity and security criteria.
	//
	// At any point in time, for some types of Sessions, there might be
	// several active routes leading to a same pivot session.
	Routes []route.Chain

	// Environment - The Environment object stores environment variables
	// on the session host. The latter might be queried through simple
	// command one-liners and parsed, all that automatically at startup.

	// Logger
	Log *logrus.Entry
}

// New - A Transport module has handled a connection and creates a session.
func New() (s *Session) {
	s = &Session{
		&serverpb.Session{
			ID:   uuid.New().String(),
			Type: serverpb.SessionType_UNKNOWN,
			// Owner
			// Permissions
			// RoutePermissions
			Status: serverpb.Status_ALIVE,
		},
		nil, // No routes yet
		nil, // Logger is set up later
	}

	s.SetupLog()

	return
}

// Info - The Session can push all of its base information.
func (s *Session) Info() (sess *serverpb.Session) {
	return
}

// SetupLog - The session instantiates and setup its log, which all embedders can use.
func (s *Session) SetupLog() (log *logrus.Entry, err error) {
	// Many fields to pass in: session uuid, log files to set/get for later,
	// Command history file
	s.Log = logrus.StandardLogger().WithField("session", "test")

	return s.Log, nil
}

// FromExploit - When this session has been spawned from an Exploit module,
// we derive its information, just in case.
func (s *Session) FromExploit(m *exploit.Module) {

}

// FromTransport - Associates a transport to this session. This is useful
// when the session is a simple command shell listening on a port: we might
// have cut the connection but not the shell, etc. Anyway, we keep a reference
// to this transport, which will be reused if needed.
func (s *Session) FromTransport(t *transport.Transport) {

}

// Kill - At this point, this function just asks deletion of the Session from a list.
func (s *Session) Kill() (err error) {
	// Metasploit does "deregestering" here.
	return
}

// Cleanup - Clean any state needed for this Session. This function is here more to remind
// all types embedding this Session that they may override it, as a good practice of cleaning.
func (s *Session) Cleanup() (err error) {
	return
}

// ToProtobuf - The Session can push all of its base information .
func (s *Session) ToProtobuf() (sess *serverpb.Session) {
	return
}
