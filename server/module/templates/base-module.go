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

package templates

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	consts "github.com/maxlandon/wiregost/client/constants"
	pb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/core"
)

// Module - Contains all properties of a module
type Module struct {
	Name        string   `json:"name"`        // Name of the module
	Type        string   `json:"type"`        // Type of module (auxiliary, exploit, post, payload)
	Path        []string `json:"path"`        // Path to the module (ie. post/windows/x64/powershell/gather/powerview)
	Description string   `json:"description"` // Description of the module
	Notes       string   `json:"notes"`       // Notes about the module
	References  []string `json:"references"`  // A list of references to vulnerabilities/others (ie. CVEs)
	Author      []string `json:"author"`      // A list of module authors
	Credits     []string `json:"credits"`     // A list of people to credit for underlying tools/techniques
	Platform    string   `json:"platform"`    // Operating system the module can run on.
	Targets     []string `json:"targets"`     // A list of operating system versions the modules works on
	Arch        string   `json:"arch"`        // CPU architecture for which the module works
	Lang        string   `json:"lang"`        // Programming language in which the module is written
	Priviledged bool     `json:"priviledged"` // Does the module requires administrator privileges

	Options map[string]*Option
}

// Option - Module option
type Option struct {
	Name        string `json:"name"`        // Name of the option
	Value       string `json:"value"`       // Value of the option (default is filled here)
	Required    bool   `json:"required"`    // Is this a required option ?
	Flag        string `json:"flag"`        // Flag value of the option, used for execution
	Description string `json:"description"` // A description of the option
}

// ToProtobuf - Returns the protobuf version of a module
func (m *Module) ToProtobuf() *pb.Module {
	mod := &pb.Module{
		Name:        m.Name,
		Type:        m.Type,
		Path:        m.Path,
		Description: m.Description,
		Notes:       m.Notes,
		References:  m.References,
		Author:      m.Author,
		Credits:     m.Credits,
		Platform:    m.Platform,
		Targets:     m.Targets,
		Arch:        m.Arch,
		Lang:        m.Lang,
		Priviledged: m.Priviledged,
		Options:     map[string]*pb.Option{},
	}

	for name, opt := range m.Options {
		mod.Options[name] = opt.ToProtobuf()
	}

	return mod
}

// ToProtobuf - Returns the protobuf version of a module option
func (o *Option) ToProtobuf() *pb.Option {
	return &pb.Option{
		Name:        o.Name,
		Value:       o.Value,
		Required:    o.Required,
		Flag:        o.Flag,
		Description: o.Description,
	}
}

// Init - Module initialization, loads metadata.
func (m *Module) Init() error {
	// func (m *Module) Init() error {

	m.Options = make(map[string]*Option)

	file, err := os.Open(filepath.Join(assets.GetModulesDir(),
		strings.Join(m.Path, "/"), "metadata.json"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	metadata, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = json.Unmarshal(metadata, &m)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// ToProtobuf - Returns the protobuf version of a module
func (m *Module) ParseProto(pbmod *pb.Module) {
	m.Name = pbmod.Name
	m.Type = pbmod.Type
	m.Path = pbmod.Path
	m.Description = pbmod.Description
	m.Notes = pbmod.Notes
	m.References = pbmod.References
	m.Author = pbmod.Author
	m.Credits = pbmod.Credits
	m.Platform = pbmod.Platform
	m.Targets = pbmod.Targets
	m.Arch = pbmod.Arch
	m.Lang = pbmod.Lang
	m.Priviledged = pbmod.Priviledged
	m.Options = map[string]*Option{}

	for name, opt := range pbmod.Options {
		option := Option{}
		option.Name = opt.Name
		option.Value = opt.Value
		option.Flag = opt.Flag
		option.Required = opt.Required
		option.Description = opt.Description
		m.Options[name] = &option
	}
}

func (m *Module) SetOption(option, value string) {
	opt := m.Options[option]
	opt.Value = value
}

func (m *Module) CheckRequiredOptions() (ok bool, err error) {
	// Check every 'required' option to make sure it isn't null
	for _, v := range m.Options {
		if v.Required {
			if v.Value == "" {
				return false, errors.New(v.Name + " is required")
			}
		}
	}

	return true, nil
}

func (m *Module) GetSession() (session *core.Ghost, err error) {

	// Check empty session
	if m.Options["Session"].Value == "" {
		return nil, errors.New("Provide a Session to run this module on.")
	}

	// Check connected session
	if 0 < len(*core.Wire.Ghosts) {
		for _, g := range *core.Wire.Ghosts {
			if g.Name == m.Options["Session"].Value {
				session = g
			}
		}
	}

	if session == nil {
		invalid := fmt.Sprintf("Invalid or non-connected session: %s", m.Options["Session"].Value)
		return nil, errors.New(invalid)
	}

	// Check valid platform
	platform := ""
	switch m.Platform {
	case "windows", "win", "Windows":
		platform = "windows"
	case "darwin", "ios", "macos", "MacOS", "Apple":
		platform = "darwin"
	case "Linux", "linux":
		platform = "linux"
	}

	if platform != session.OS {
		return nil, errors.New("The session's target OS is not supported by this module")
	}

	return session, nil
}

func (m *Module) ModuleEvent(requestID int32, event string) {
	core.EventBroker.Publish(core.Event{
		EventType:       consts.ModuleEvent,
		EventSubType:    "run",
		ModuleRequestID: requestID,
		Data:            []byte(event),
	})
}
