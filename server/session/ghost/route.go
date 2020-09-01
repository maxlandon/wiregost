package ghost

import (
	routepb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport/route"
)

// AddRoute - We add a route chain to this client, for connecting to remote networks.
// This has no effect on the currently existing RouteStream and its conns, at best it
// will open a new physical connection/listener on the remote side.
func (c *Client) AddRoute(r *routepb.Route) (err error) {

	// Check the implant supports this route, and permissions are ok.

	// Send request to implant: it will return success with much info, or error

	// Add to client route list, with the returned id, etc.

	return
}

// RemoveRoute - We remove a route chain from this client.
func (c *Client) RemoveRoute(r *routepb.Route) (err error) {

	// Check route exists in client list

	// Send request to implant. Check result has we might have to wait for conns to end

	// Handle either wait time (with a log event, for instance), or return

	// If we went here we have cut the route, so we delete it from the client list.

	return
}
