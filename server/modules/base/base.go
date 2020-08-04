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

package base

import (
	"fmt"

	"github.com/sirupsen/logrus"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	"github.com/maxlandon/wiregost/server/log"
)

// Module - The base module, embedding a protobuf object
type Module struct {
	Info   *modulepb.Module    // Base module information
	Opts   *map[string]*Option // Module options
	Client *clientpb.Client    // The console client (and its user) using the module
	Log    *logrus.Entry       // Each module can log its output to the console and to log files.
}

// Option - Returns one of the module's options, by name
func (m *Module) Option(name string) (opt *Option, err error) {
	if opt, found := (*m.Opts)[name]; found {
		return opt, nil
	}

	return nil, fmt.Errorf("invalid option: %s", name)
}

// ParseMetadata - Looks for the module metadata in its respective source directory.
func (m *Module) ParseMetadata() error {
	return nil
}

// SetLogger - Initializes logging for the module
func (m *Module) SetLogger(client *clientpb.Client) {
	m.Log = log.ModuleLogger(m.Info.Path, m.Client)
}

// PreRunChecks - All checks for session, commands, options, etc. are done in this function.
// IT IS MANDATORY TO CALL THIS FUNCTION at the beginning of the any module.
func (m *Module) PreRunChecks(cmd string) (err error) {

	err = m.CheckCommand(cmd)
	if err != nil {
		return err
	}
	err = m.CheckRequiredOptions()
	if err != nil {
		return err
	}

	return
}

// CheckRequiredOptions - Checks that all required options have a value
func (m *Module) CheckRequiredOptions() (err error) {
	return
}

// CheckCommand - Verifies the command run by the user (like 'run exploit' or 'run check') is valid
func (m *Module) CheckCommand(command string) (err error) {
	// If we don't have commands it means there is no need for them, so no error
	if len(m.Info.Commands) == 0 || m.Info.Commands == nil {
		return nil
	}

	// Else check for it
	for com, desc := range m.Info.Commands {
		if com == command {
			return nil
		}
	}
	return fmt.Errorf("invalid command: %s", command)
}

// Event - Pushes an event message (ex: for status) back to the console running the module.
// It also logs the event to the module user's log file.
// NOTE: WE SHOULD USE THIS LOGGER INTERFACE AS A WAY TO PUSH EVENT MESSAGES
//       THIS METHOD SHOULD BE MOFIFIED SO AS TO USE THE MODULE LOGGER TRANSPARENTLY
func (m *Module) Event(event string, pending bool) {
}

// Asset - Find the path of an asset in the module source directory.
func (m *Module) Asset(path string) (filePath string, err error) {
	return
}

// ToProtobuf - A user requested the module information and options.
func (m *Module) ToProtobuf() (modpb *modulepb.Module) {
	modpb = m.Info
	for _, opt := range *m.Opts {
		modpb.Options[opt.proto.Name] = &opt.proto
	}
	return
}

// OptionsToProtobuf - A user requested the module options.
func (m *Module) OptionsToProtobuf() (options map[string]modulepb.Option) {
	for _, opt := range *m.Opts {
		options[opt.proto.Name] = opt.proto
	}
	return
}
