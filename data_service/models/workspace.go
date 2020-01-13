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
)

// Workspace is the exported object
type Workspace struct {
	ID             int          `json:"id"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Boundary       string       `json:"boundary"`
	LimitToNetwork bool         `json:"limit_to_network"`
	ModuleStack    *ModuleStack `json:"module_stack"`
	CreatedAt      string       `json:"created_at"`
	UpdatedAt      string       `json:"updated_at"`
}

func NewWorkspace(name string) *Workspace {
	w := &Workspace{
		ID:   rand.Int(),
		Name: name,
	}
	return w
}

func (db *DB) Workspaces() ([]*Workspace, error) {
	var workspaces []*Workspace
	err := db.Model(&workspaces).Select()
	if err != nil {
		return nil, err
	}
	return workspaces, err
}

func (db *DB) FindWorkspace(name string) (*Workspace, error) {
	workspace := new(Workspace)
	err := db.Model(workspace).Where("name = ?", name).Select()
	if err != nil {
		return nil, err
	}
	return workspace, err
}

func (db *DB) AddWorkspaces(names []string) error {
	for _, name := range names {
		workspace := NewWorkspace(name)
		err := db.Insert(workspace)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) DeleteWorkspaces(ids []int) (rows int, err error) {
	w := new(Workspace)
	var deleted int
	for _, id := range ids {
		res, err := db.Model(w).Where("id = ?", id).Delete()
		deleted = res.RowsAffected()
		if err != nil {
			return deleted, err
		}
	}
	return deleted, nil
}

func (db *DB) UpdateWorkspace(id int) error {
	return nil
}
