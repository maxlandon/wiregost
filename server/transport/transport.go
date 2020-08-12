package transport

import (
	modpb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport"
	"github.com/maxlandon/wiregost/server/module"
	"github.com/maxlandon/wiregost/server/module/stack"
)

// Module - A transport is a module in Wiregost because a console may use them in several contexts:
//
// EXPLOITS --------------------
// 1) Associated with an Exploit module: for remote exploits, we need a listener/dialer with
//    a capability either to register a Ghost implant, or to stage it, or even simply to handle
//    a non-Ghost session type, which still may be used by consoles (like a remote SSH session)
//
// GHOSTS ----------------------
// 2) Ghost implants may have several transports to use/cycle through. Users can add some by
//    using some transport modules. Generally, these transport modules would be single, because
//    the "stage" is already up and running. However, transports might be configured to stage
//    a new ghost, for whatever reason.
//
// 3) When a user wants to compile a Ghost implant, he uses the appropriate payload module to
//    prepare the implant. He may use a transport module to add configuration before compilation.
//
// NOTE: Routing Abstractions
// All types (and therefore, modules) embedding this type do so because it provides methods
// do deal with abstractions like routing (which actual physical/logical connection to use?),
// or like the concerned host (do we start a listener on the server or on an implant?), etc.
//
// Maybe this base module might include a net.Conn/net.Listener/net.Dialer, used by subtypes
// for interacting with the other end of the transport.
type Module struct {
	// Base module. Makes this Exploit a valid module in Wiregost, with full access to UI.
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

// Run - Execute the main function of this transport, which may depend on which
// settings and details are provided, such as: is the transport a Dialer or a Listener ?
// or: is the Transport set for a remote implant ? or is the Transport has to be started
// now or later ?
func (m *Module) Run(cmd string, args []string) (result string, err error) {

	return
}

// Start - Start monitoring a logical/physical connection.
func (m *Module) StartHandler() (err error) {
	return
}

// AddHandler - Adds another connection monitor.
func (m *Module) AddHandler() (err error) {
	return
}

// Stop - stop monitoring a logical/physical connection.
func (m *Module) StopHandler() (err error) {
	return
}

// HandleConnection - Handles an established (logical/physical connection). The default
// path is to attempt to create a Session, but it will be overriden by some subtypes.
func (m *Module) HandleConnection() (err error) {
	// Create Session
	return
}

// Waits for a session to be created as the result of a handler connection coming in.
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
func (t *Module) OnSession() (err error) {

	// If there is an associated exploit, notify him so that he can do
	// his things if he needs to.
	return
}

// CreateSession - Creates a session, if necessary, for the connection handled.
func (t *Module) CreateSession() (err error) {

	return
}

// Cleanup - Clean any state needed for this module. This function is here more to remind
// all types embedding this module that they may override it, as a good practice of cleaning.
func (m *Module) Cleanup() (err error) {

	return
}
