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

import (
	"math/rand"
	"time"
)

// Workspace is the exported object, needed for JSON marshalling in various
// places of Wiregost.
type Workspace struct {
	ID             int          `json:"workspace_id" sql:"workspace_id,pk"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Boundary       string       `json:"boundary"`
	LimitToNetwork bool         `json:"limit_to_network"`
	ModuleStack    *ModuleStack `json:"module_stack"`
	CreatedAt      string       `json:"created_at"`
	UpdatedAt      string       `json:"updated_at"`
}

// Instantiate a new workspace, with a unique ID. Only used by AddWorkspaces().
func newWorkspace(name string) *Workspace {
	// Get good random id
	rand.Seed(time.Now().Unix())
	id := rand.Int()

	w := &Workspace{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05"),
	}
	return w
}

// Workspaces returns all workspaces in database
func (db *DB) Workspaces() ([]*Workspace, error) {
	var workspaces []*Workspace
	err := db.Model(&workspaces).Select()
	if err != nil {
		return nil, err
	}
	return workspaces, err
}

// FindWorkspace returns a workspace queried by its name
func (db *DB) FindWorkspace(name string) (*Workspace, error) {
	workspace := new(Workspace)
	err := db.Model(workspace).Where("name = ?", name).Select()
	if err != nil {
		return nil, err
	}
	return workspace, err
}

// AddWorkspaces adds workspaces to database, using names supplied.
func (db *DB) AddWorkspaces(names []string) error {
	for _, name := range names {
		workspace := newWorkspace(name)
		err := db.Insert(workspace)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteWorkspaces adds workspaces to database, using ids supplied.
func (db *DB) DeleteWorkspaces(ids []int) (rows int, err error) {
	w := new(Workspace)
	var deleted int
	for _, id := range ids {
		res, err := db.Model(w).Where("workspace_id = ?", id).Delete()
		deleted += res.RowsAffected()
		if err != nil {
			return deleted, err
		}
	}
	return deleted, nil
}

// UpdateWorkspace updates a workspace, using the id supplied.
func (db *DB) UpdateWorkspace(ws Workspace) error {
	ws.UpdatedAt = time.Now().Format("2006-01-02T15:04:05")
	err := db.Update(&ws)
	if err != nil {
		return err
	}
	return nil
}
