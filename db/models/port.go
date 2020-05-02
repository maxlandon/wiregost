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

// Ports - Returns a list of ports according to filters provided in the port template
func Ports(port dbpb.Port) (found []dbpb.Port) {
	return
}

// CREATE ----------------------------------------------------------------------

// AddPort - Add a port to the database, if no existing ports match it.
func AddPort(port dbpb.Port) (added dbpb.Port) {
	return
}

// UPDATE ----------------------------------------------------------------------

// UpdatePort - Update a Port provided with an ID (we thus know it already exists).
func UpdatePort(port dbpb.Port) (updated dbpb.Port) {
	return
}

// DELETE ----------------------------------------------------------------------

// DeletePorts - Given a port model, this function checks all fields and builds a
// query that may match one or more ports, and deletes the ones found.
// For instance, if the port has an ID, this one only will be deleted, but if it has
// several addresses, all matches will be deleted at once.
func DeletePorts(port dbpb.Port) (deleted []dbpb.Port) {
	return
}

// DeletePortByID - Delete a port given its ID
func DeletePortByID(id uint32) (deleted bool) {
	return
}
