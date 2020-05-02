// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package models

import routepb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport/route"

// SELECT ----------------------------------------------------------------------

// Routes - Given a route model, this function checks all fields and builds a
// query that may match one or more routes, and returns them.
// Some fields will automatically return only one route, such as ID
func Routes(route routepb.Route) (routes []routepb.Route) {
	return
}

// RoutesByWorkspace - Returns all routes for a given workspace
func RoutesByWorkspace(id uint32) (routes []routepb.Route) {
	return
}

// RouteByID - Return a route identified with an ID
func RouteByID(id uint32) (route routepb.Route) {
	return
}

// RouteByAddress - Returns a route EXACTLY matching the provided IP (v4 or v6)
func RouteByAddress(addr string) (route routepb.Route) {
	return
}

// CREATE ----------------------------------------------------------------------

// AddRoute - Add a route to the database, if no existing route matches it.
func AddRoute(route routepb.Route) (added routepb.Route) {
	// We must find a few fields that, in common, should always produce
	// a unique fingerprint for a Route. Then we know if we need to update
	// an already saved route, or add a new one.
	// - Workspace ID
	// - MAC Address
	// - OS Name & arch
	// - List of IP Addresses is identical
	// - Route name
	// - List of users is identical
	return
}

// UPDATE ----------------------------------------------------------------------

// UpdateRoute - Update a route provided with an ID (we thus know it already exists).
func UpdateRoute(route routepb.Route) (updated routepb.Route) {
	return
}

// DELETE ----------------------------------------------------------------------

// DeleteRoutes - Given a route model, this function checks all fields and builds a
// query that may match one or more routes, and deletes the ones found.
// For instance, if the route has an ID, this one only will be deleted, but if it has
// several addresses, all matches will be deleted at once.
func DeleteRoutes(route routepb.Route) (deleted []routepb.Route) {
	return
}

// DeleteRoutesInWorkspace - Delete all routes in workspace
func DeleteRoutesInWorkspace(id uint32) (deleted int) {
	return
}

// DeleteRoutesByIP - Delete routes matching the IP
func DeleteRoutesByIP(addr string) (deleted int) {
	return
}

// DeleteRouteByID - Delete a route given its ID
func DeleteRouteByID(id uint32) (deleted bool) {
	return
}
