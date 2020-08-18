package session

import (
	"sync"
)

var (
	// Sessions - Holds all registered Sessions in the lifetime of this server. This package
	// (root of session/ packages) is meant to manage the complete lifetime of a Session no
	// matter its type, target host, transports or abilities.
	//
	// As well, it is from here that many pieces of the system will be fired up and organized,
	// because Sessions will be accessed by some way or another for most of Wiregost abilities:
	// - Transports needing to operate behind enemy lines.
	// - Post-exploitation modules triggering RPC calls
	// - Console users interacting with sessions.
	// Depending on various checks, Sessions are "registered" to all the interfaces they implement.
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

	// Metasploit sends here an event (should be detailed) notifying we are opening a Session.
	// If populated correctly, exploits, transports, payloads and routes might be alerted for
	// certain things to do, and others to not do...

	return
}

// Deregister - Deregister the Session identified by ID. This function will obviously
// call the last embeder's implementation of Cleanup, etc...
func (s *sessions) Deregister(sessionUUID string) (ok bool, err error) {

	// Check state of session (lost or down ?)

	// The Session should unsubscribe from all Events

	// Notify everyone were are deregistering session. We send detailed identification. Now:
	// We might give either some time to everyone to stop their business with this session,
	// if they have any.
	// Or we might call various functions on objects, concurrently or sequentially, to verify
	// from HERE that the session has no debts to anyone around.

	// Routing System:
	// It's another part that will have to handle and synchronize many things when an event
	// occurs, especially when sessions pop/quit: routes and pivots to verify, routed traffic.
	// Here we might have to do some go-back-and-forth with the Routing system, to handle
	// everything nicely.

	return
}

// Function for concurrently monitoring the state of sessions, and taking according steps
// (logging, automatic actions, user messages/input, etc.).
