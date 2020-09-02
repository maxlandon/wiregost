package ghost

import (
	"io"

	"github.com/hashicorp/yamux"
	transportpb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport"
)

// AddTransport - Instructs the implant to add a transport to its list, if possible.
func (c *Client) AddTransport(t *transportpb.Transport) (err error) {

	return
}

// RemoveTransport - Instructs the implant to remove a transport from its list.
func (c *Client) RemoveTransport(t *transportpb.Transport) (err error) {
	return
}

// GetSupportedTransports - Get the list of all transport protocols supported by this implant
func (c *Client) GetSupportedTransports() (ts []transportpb.Transport) {
	return
}

// HandleRouteMux - This function is an equivalent of the HandleRouteMux() goroutine used by implants for their active transports.
// This will be called when, server-side, we receive a request from an implant that must forward traffic back to us. We thus start
// to listener for its incoming mux requests.
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
			// 1) Streams used by pivoted implants to pass core C2 requests/responses
			// 2) All non-implant traffic going through, no matter what happens before/after.

			// Add stream to the streams channel: these streams will be handled by the implant's
			// routing system and its handlers.
			streams <- stream
		}
	}
}
