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
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"strconv"
	"time"
)

// Run represents an nmap scanning run.
type Run struct {
	XMLName xml.Name `xml:"nmaprun"`

	Args             string         `xml:"args,attr"`
	ProfileName      string         `xml:"profile_name,attr"`
	Scanner          string         `xml:"scanner,attr"`
	StartStr         string         `xml:"startstr,attr"`
	Version          string         `xml:"version,attr"`
	XMLOutputVersion string         `xml:"xmloutputversion,attr"`
	Debugging        Debugging      `xml:"debugging"`
	Stats            Stats          `xml:"runstats"`
	ScanInfo         ScanInfo       `xml:"scaninfo"`
	Start            Timestamp      `xml:"start,attr"`
	Verbose          Verbose        `xml:"verbose"`
	Hosts            []Host         `xml:"host"`
	PostScripts      []Script       `xml:"postscript>script"`
	PreScripts       []Script       `xml:"prescript>script"`
	Targets          []Target       `xml:"target"`
	TaskBegin        []Task         `xml:"taskbegin"`
	TaskProgress     []TaskProgress `xml:"taskprogress"`
	TaskEnd          []Task         `xml:"taskend"`

	NmapErrors []string
	rawXML     []byte
}

// ToFile writes a Run as XML into the specified file path.
func (r Run) ToFile(filePath string) error {
	return ioutil.WriteFile(filePath, r.rawXML, 0666)
}

// ToReader writes the raw XML into an streamable buffer.
func (r Run) ToReader() io.Reader {
	return bytes.NewReader(r.rawXML)
}

// ScanInfo represents the scan information.
type ScanInfo struct {
	NumServices int    `xml:"numservices,attr"`
	Protocol    string `xml:"protocol,attr"`
	ScanFlags   string `xml:"scanflags,attr"`
	Services    string `xml:"services,attr"`
	Type        string `xml:"type,attr"`
}

// Verbose contains the verbosity level of the scan.
type Verbose struct {
	Level int `xml:"level,attr"`
}

// Debugging contains the debugging level of the scan.
type Debugging struct {
	Level int `xml:"level,attr"`
}

// Task contains information about a task.
type Task struct {
	Time      Timestamp `xml:"time,attr"`
	Task      string    `xml:"task,attr"`
	ExtraInfo string    `xml:"extrainfo,attr"`
}

// TaskProgress contains information about the progression of a task.
type TaskProgress struct {
	Percent   float32   `xml:"percent,attr"`
	Remaining int       `xml:"remaining,attr"`
	Task      string    `xml:"task,attr"`
	Etc       Timestamp `xml:"etc,attr"`
	Time      Timestamp `xml:"time,attr"`
}

// Target represents a target, how it was specified when passed to nmap,
// its status and the reason for its status. Example:
// <target specification="domain.does.not.exist" status="skipped" reason="invalid"/>
type Target struct {
	Specification string `xml:"specification,attr"`
	Status        string `xml:"status,attr"`
	Reason        string `xml:"reason,attr"`
}

// Times contains time statistics for an nmap scan.
type Times struct {
	SRTT string `xml:"srtt,attr"`
	RTT  string `xml:"rttvar,attr"`
	To   string `xml:"to,attr"`
}

// Stats contains statistics for an nmap scan.
type Stats struct {
	Finished Finished  `xml:"finished"`
	Hosts    HostStats `xml:"hosts"`
}

// Finished contains detailed statistics regarding a finished scan.
type Finished struct {
	Time     Timestamp `xml:"time,attr"`
	TimeStr  string    `xml:"timestr,attr"`
	Elapsed  float32   `xml:"elapsed,attr"`
	Summary  string    `xml:"summary,attr"`
	Exit     string    `xml:"exit,attr"`
	ErrorMsg string    `xml:"errormsg,attr"`
}

// HostStats contains the amount of up and down hosts and the total count.
type HostStats struct {
	Up    int `xml:"up,attr"`
	Down  int `xml:"down,attr"`
	Total int `xml:"total,attr"`
}

// Timestamp represents time as a UNIX timestamp in seconds.
type Timestamp time.Time

// ParseTime converts a UNIX timestamp string to a time.Time.
func (t *Timestamp) ParseTime(s string) error {
	timestamp, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(timestamp, 0))

	return nil
}

// FormatTime formats the time.Time value as a UNIX timestamp string.
func (t Timestamp) FormatTime() string {
	return strconv.FormatInt(time.Time(t).Unix(), 10)
}

// MarshalJSON implements the json.Marshaler interface.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(t.FormatTime()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	return t.ParseTime(string(b))
}

// MarshalXMLAttr implements the xml.MarshalerAttr interface.
func (t Timestamp) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if time.Time(t).IsZero() {
		return xml.Attr{}, nil
	}

	return xml.Attr{Name: name, Value: t.FormatTime()}, nil
}

// UnmarshalXMLAttr implements the xml.UnmarshalXMLAttr interface.
func (t *Timestamp) UnmarshalXMLAttr(attr xml.Attr) (err error) {
	return t.ParseTime(attr.Value)
}

// Script represents an Nmap Scripting Engine script.
// The inner elements can be an arbitrary collection of Tables and Elements. Both of them can also be empty.
type Script struct {
	// General
	PortID uint16 `gorm:"not null"`

	// Nmap attributes
	ID       string    `xml:"id,attr" json:"id"`
	Output   string    `xml:"output,attr" json:"output"`
	Elements []Element `xml:"elem,omitempty" json:"elements,omitempty"`
	Tables   []Table   `xml:"table,omitempty" json:"tables,omitempty"`

	// Timestamp
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Table is an arbitrary collection of (sub-)Tables and Elements. All its fields can be empty.
type Table struct {
	Key      string    `xml:"key,attr,omitempty" json:"key,omitempty"`
	Tables   []Table   `xml:"table,omitempty" json:"tables,omitempty"`
	Elements []Element `xml:"elem,omitempty" json:"elements,omitempty"`
}

// Element is the smallest building block for scripts/tables. It can optionally(!) have a key.
type Element struct {
	Key   string `xml:"key,attr,omitempty" json:"key,omitempty"`
	Value string `xml:",innerxml" json:"value"`
}

// Status represents a host's status
type Status struct {
	ID        uint
	HostID    uint    `gorm:"not null"`
	State     string  `xml:"state,attr"`
	Reason    string  `xml:"reason,attr"`
	ReasonTTL float32 `xml:"reason_ttl,attr"`
}

func (s Status) String() string {
	return s.State
}

// Hostname is a name for a host.
type Hostname struct {
	ID     uint
	HostID uint   `gorm:"not null"`
	Name   string `xml:"name,attr"`
	Type   string `xml:"type,attr"`
}

func (h Hostname) String() string {
	return h.Name
}

// Smurf contains responses from a smurf attack.
type Smurf struct {
	Responses string `xml:"responses,attr"`
}

// Distance is the amount of hops to a particular host.
type Distance struct {
	Value int `xml:"value,attr"`
}

// Uptime is the amount of time the host has been up.
type Uptime struct {
	Seconds  int    `xml:"seconds,attr"`
	Lastboot string `xml:"lastboot,attr"`
}

// Sequence represents a detected sequence.
type Sequence struct {
	Class  string `xml:"class,attr"`
	Values string `xml:"values,attr"`
}

// TCPSequence represents a detected TCP sequence.
type TCPSequence struct {
	Index      int    `xml:"index,attr"`
	Difficulty string `xml:"difficulty,attr"`
	Values     string `xml:"values,attr"`
}

// IPIDSequence represents a detected IP ID sequence.
type IPIDSequence Sequence

// TCPTSSequence represents a detected TCP TS sequence.
type TCPTSSequence Sequence

// Trace represents the trace to a host, including the hops.
type Trace struct {
	Proto string `xml:"proto,attr"`
	Port  int    `xml:"port,attr"`
	Hops  []Hop  `xml:"hop"`
}

// Hop is an IP hop to a host.
type Hop struct {
	TTL    float32 `xml:"ttl,attr"`
	RTT    string  `xml:"rtt,attr"`
	IPAddr string  `xml:"ipaddr,attr"`
	Host   string  `xml:"host,attr"`
}
