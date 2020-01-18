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
	"errors"
	"math/rand"
	"strconv"
	"time"
)

// Host represents a computer system
type Host struct {
	ID                  int `json:"host_id" sql:"host_id,pk"`
	Address             string
	MAC                 string
	Comm                string
	Name                string
	State               string
	OSName              string
	OSFlavor            string
	OSSp                string
	OSLang              string
	OSFamily            string
	Arch                string
	DetectedArch        string
	WorkspaceID         int `sql:"workspace_id,notnull,on_delete:CASCADE"`
	Workspace           *Workspace
	Purpose             string
	Info                string
	Comments            string
	Scope               string
	VirtualHost         string
	NoteCount           int
	VulnCount           int
	ServiceCount        int
	HostDetailCount     int
	ExploitAttemptCount int
	CredCount           int
	CreatedAt           string
	UpdatedAt           string
}

// NewHost instantiates a new Host. Exported so that each Host instantiated elsewhere
// in Wiregost will be immediately given an ID.
func NewHost() *Host {
	host := &Host{
		ID:        rand.Int(),
		CreatedAt: time.Now().Format("2006-01-02T15:04:05"),
	}

	return host
}

// Hosts returns all Host entries in the database
func (db *DB) Hosts() ([]*Host, error) {
	var hosts []*Host
	err := db.Model(&hosts).Select()
	if err != nil {
		return nil, err
	}

	return hosts, err
}

// FindOrCreateHost searches through the Database for a Host entry: reports one if none found.
func (db *DB) FindOrCreateHost(opts map[string]string) (*Host, error) {
	host, err := db.GetHost(opts)
	// If not host is found, create one and fill values given
	if host == nil {
		h := NewHost()
		ws, found := opts["workspace_id"]
		if found {
			h.WorkspaceID, _ = strconv.Atoi(ws)
		}
		addr, found := opts["address"]
		if found {
			h.Address = addr
		}

		return db.ReportHost(*h)
	}

	return host, err
}

// DeleteHost deletes a single host, based on the id passed as argument
func (db *DB) DeleteHost(hostID int) (int, error) {
	h := new(Host)
	res, err := db.Model(h).Where("host_id = ?", hostID).Delete()
	deleted := res.RowsAffected()
	if err != nil {
		return deleted, err
	}
	return deleted, nil
}

// DeleteHosts deletes Host entries based on the IDs passed as arguments
func (db *DB) DeleteHosts(ids []int) (rows int, err error) {
	h := new(Host)
	var deleted int
	for _, id := range ids {
		res, err := db.Model(h).Where("host_id = ?", id).Delete()
		deleted += res.RowsAffected()
		if err != nil {
			return deleted, err
		}
	}

	return deleted, nil
}

// GetHost returns a host based on options passed as argument
func (db *DB) GetHost(opts map[string]string) (*Host, error) {
	h := new(Host)

	// Find host by ID, and return it if found, return error otherwise.
	id, found := opts["host_id"]
	if found {
		id, _ := strconv.Atoi(id)
		err := db.Model(h).Where("host_id= ?", id).Select()
		if err != nil {
			return nil, err
		}
		return h, nil
	}

	// Workspace ID is required if no HostID is given, and needs to be cast
	ws, found := opts["workspace_id"]
	if !found {
		return nil, errors.New("Workspace ID is required")
	}
	wsID, _ := strconv.Atoi(ws)

	// Find host by address
	addr, found := opts["address"]
	if found {
		err := db.Model(h).Where("workspace_id = ?", wsID).
			Where("address = ?", addr).Select()
		if err != nil {
			return nil, err
		}
		return h, nil
	}

	return nil, nil
}

// UpdateHost updates a Host entry with the Host struct passed as argument.
func (db *DB) UpdateHost(h Host) (*Host, error) {
	h.UpdatedAt = time.Now().Format("2006-01-02T15:04:05")
	err := db.Update(&h)
	if err != nil {
		return nil, err
	}

	return &h, err
}

// HasHost checks if a Host entry exists in the workspace passed as argument, that
// matches the IP Address passed as argument
func (db *DB) HasHost(workspaceID int, address string) (bool, error) {
	h := new(Host)
	err := db.Model(&h).Where("workspace_id = ?", workspaceID).
		Where("address", address).Select()
	if err != nil {
		return false, err
	}

	return true, nil
}

// ReportHost adds a Host entry to the database
func (db *DB) ReportHost(h Host) (*Host, error) {
	// Add Host (no need to set ID, already exists with NewHost())
	err := db.Insert(&h)
	if err != nil {
		return nil, err
	}

	return &h, err
}

// HostStateChanged updates the state of a Host entry
// func (db *DB) HostStateChanged(hostID int, ostate string) error {
//
// }
