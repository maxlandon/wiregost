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
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
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

	// Network
	Addresses []Address `xml:"address" gorm:"foreignkey:HostID"`

	// Nmap Temporary Attributes
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

// Hosts returns all Host entries in the database, with sequential chaining of options
func (db *DB) Hosts(wsID uint, opts map[string]interface{}) ([]*Host, error) {
	var hosts []*Host

	// If ID is given, return corresponding host
	id, found := opts["host_id"].(int)
	if found {
		return db.hostByID(id)
	}

	// Queries are always in a workspace context:
	tx := db.Where("workspace_id = ?", wsID)

	// Parse search options and add them to methods chain.
	// if none are provided, return all hosts in workspace
	if opts == nil {
		return db.hostsByWorkspace(tx)
	} else {
		tx = parseOptions(opts, tx)
	}

	// This section, once all options filters are coded, will not be useful, as
	// each successive search below it will populate/refine the list of hosts.
	// All of these successive searches should also, a priori, take care of populating
	// Host struct fields, such as addresses, Hostnames, etc...
	// -------------------------------------------------------------------------
	// Perform query based on parsed options chained method.
	// Returns a first batch of hosts, which can be refined after with other
	// options such as IP addresses
	if tx.Find(&hosts); tx.Error != nil {
		return nil, tx.Error
	}
	// Load addresses into each host found
	for _, h := range hosts {
		if tx := db.Model(&h).Related(&h.Addresses); tx.Error != nil {
			continue
		}
	}
	// -------------------------------------------------------------------------

	// Addresses: need to perform query here (because query logic
	// involves going back and forth).
	addrs, found := opts["addresses"]
	if found {
		hosts, _ = db.hostsByAddress(wsID, addrs, tx)
		if hosts == nil {
			// If no hosts have been returned, search is not sucessful.
			return nil, nil
		}
	}

	return hosts, nil
}

// GetHost returns a host based on options passed as argument
func (db *DB) GetHost(wsID uint, opts map[string]interface{}) (*Host, error) {
	var host *Host

	// If ID, no need to search with other arguments, immediatly return result
	id, found := opts["host_id"].(float64)
	if found {
		hostID := int(id)
		hosts, _ := db.hostByID(hostID)
		host = hosts[0]
		return host, nil
	}

	return host, nil
}

// ReportHost adds a Host to the database, and returns it with an ID
func (db *DB) ReportHost(wsID uint) (*Host, error) {
	host := NewHost(wsID)

	if err := db.Create(&host).Select(&host).Error; err != nil {
		return nil, err
	} else {
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
	if err := db.Save(&h).Select(&host); err.Error != nil {
		return &h, err.Error
	}

	// Update and load IP addresses
	db.Model(&h).Association("Addresses").Replace(h.Addresses)
	db.Where(&Address{HostID: host.ID}).Find(&h)

	return &h, nil
}

// GetHost returns a host based on options passed as argument
func (db *DB) hostByID(ID int) ([]*Host, error) {
	var hosts []*Host

	hostID := uint(ID)
	if tx := db.Where(&Host{ID: hostID}).Find(&hosts); tx.Error != nil {
		return nil, tx.Error
	}
	if err := db.Model(&hosts[0]).Related(&hosts[0].Addresses); err.Error != nil {
		return hosts, err.Error
	}
	return hosts, nil
}

// workspaceHosts queries all hosts in a workspace
func (db *DB) hostsByWorkspace(tx *gorm.DB) ([]*Host, error) {
	var hosts []*Host
	if tx = tx.Find(&hosts); tx.Error != nil {
		return nil, tx.Error
	}
	// Load all adddresses for each host
	for _, h := range hosts {
		err := db.Model(&h).Related(&h.Addresses).Error
		if err != nil {
			continue
		}
	}
	return hosts, nil
}

// hostsByAddress is given a workspaceID, a list of addresses to process and a tx context (with possibly
// other required search filters). It then refines a list based on these addresses, and returns results.
func (db *DB) hostsByAddress(workspaceID uint, addrs interface{}, tx *gorm.DB) ([]*Host, error) {

	var hosts []*Host

	// Convert addrs to []string{}
	s := reflect.ValueOf(addrs)
	a := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		a[i] = s.Index(i).Interface()
	}
	addrStr := make([]string, 0)
	for _, item := range a {
		addrStr = append(addrStr, item.(string))
	}

	// Load addresses
	var addresses []Address
	for _, addr := range addrStr {
		var tempAddr []Address
		db.Where("addr = ?", addr).Find(&tempAddr)
		for _, a := range tempAddr {
			addresses = append(addresses, a)
		}
	}

	var unfiltered []Host
	// load hosts for each address, and addresses for each host
	if len(addresses) != 0 {
		for _, addr := range addresses {
			h := Host{}
			if tx.Model(&addr).Related(&h).RecordNotFound() {
				continue
			} else {
				if tx := db.Model(&h).Related(&h.Addresses); tx.Error != nil {
					continue
				}
				unfiltered = append(unfiltered, h)
			}
		}
		// Filter hosts for redundant elements
		found := map[uint]bool{}
		for _, h := range unfiltered {
			if found[h.ID] {
				break
			} else {
				if h.WorkspaceID == workspaceID {
					found[h.ID] = true
					hosts = append(hosts, &h)
				}
			}
		}
		return hosts, nil
	} else {
		return nil, nil
	}
}

// FindOrCreateHost searches through the Database for a Host entry: reports one if none found.
// func (db *DB) FindOrCreateHost(opts map[string]string) (*Host, error) {
// }

// HasHost checks if a Host entry exists in the workspace passed as argument, that
// matches the IP Address passed as argument
func (db *DB) HasHost(workspaceID uint, address string) bool {

	addrs := []string{address}
	tx := db.Where("workspace_id = ?", workspaceID)

	hosts, err := db.hostsByAddress(workspaceID, addrs, tx)
	if err != nil {
		return false
	}
	if hosts == nil {
		return false
	} else {
		return true
	}
}

// parseOptions extracts search options and construct and chain of query conditions
// that is passed to functions needing it.
func parseOptions(opts map[string]interface{}, tx *gorm.DB) *gorm.DB {

	// OS Name
	osName, found := opts["os_name"].(string)
	if found {
		tx = tx.Where("os_name ILIKE (?)", osName)
	}
	// OS Flavor
	osFlav, found := opts["os_flavor"].(string)
	if found {
		tx = tx.Where("os_flavor ILIKE (?)", osFlav)
	}
	// OS Family
	osFam, found := opts["os_family"].(string)
	if found {
		tx = tx.Where("os_family ILIKE (?)", osFam)
	}
	// Architecture
	arch, found := opts["arch"].(string)
	if found {
		tx = tx.Where("arch ILIKE (?)", arch)
	}

	return tx
}
