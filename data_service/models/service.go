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

type Service struct {
	ID        int `sql:"service_id,pk"`
	HostID    int `sql:"host_id,notnull,on_delete:CASCADE"`
	Host      *Host
	Port      int
	Proto     string
	Name      string
	Info      string
	State     string
	CreatedAt string
	UpdatedAt string
}

// NewService instantiates a new Service. Exported so that each Service instantiated elsewhere
// in Wiregost will be immediately given an ID.
func NewService() *Service {
	// Get good random id
	rand.Seed(time.Now().Unix())
	id := rand.Int()

	service := &Service{
		ID: id,
		// Maybe we will need to input a HostID right now,
		// so that we don't forget later and create bugs.
		CreatedAt: time.Now().Format("2006-01-02T15:04:05"),
	}

	return service
}

// Services returns all Service entries in the database
func (db *DB) Services() ([]*Service, error) {
	var services []*Service
	err := db.Model(&services).Select()
	if err != nil {
		return nil, err
	}

	return services, err
}

// FindOrCreateService searches through the Database for a Service entry: reports one if none found.
func (db *DB) FindOrCreateService(opts map[string]string) (*Service, error) {
	service, err := db.GetService(opts)
	// If not host is found, create one and fill values given
	if service == nil {
		s := NewService()
		ws, found := opts["host_id"]
		if found {
			s.HostID, _ = strconv.Atoi(ws)
		}
		proto, found := opts["proto"]
		if found {
			s.Proto = proto
		}
		port, found := opts["port"]
		if found {
			s.Port, _ = strconv.Atoi(port)
		}

		return db.ReportService(*s)
	}

	return service, err
}

// DeleteServices deletes Service entries based on the IDs passed as arguments
func (db *DB) DeleteServices(ids []int) (int, error) {
	s := new(Service)
	var deleted int
	for _, id := range ids {
		res, err := db.Model(s).Where("host_id = ?", id).Delete()
		deleted += res.RowsAffected()
		if err != nil {
			return deleted, err
		}
	}

	return deleted, nil
}

// GetService returns a Service based on options passed as argument
func (db *DB) GetService(opts map[string]string) (*Service, error) {
	s := new(Service)

	// Find host by ID, and return it if found, return error otherwise.
	id, found := opts["host_id"]
	if found {
		id, _ := strconv.Atoi(id)
		err := db.Model(s).Where("host_id= ?", id).Select()
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	// Workspace ID is required if no HostID is given, and needs to be cast
	ws, found := opts["workspace_id"]
	if !found {
		return nil, errors.New("Workspace ID is required")
	}
	wsID, _ := strconv.Atoi(ws)

	// Find host by address
	proto, found := opts["proto"]
	port, found := opts["port"]
	if found {
		err := db.Model(s).Where("workspace_id = ?", wsID).
			Where("proto = ?", proto).
			Where("port = ?", port).Select()
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	return nil, nil
}

// UpdateService updates a Service entry with the Service struct passed as argument.
func (db *DB) UpdateService(s Service) (*Service, error) {
	s.UpdatedAt = time.Now().Format("2006-01-02T15:04:05")
	err := db.Update(&s)
	if err != nil {
		return nil, err
	}

	return &s, err
}

// ReportService adds a Service entry to the database
func (db *DB) ReportService(s Service) (*Service, error) {
	// Add Service (no need to set ID, already exists with NewService())
	err := db.Insert(&s)
	if err != nil {
		return nil, err
	}

	return &s, err
}
