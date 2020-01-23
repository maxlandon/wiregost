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

// Port contains all the information about a scanned port.
type Port struct {
	ID       uint16   `xml:"portid,attr"`
	Protocol string   `xml:"protocol,attr"`
	Owner    Owner    `xml:"owner"`
	Service  Service  `xml:"service"`
	State    State    `xml:"state"`
	Scripts  []Script `xml:"script"`
}

// ExtraPort contains the information about the closed and filtered ports.
type ExtraPort struct {
	State   string   `xml:"state,attr"`
	Count   int      `xml:"count,attr"`
	Reasons []Reason `xml:"extrareasons"`
}

// Reason represents a reason why a port is closed or filtered.
// This won't be in the scan results unless WithReason is used.
type Reason struct {
	Reason string `xml:"reason,attr"`
	Count  int    `xml:"count,attr"`
}

// PortStatus represents a port's state.
type PortStatus string

// Enumerates the different possible state values.
const (
	Open       PortStatus = "open"
	Closed     PortStatus = "closed"
	Filtered   PortStatus = "filtered"
	Unfiltered PortStatus = "unfiltered"
)

// Status returns the status of a port.
func (p Port) Status() PortStatus {
	return PortStatus(p.State.State)
}

// State contains information about a given port's status.
// State will be open, closed, etc.
type State struct {
	State     string  `xml:"state,attr"`
	Reason    string  `xml:"reason,attr"`
	ReasonIP  string  `xml:"reason_ip,attr"`
	ReasonTTL float32 `xml:"reason_ttl,attr"`
}

func (s State) String() string {
	return s.State
}

// Owner contains the name of a port's owner.
type Owner struct {
	Name string `xml:"name,attr"`
}

func (o Owner) String() string {
	return o.Name
}
