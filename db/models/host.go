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

import dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"

// SELECT ----------------------------------------------------------------------

// Hosts - Given a host model, this function checks all fields and builds a
// query that may match one or more hosts, and returns them.
// Some fields will automatically return only one host, such as ID
func Hosts(host dbpb.Host) (hosts []dbpb.Host) {
	return
}

// HostsByWorkspace - Returns all hosts for a given workspace
func HostsByWorkspace(id uint32) (hosts []dbpb.Host) {
	return
}

// HostByID - Return a host identified with an ID
func HostByID(id uint32) (host dbpb.Host) {
	return
}

// HostByAddress - Returns a host EXACTLY matching the provided IP (v4 or v6)
func HostByAddress(addr string) (host dbpb.Host) {
	return
}

// CREATE ----------------------------------------------------------------------

// AddHost - Add a host to the database, if no existing host matches it.
func AddHost(host dbpb.Host) (added dbpb.Host) {
	// We must find a few fields that, in common, should always produce
	// a unique fingerprint for a Host. Then we know if we need to update
	// an already saved host, or add a new one.
	// - Workspace ID
	// - MAC Address
	// - OS Name & arch
	// - List of IP Addresses is identical
	// - Host name
	// - List of users is identical
	return
}

// UPDATE ----------------------------------------------------------------------

// UpdateHost - Update a host provided with an ID (we thus know it already exists).
func UpdateHost(host dbpb.Host) (updated dbpb.Host) {
	return
}

// DELETE ----------------------------------------------------------------------

// DeleteHosts - Given a host model, this function checks all fields and builds a
// query that may match one or more hosts, and deletes the ones found.
// For instance, if the host has an ID, this one only will be deleted, but if it has
// several addresses, all matches will be deleted at once.
func DeleteHosts(host dbpb.Host) (deleted []dbpb.Host) {
	return
}

// DeleteHostsInWorkspace - Delete all hosts in workspace
func DeleteHostsInWorkspace(id uint32) (deleted int) {
	return
}

// DeleteHostsByIP - Delete hosts matching the IP
func DeleteHostsByIP(addr string) (deleted int) {
	return
}

// DeleteHostByID - Delete a host given its ID
func DeleteHostByID(id uint32) (deleted bool) {
	return
}
