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

import transportpb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport"

// SELECT ----------------------------------------------------------------------

// Listeners - Given a listener model, this function checks all fields and builds a
// query that may match one or more listeners, and returns them.
// Some fields will automatically return only one listener, such as ID
func Listeners(listener transportpb.Listener) (listeners []transportpb.Listener) {
	return
}

// ListenersByWorkspace - Returns all listeners for a given workspace
func ListenersByWorkspace(id uint32) (listeners []transportpb.Listener) {
	return
}

// ListenerByID - Return a listener identified with an ID
func ListenerByID(id uint32) (listener transportpb.Listener) {
	return
}

// ListenerByAddress - Returns a listener EXACTLY matching the provided IP (v4 or v6)
func ListenerByAddress(addr string) (listener transportpb.Listener) {
	return
}

// CREATE ----------------------------------------------------------------------

// AddListener - Add a listener to the database, if no existing listener matches it.
func AddListener(listener transportpb.Listener) (added transportpb.Listener) {
	// We must find a few fields that, in common, should always produce
	// a unique fingerprint for a Listener. Then we know if we need to update
	// an already saved listener, or add a new one.
	// - Workspace ID
	// - MAC Address
	// - OS Name & arch
	// - List of IP Addresses is identical
	// - Listener name
	// - List of users is identical
	return
}

// UPDATE ----------------------------------------------------------------------

// UpdateListener - Update a listener provided with an ID (we thus know it already exists).
func UpdateListener(listener transportpb.Listener) (updated transportpb.Listener) {
	return
}

// DELETE ----------------------------------------------------------------------

// DeleteListeners - Given a listener model, this function checks all fields and builds a
// query that may match one or more listeners, and deletes the ones found.
// For instance, if the listener has an ID, this one only will be deleted, but if it has
// several addresses, all matches will be deleted at once.
func DeleteListeners(listener transportpb.Listener) (deleted []transportpb.Listener) {
	return
}

// DeleteListenersInWorkspace - Delete all listeners in workspace
func DeleteListenersInWorkspace(id uint32) (deleted int) {
	return
}

// DeleteListenersByIP - Delete listeners matching the IP
func DeleteListenersByIP(addr string) (deleted int) {
	return
}

// DeleteListenerByID - Delete a listener given its ID
func DeleteListenerByID(id uint32) (deleted bool) {
	return
}
