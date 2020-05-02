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

// Workspaces - Returns a list of workspaces according to filters provided in the workspace template
func Workspaces(workspace dbpb.Workspace) (found []dbpb.Workspace) {
	return
}

// CREATE ----------------------------------------------------------------------

// AddWorkspace - Add a workspace to the database, if no existing workspaces match it.
func AddWorkspace(workspace dbpb.Workspace) (added dbpb.Workspace) {
	return
}

// UPDATE ----------------------------------------------------------------------

// UpdateWorkspace - Update a Workspace provided with an ID (we thus know it already exists).
func UpdateWorkspace(workspace dbpb.Workspace) (updated dbpb.Workspace) {
	return
}

// DELETE ----------------------------------------------------------------------

// DeleteWorkspaces - Given a workspace model, this function checks all fields and builds a
// query that may match one or more workspaces, and deletes the ones found.
// For instance, if the workspace has an ID, this one only will be deleted, but if it has
// several addresses, all matches will be deleted at once.
func DeleteWorkspaces(workspace dbpb.Workspace) (deleted []dbpb.Workspace) {
	return
}

// DeleteWorkspaceByID - Delete a workspace given its ID
func DeleteWorkspaceByID(id uint32) (deleted bool) {
	return
}
