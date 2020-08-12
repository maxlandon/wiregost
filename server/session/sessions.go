package session

import (
	"sync"

	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

// Session - Plays a role similar to the module.Module interface, which allows us
// to load various types of modules on a same stack. Sessions might be of several
// subtypes, and therefore with very different skill sets.
// Thus the Session interface aims to remain relatively small, or at least being
// implemented by ALL session types.
//
// NOTE: In order to use the extended functionality of all Sessions, we will make
// use of other, skilled-specific interfaces in other packages.
type Session interface {
	// ToProtobuf - The Session can push all of its base information .
	ToProtobuf() (sess *serverpb.Session)
}

var (
	// Sessions - Holds all registered Sessions in the lifetime of this server.
	Sessions = &sessions{
		Registered: []Session{},
		mutex:      &sync.Mutex{},
	}
)

type sessions struct {
	Registered []Session
	mutex      *sync.Mutex
}

// Register - Registers the supplied Session object to the server and returns a UUID to the caller
func (s *sessions) Register(sess Session) (id string, err error) {
	return
}

// Deregister - Deregister the Session identified by ID. This function will obviously
// call the last embeder's implementation of Cleanup, etc...
func (s *sessions) Deregister(sessionUUID string) (ok bool, err error) {
	return
}

// Function for concurrently monitoring the state of sessions, and taking according steps
// (logging, automatic actions, user messages/input, etc.).
