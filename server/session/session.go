package session

import (
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/sirupsen/logrus"
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
	// Info - The Session can push all of its base information.
	// This should encapsulate most of the data concerning the core Session.
	// Non-core info, like routes, transports, etc, should be handled by more
	// specialized Session types, if any.
	Info() (sess *serverpb.Session)
	// Setup - The Session can setup itself, no matter its type.
	Setup() (err error)
	// SetupLog - All Session must be able to set their own log depending
	// on various things, such as session owners, workspaces, skills, etc.
	// There some cases where implants might need to log to several files.
	// The function can return a log instance in case its needed during
	// registration, etc.
	SetupLog() (log *logrus.Entry, err error)

	// ToProtobuf - The Session can push all of its base information .
	ToProtobuf() (sess *serverpb.Session)
}
