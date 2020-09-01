package transport

import (
	"io"

	modpb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport"
	"github.com/maxlandon/wiregost/server/module"
	"github.com/maxlandon/wiregost/server/module/stack"
)

// Transport - A transport object able to establish a connection between 2 nodes.
// This object will primarily be used by Payload modules and Session objects, and is never used
// by console users: the Module type below is used for this instead.
type Transport interface {
	// Base methods
	Type() pb.Type                                    // Is this transport a bind or reverse type ?
	Info() *pb.Transport                              // All information pertaining to this transport
	Multiplex() bool                                  // Is this transport able to multiplex connections ?
	Start() (err error)                               // Either dial to or listen on an address, with current options
	HandleConn() (conn io.ReadWriteCloser, err error) // Same as start, but return the raw conn/stream
	Stop() error                                      // Stop the current transport.
	WaitForSession() error                            // Waits for a session to be created as the result of a handler connection coming in.
	CreateSession() error                             // Creates a session, if necessary, for the connection handled.
	ToCompileConfig() (conf string, err error)        // Get a string we will inject in an implant source code
}

// Module - A transport is a module in Wiregost because a console may use them in several contexts:
//
// 1) Associated with an Exploit module: for remote exploits, we need a listener/dialer with
//    a capability either to register a Ghost implant, or to stage it, or even simply to handle
//    a non-Ghost session type, which still may be used by consoles (like a remote SSH session)
//
// 2) Ghost implants may have several transports to use/cycle through. Users can add some by
//    using some transport modules. Generally, these transport modules would be single, because
//    the "stage" is already up and running. However, transports might be configured to stage
//    a new ghost, for whatever reason.
//
// 3) When a user wants to compile a Ghost implant, he uses the appropriate payload module to
//    prepare the implant. He may use a transport module to add configuration before compilation.
type Module struct {
	// Base module. Makes this Transport a valid module in Wiregost, with full access to UI.
	*module.Base

	// Base information for this transport: check if needed here.
	*pb.Transport
}

// New - Instantiates a new Transport module. This function is called
// by Exploit and Payload modules as well as console users.
func New(meta *modpb.Info) (m *Module) {

	m = &Module{
		module.New(meta), // Create base module
		nil,              // Don't know if we're using this yet.
	}

	// Defaul module settings
	m.Info.Type = modpb.Type_TRANSPORT    // Set module type
	m.StagingType = pb.StagingType_SINGLE // The transport is single by default
	m.WFSDelay = int32(5)                 // 5 seconds is the default for session wait time upon a connection.

	// Add specific fields to the Transport logger. Overwrites "module":"module" key/val pair.
	m.Log = m.Log.WithField("module", "transport")

	// Default options and commands
	m.AddOption("HandlerName", "", "", "A human name for this handler", false)
	return
}

// AddModule - Implements the stack.Module interface. This Transport will always return
// false, because until needed otherwise, we don't embed modules into Transports.
func (m *Module) AddModule(mod stack.Module) (ok bool, err error) {
	// Do nothing and change nothing.
	return
}

// SetOption - Implements the stack.Module interface. The Transport does not
// embed any submodule so there is no special logic needed here.
func (m *Module) SetOption(opt *modpb.Option) (err error) {

	if option, found := m.Opts[opt.Name]; found {
		err = option.Set(opt.Value)
	}
	return
}

// Run - Execute the main function of this transport, which may depend on which settings and details are provided, such as:
// Is the transport a Dialer or a Listener ?
// Is the Transport set for a remote implant ?
// Is the Transport has to be started now or later ?
// The module takes care of requesting all details from its Transport object, checks them, and acts accordingly
func (m *Module) Run(cmd string, args []string) (result string, err error) {

	// 1) Is the transport to be started remotely ? Check routing table

	// If yes {}
	// Get the route chain, and for each ghost ID, request the appropriate client to mux its physical conn.
	// The last node has started the appropriate bind/reverse handler, and returned the conn to the n-1 node.
	// We get here a logical conn returned by the first node in the chain

	return
}

// StartHandler - Start monitoring a logical/physical connection.
func (m *Module) StartHandler() (err error) {
	return
}

// AddHandler - Adds another connection monitor.
func (m *Module) AddHandler() (err error) {
	return
}

// StopHandler - stop monitoring a logical/physical connection.
func (m *Module) StopHandler() (err error) {
	return
}

// HandleConnection - Handles an established (logical/physical connection). The default
// path is to attempt to create a Session, but it will be overriden by some subtypes.
// This might be used for MANY THINGS, but should generally return at least an io.ReadWriteCloser
func (m *Module) HandleConnection() (err error) {

	// Create Session:
	// 1) Instantiates a new session object, with a base stream passed on to it.

	// 2) Check Session type:
	// We can directly add some RPC registration boilerplate, with a timeout:
	// if the timeout is reached without receiving a registration message, or an error
	// we just stay with the core Interactive session.

	// If the registration is received, we fill everything needed for the Ghost session:
	// We might have to call a few functions to setup RPC for main channel, etc.

	// 3) Finally add the Session to Sessions (Sessions.Register(session))

	return
}

// WaitForSession - Waits for a session to be created as the result of a handler connection coming in.
// The return value is either a Session object, or nil if the timeout expires
func (m *Module) WaitForSession() (err error) {
	return
}

// OnSession - Equivalent to Metasploit's on_session function. Here is its description:
//
// "Once an exploit completes and a session has been created on behalf of the
// {transport}, the framework will call the {transport}'s on_session notification
// routine to allow it to manipulate the session prior to handing off
// control to the user."
//
// Because in Wiregost, the Transport is responsible for monitoring connections and Session
// management (registration first and foremost), this function is here.
func (m *Module) OnSession() (err error) {

	// If there is an associated exploit, notify him so that he can do
	// his things if he needs to.
	return
}

// CreateSession - Creates a session, if necessary, for the connection handled.
func (m *Module) CreateSession() (err error) {

	// Here, Metasploit asks a payload module if it has a "session factory":
	// This allows a payload module, previously used for compiling an implant
	// to also setup the client object for this payload.
	// The transport is thus agnostic to what type of session is created and
	// how to populate/set it up accordingly.

	// Set FromExploit if needed, or call it through the Session object which,
	// again, might know better than us how to link exploits and sessions.

	// Pass along various context like workspace

	// Here, Metasploit allows to pass custom UUIDs. Check if that's useful.
	return
}

// Cleanup - Clean any state needed for this module. This function is here more to remind
// all types embedding this module that they may override it, as a good practice of cleaning.
func (m *Module) Cleanup() (err error) {

	return
}
