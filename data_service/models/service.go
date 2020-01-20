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

// Service contains detailed information about a service on an open port.
type Service struct {

	//Nmap attributes
	DeviceType    string `xml:"devicetype,attr"`
	ExtraInfo     string `xml:"extrainfo,attr"`
	HighVersion   string `xml:"highver,attr"`
	Hostname      string `xml:"hostname,attr"`
	LowVersion    string `xml:"lowver,attr"`
	Method        string `xml:"method,attr"`
	Name          string `xml:"name,attr"`
	OSType        string `xml:"ostype,attr"`
	Product       string `xml:"product,attr"`
	Proto         string `xml:"proto,attr"`
	RPCNum        string `xml:"rpcnum,attr"`
	ServiceFP     string `xml:"servicefp,attr"`
	Tunnel        string `xml:"tunnel,attr"`
	Version       string `xml:"version,attr"`
	Configuration int    `xml:"conf,attr"`
	CPEs          []CPE  `xml:"cpe"`
}
