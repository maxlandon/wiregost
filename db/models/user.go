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

import serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"

// SELECT ----------------------------------------------------------------------

// Users - Given a user model, this function checks all fields and builds a
// query that may match one or more users, and returns them.
// Some fields will automatically return only one user, such as ID
func Users(user *serverpb.User) (users []serverpb.User) {
	return
}

// UserByID - Return a user identified with an ID
func UserByID(id uint32) (user serverpb.User) {
	return
}

// CREATE ----------------------------------------------------------------------

// AddUser - Add a user to the database, if no existing user matches it.
func AddUser(user *serverpb.User) (added serverpb.User) {
	return
}

// UPDATE ----------------------------------------------------------------------

// UpdateUser - Update a user provided with an ID (we thus know it already exists).
func UpdateUser(user *serverpb.User) (updated serverpb.User) {
	return
}

// DELETE ----------------------------------------------------------------------

// DeleteUsers - Given a user model, this function checks all fields and builds a
// query that may match one or more users, and deletes the ones found.
// For instance, if the user has an ID, this one only will be deleted, but if it has
// several addresses, all matches will be deleted at once.
func DeleteUsers(user *serverpb.User) (deleted []serverpb.User) {
	return
}

// DeleteUserByID - Delete a user given its ID
func DeleteUserByID(id uint32) (deleted bool) {
	return
}
