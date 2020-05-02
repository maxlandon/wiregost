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

// Services - Returns a list of services according to filters provided in the service template
func Services(service dbpb.Service) (found []dbpb.Service) {
	return
}

// CREATE ----------------------------------------------------------------------

// AddService - Add a service to the database, if no existing services match it.
func AddService(service dbpb.Service) (added dbpb.Service) {
	return
}

// UPDATE ----------------------------------------------------------------------

// UpdateService - Update a Service provided with an ID (we thus know it already exists).
func UpdateService(service dbpb.Service) (updated dbpb.Service) {
	return
}

// DELETE ----------------------------------------------------------------------

// DeleteServices - Given a service model, this function checks all fields and builds a
// query that may match one or more services, and deletes the ones found.
// For instance, if the service has an ID, this one only will be deleted, but if it has
// several addresses, all matches will be deleted at once.
func DeleteServices(service dbpb.Service) (deleted []dbpb.Service) {
	return
}

// DeleteServiceByID - Delete a service given its ID
func DeleteServiceByID(id uint32) (deleted bool) {
	return
}
