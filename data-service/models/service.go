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
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Service contains detailed information about a service on an open port.
type Service struct {
	// General
	ID        uint
	PortID    uint   `gorm:"not null"`
	Proto     string `xml:"proto,attr"`
	Name      string `xml:"name,attr"`
	ExtraInfo string `xml:"extrainfo,attr"`

	//Nmap attributes
	DeviceType    string `xml:"devicetype,attr"`
	HighVersion   string `xml:"highver,attr"`
	Hostname      string `xml:"hostname,attr"`
	LowVersion    string `xml:"lowver,attr"`
	Method        string `xml:"method,attr"`
	OSType        string `xml:"ostype,attr"`
	Product       string `xml:"product,attr"`
	RPCNum        string `xml:"rpcnum,attr"`
	ServiceFP     string `xml:"servicefp,attr"`
	Tunnel        string `xml:"tunnel,attr"`
	Version       string `xml:"version,attr"`
	Configuration int    `xml:"conf,attr"`
	// CPEs          []CPE  `xml:"cpe"` // SOLVE THIS: CANNOT INSERT IN SCHEMA. CHECK NEED TO PERSIST

	// Timestamp
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewService instantiates a Service and gives it a PortID
func NewService(hostID uint) *Port {

	service := &Port{
		HostID: hostID,
	}
	return service
}

// Services returns all or part of Services in a workspace, with sequential chaining of options
func (db *DB) Services(wsID uint, opts map[string]interface{}) (services []*Port, err error) {

	// Prepare simple search filters
	tx := db.Where("workspace_id = ?", wsID)
	if opts == nil {
		return db.servicesByWorkspace(tx)
	} else {
		services, err = db.servicesByWorkspace(tx)
	}

	// Return if IDs found
	ids, idFound := opts["host_id"]
	if idFound {
		return db.servicesByHosts(&wsID, ids)
	}

	// Populate service list if addresses found
	addrs, addrFound := opts["addresses"]
	if addrFound {
		hosts, err := db.hostsByAddress(wsID, addrs)
		if err != nil {
			return nil, err
		}
		if len(hosts) == 0 {
			return nil, errors.New("No hosts match given IDs")
		}
		var ids []uint
		for i, _ := range hosts {
			ids = append(ids, hosts[i].ID)
		}
		services, err = db.servicesByHosts(&wsID, ids)
	}

	// Refine service list if ports found
	ports, portFound := opts["port"]
	if portFound {
		portServs, err := db.servicesByPort(&wsID, ports)
		if err != nil {
			return nil, err
		}
		var filteredServs []*Port
		for _, ps := range portServs {
			for _, s := range services {
				if s.HostID == ps.HostID {
					filteredServs = append(filteredServs, ps)
					break
				}
			}
		}
		services = filteredServs
	}

	// Refine with service name (http/https/dns, etc...)
	names, namesFound := opts["name"]
	if namesFound {
		servNames := strings.Split(names.(string), ",")
		var filteredServs []*Port
		for _, s := range services {
			for _, n := range servNames {
				if strings.Contains(s.Service.Name, n) {
					filteredServs = append(filteredServs, s)
					break
				}
			}
		}
		services = filteredServs
	}

	// Refine with service info
	info, infoFound := opts["info"]
	if infoFound {
		infoStr := info.(string)
		var filteredServs []*Port
		for _, s := range services {
			if strings.Contains(s.Service.ExtraInfo, infoStr) {
				filteredServs = append(filteredServs, s)
			}
		}
		services = filteredServs
	}

	return services, nil
}

// ReportService either adds a Service to the database and returns it, or returns an already existing
// Service matching provided options
func (db *DB) ReportService(hostID *uint, opts map[string]interface{}) (service *Port, err error) {

	tx := db.Where("host_id = ?", hostID)
	tx = parseServiceOptions(opts, tx)

	hostPorts, err := db.servicesByHosts(nil, []uint{*hostID})

	ports, portFound := opts["port"]
	var port uint16
	if portFound {
		// A unique port is given when reporting host
		pStr := strings.Split(ports.(string), ",")[0]
		pInt, _ := strconv.Atoi(pStr)
		port = uint16(pInt)

		for _, p := range hostPorts {
			if p.Number == port {
				return p, nil
			}
		}
		// If port is new
		goto newPort
	} else {
		return nil, errors.New("No port number given")
	}

newPort:
	// Create service and populate fields
	service = NewService(*hostID)
	if tx = db.Create(&service).Select(&service); tx.Error != nil {
		return nil, tx.Error
	}
	service.Number = port
	service.Service = Service{
		PortID: service.ID,
	}

	name, found := opts["name"]
	if found {
		service.Service.Name = name.(string)
	}
	proto, found := opts["proto"]
	if found {
		service.Protocol = proto.(string)
		service.Service.Proto = proto.(string)
	}
	info, found := opts["info"]
	if found {
		service.Service.ExtraInfo = info.(string)
	}
	db.Save(&service)

	return service, nil
}

// DeleteServices deletes one or more Services based on the IDs passed as argument
func (db *DB) DeleteServices(hostID uint, opts map[string]interface{}) (deleted int64, err error) {
	s := new(Service)

	ids, found := opts["port_id"]
	if found {
		switch idList := ids.(type) {
		case []float64:
			for i, _ := range idList {
				portID := uint16(idList[i])
				if tx := db.Model(s).Where("id = ?", portID).Delete(&s); tx.Error != nil {
					continue
				}
				deleted += 1
			}
		}
	} else {
		return 0, errors.New("Error: No PortID provided")
	}

	return deleted, nil
}

// UpdateService updates a Service entity passed as argument
func (db *DB) UpdateService(s Port) (updated *Port, err error) {
	if tx := db.Save(&s).Select(&s); tx.Error != nil {
		return &s, tx.Error
	}

	db.Model(&s).Association("Services").Replace(s.Service)
	db.Find(&s)

	return &s, nil
}

// servicesByWorkspace returns all services attached to hosts in a workspace
func (db *DB) servicesByWorkspace(tx *gorm.DB) (services []*Port, err error) {

	var hosts []*Host
	if tx = tx.Find(&hosts); tx.Error != nil {
		return nil, tx.Error
	}

	for _, h := range hosts {
		// Load ports
		var ports []*Port
		err := db.Model(&h).Related(&ports).Error
		if err != nil {
			continue
		}
		// Load service for each port
		for _, p := range ports {
			err := db.Model(&p).Related(&p.Service).Error
			if err != nil {
				continue
			}
		}
		// Append ports/service to services
		services = append(services, ports...)

	}
	return services, nil
}

// servicesByHosts returns all services attached to a list of hosts
func (db *DB) servicesByHosts(wsID *uint, hostIDs interface{}) (services []*Port, err error) {
	// Find hosts
	hosts, err := db.Hosts(wsID, map[string]interface{}{"host_id": hostIDs})
	if err != nil {
		return nil, err
	}
	if len(hosts) == 0 {
		return nil, errors.New("No host for these IDs")
	}

	// Get service and port for each host
	for i, _ := range hosts {
		var ports []*Port
		if tx := db.Where("host_id = ?", hosts[i].ID).Related(&ports); tx.Error != nil {
			return nil, tx.Error
		}

		for i, _ := range ports {
			if tx := db.Model(ports[i]).Related(&ports[i].Service); tx.Error != nil {
				return nil, tx.Error
			}
		}
		// Append host ports/services to list of port/services
		services = append(services, ports...)
	}
	return services, nil

}

// servicesByPort returns all services matching provided ports for a workspace
func (db *DB) servicesByPort(wsID *uint, ports interface{}) (services []*Port, err error) {

	portList := strings.Split(ports.(string), ",")

	hosts, err := db.Hosts(wsID, nil)
	if err != nil {
		return nil, err
	}
	if len(hosts) == 0 {
		return nil, errors.New("No hosts in workspace")
	}

	var allPorts []Port
	for i, _ := range hosts {
		if tx := db.Model(hosts[i]).Related(hosts[i].Ports); tx.Error != nil {
			return nil, tx.Error
		}
		allPorts = append(allPorts, hosts[i].Ports...)
	}

	for _, p := range portList {
		pInt, _ := strconv.Atoi(p)
		pUint := uint16(pInt)
		for _, port := range allPorts {
			if port.Number == pUint {
				db.Model(port).Related(port.Service)
				services = append(services, &port)
			}
		}
	}

	return services, nil
}

func parseServiceOptions(opts map[string]interface{}, tx *gorm.DB) *gorm.DB {

	name, found := opts["name"].(string)
	if found {
		tx = tx.Where("name ILIKE (?)", name)
	}
	info, found := opts["info"].(string)
	if found {
		tx = tx.Where("extra_info ILIKE (?)", info)
	}
	proto, found := opts["proto"].(string)
	if found {
		tx = tx.Where("proto ILIKE (?)", proto)
	}

	return tx
}
