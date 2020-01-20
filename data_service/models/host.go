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
	"time"
)

// Host represents a computer system. It holds all the information necessary
// to all tools acting upon a Host entity.
type Host struct {
	// Identification
	ID          uint
	WorkspaceID uint `gorm:"not null"`

	// General
	MAC  string
	Comm string
	Name string

	// Operating system (filled by other tools when OS is determined)
	OSName   string
	OSFlavor string
	OSSp     string
	OSLang   string
	OSFamily string
	Arch     string

	// Scope
	Purpose     string
	Info        string
	Scope       string
	VirtualHost string

	// Nmap Attributes
	Distance      Distance      `xml:"distance"`
	EndTime       Timestamp     `xml:"endtime,attr,omitempty"`
	IPIDSequence  IPIDSequence  `xml:"ipidsequence" json:"ip_id_sequence"`
	OS            OS            `xml:"os"`
	TCPSequence   TCPSequence   `xml:"tcpsequence"`
	TCPTSSequence TCPTSSequence `xml:"tcptssequence" json:"tcp_ts_sequence"`
	Times         Times         `xml:"times"`
	Trace         Trace         `xml:"trace"`
	Uptime        Uptime        `xml:"uptime"`
	Comment       string        `xml:"comment,attr"`
	StartTime     Timestamp     `xml:"starttime,attr,omitempty"`
	Addresses     []Address     `xml:"address" gorm:"foreignkey:HostID"`
	Status        Status        `xml:"status"`
	ExtraPorts    []ExtraPort   `xml:"ports>extraports"`
	Hostnames     []Hostname    `xml:"hostnames>hostname"`
	HostScripts   []Script      `xml:"hostscript>script"`
	Ports         []Port        `xml:"ports>port"`
	Smurfs        []Smurf       `xml:"smurf"`

	// Timestamp
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewHost instantiates a Host and gives it an ID and a workspaceID
func NewHost(workspaceID uint) *Host {
	host := &Host{
		WorkspaceID: workspaceID,
	}

	return host
}

// Hosts returns all Host entries in the database
func (db *DB) Hosts(wsID uint, opts map[string]interface{}) ([]*Host, error) {
	var hosts []*Host

	// If ID, no need to search with other arguments
	id, found := opts["host_id"].(int)
	if found {
		hostID := uint(id)
		if err := db.Where(&Host{ID: hostID}).Find(&hosts).Error; err != nil {
			return nil, err
		}
		return hosts, nil
	}

	// Load all hosts in workspace
	if err := db.Where(&Host{WorkspaceID: wsID}).Find(&hosts).Error; err != nil {
		return nil, err
	}
	// Load all adddresses for each host
	for _, h := range hosts {
		err := db.Where(&Address{HostID: h.ID}).Find(&h).Error
		if err != nil {
			break
		}

	}
	return hosts, nil
}

// ReportHost adds a Host to the database, and returns it with an ID
func (db *DB) ReportHost(wsID uint) (*Host, error) {
	host := NewHost(wsID)

	if err := db.Create(&host).Select(&host).Error; err != nil {
		return nil, err
	} else {
		// db.Last(&host)
		return host, nil
	}
}

// DeleteHost deletes a Host based on the ID passed in
func (db *DB) DeleteHost(wsID uint, opts map[string]interface{}) (int64, error) {
	h := new(Host)
	var deleted int64

	id, found := opts["host_id"].(int)
	if found {
		hostID := uint(id)
		err := db.Model(h).Where("id = ?", hostID).Delete(h)
		deleted += err.RowsAffected
		if len(err.GetErrors()) != 0 {
			return deleted, err.GetErrors()[0]
		}
	} else {
		return 0, errors.New("Error: No Host ID provided")
	}

	return deleted, nil
}

// UpdateHost updates a Host, using the Host object supplied
func (db *DB) UpdateHost(h Host) (*Host, error) {
	host := &Host{}
	err := db.Save(&h).Select(&host)
	if len(err.GetErrors()) != 0 {
		return &h, err.GetErrors()[0]
	}
	// Load IP addresses
	db.Where(&Address{HostID: host.ID}).Find(&h)

	return host, nil
}
