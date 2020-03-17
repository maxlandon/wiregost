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

import "time"

// Workspace is the exported object
type Workspace struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Boundary       string    `json:"boundary"`
	LimitToNetwork bool      `json:"limit_to_network"`
	IsDefault      bool      `json:"is_default"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Instantiate a new workspace, with a unique ID. Only used by AddWorkspaces().
func newWorkspace(name string) *Workspace {
	w := &Workspace{
		Name: name,
	}
	return w
}

// Workspaces returns all workspaces in database
func (db *DB) Workspaces() ([]*Workspace, error) {
	var workspaces []*Workspace

	err := db.Find(&workspaces)
	if len(err.GetErrors()) != 0 {
		return nil, err.GetErrors()[0]
	}

	return workspaces, nil
}

// AddWorkspaces adds workspaces to database, using names supplied.
func (db *DB) AddWorkspaces(names []string) error {
	for _, name := range names {
		workspace := newWorkspace(name)

		err := db.Create(workspace)
		if len(err.GetErrors()) != 0 {
			return err.GetErrors()[0]
		}
	}

	return nil
}

// DeleteWorkspaces adds workspaces to database, using ids supplied.
func (db *DB) DeleteWorkspaces(ids []uint) (rows int64, err error) {
	w := new(Workspace)
	var deleted int64

	for _, id := range ids {

		err := db.Model(w).Where("id = ?", id).Delete(w)
		deleted += err.RowsAffected
		if len(err.GetErrors()) != 0 {
			return deleted, err.GetErrors()[0]
		}
	}

	return deleted, nil
}

// UpdateWorkspace updates a workspace, using the id supplied.
func (db *DB) UpdateWorkspace(ws Workspace) error {

	err := db.Save(&ws)
	if len(err.GetErrors()) != 0 {
		return err.GetErrors()[0]
	}

	return nil
}
