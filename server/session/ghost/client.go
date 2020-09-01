package ghost

import (
	"io"
	"sync"

	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"

	"github.com/maxlandon/wiregost/server/session/core"
	"github.com/maxlandon/wiregost/server/transport"
)

// Client - A session.ghost.Client is the client-side instance of a remote Ghost implant.
// This object in charge of implementing the most fundamental aspects of a Ghost session
// such as information, basic RPC functionality, concurrency management, extensions, etc.
//
// NOTE: It is important to signal that this object IS OS/ARCH AGNOSTIC: everything provided
// by this object can work any where, with all the transports accepted (and thus maybe held)
// by this Client ghost (and its remote peer).
type Client struct {
	// Interactive: The client embbeds reference to an Interactive session: the
	// ReadWriteCloser stream that it provides is used as the "main" channel of
	// this implant. We use it for ALL requests pertaining to the implant's base
	// management, and Transports are in charge of allowing this.
	// Eventually we might request to open new channels.
	*core.Interactive

	// Multiplexer - This object holds a multiplexable connection that has been drawned out of the
	// Interactive stream: it is used by the implant channels for core functionality and by the
	// routing system for all other non-ghost traffic.
	Multiplexer *yamux.Session

	// Log - This implant already beneficiates from a log with the Interactive base,
	// however we might modify it permanently, either to add more detailed logging
	// policy or output sources.
	Log *logrus.Entry

	// Transports - All transports currently saved and usable by the ghost implant
	// have an equivalent object in the client. These modules may pass net.Conn-like
	// objects to the Session's Interactive stream, when transports are swapped/used/etc.
	//
	// The way transports are leveraged through this client is important, and the
	// set of aims and constraints to solve is large: sending requests to implants
	// must be very simple, because transports should deal with routing abstractions.
	// Meanwhile, each of these modules handles a different transport+connect protocol,
	// they might have to make lots of adjustements/implementation in order to seamlessly
	// interface with this client RPC.
	Transports *map[string]transport.Transport

	// Channels - All implants can have concurrent processes/routines/tasks executed remotely,
	// such as pesudo-command shells, routers, etc. All of these channels have a client-side
	// object that holds its own state, history, etc.
	Channels map[string]*Channel

	// RWMutex - This client has to manage many different processes and tasks concurrently
	// (core functions, routing requests, transports, etc.), and needs maximum safety. We
	// use a RWMutex because it might give more fine-grained control with better performance.
	mutex *sync.RWMutex
}

// New - Instantiates a new Ghost client object, which is also the equivalent of a Session creation.
// The object is also instantantied with a stream, passed to the Interactive base. We then build on this.
func New(stream io.ReadWriteCloser) (c *Client, err error) {

	c = &Client{
		core.NewInteractive(stream),       // We can now send over the wire.
		nil,                               // We setup the log later.
		nil,                               // We setup the ChannelStream just below.
		&map[string]transport.Transport{}, // The transport is empty for now.
		map[string]*Channel{},             // There are no side-channels at startup.
		&sync.RWMutex{},                   // Concurrent access
	}

	// Setup Multiplexer, we return for any error as we might still set it up later
	if err = c.SetupMultiplexer(c.Stream); err != nil {
		return c, err
	}

	return
}

// SetupMultiplexer - The client is being passed a multiplexable stream, "wraps" it
// around a multiplexer: it acts as a pool of streams for implant channels and routing system.
func (c *Client) SetupMultiplexer(stream io.ReadWriteCloser) (err error) {

	// Here we might have to setup a precise yamux config for optimization.
	c.Multiplexer, err = yamux.Client(stream, yamux.DefaultConfig())

	return
}

// SetupLog - The Client object reimplements the SetupLog() of its Interactive base,
// which more elaborated logging policies and output files.
func (c *Client) SetupLog() (err error) {

	// We start by calling the Interactive.SetupLog() func, which fixes a default
	// sane logger based on the Session owner, sets basic fields like UUIDs, etc.
	// We can now either access this log as is, or modify it and reassign it to
	// the Client's own Log field.
	c.Log, err = c.Interactive.SetupLog()

	return
}

// Kill - Kill the remote Ghost implant, with all the checkups needed. This function
// should be used in conjunction with sessions.Deregister() and below's Cleanup()
func (c *Client) Kill() (err error) {
	// Metasploit does "deregestering" here.
	return
}

// Cleanup - A ghost implant has more elaborated logic and state. The Cleanup() reimplements
// the Interactive Cleanup() and adds various things like checking UUID persistence, log files,
// active routes and traffic they handle, etc. This function should work 'de concert' with the
// sessions.Deregister() function.
func (c *Client) Cleanup() (err error) {
	return
}
