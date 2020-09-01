package ghost

import (
	routepb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport/route"
)

// AddRoute - We add a route chain to this client, for connecting to remote networks.
// This has no effect on the currently existing Multiplexer and its generated conns,
// at best it will open a new physical connection/listener on the remote side.
func (c *Client) AddRoute(r *routepb.Route) (err error) {

	// Check the implant supports this route, and permissions are ok.
	// ERROR CASES:
	// - The current transport does not support the type of route asked, like a multiplexer for a DNS transport,
	// - The route requires permissions to physically open a port and these permissions are not set.

	// Send request to implant: it will return success with much info, or error
	// Here we might need to "fork" in the case where the opened route requires
	// a new physical connection opened: then we will not use the implant Multiplexer
	// when going through this route.

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
