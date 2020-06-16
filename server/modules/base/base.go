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
	"github.com/sirupsen/logrus"

	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/log"
)

// Module - The base module, embedding a protobuf object
type Module struct {
	Proto  *modulepb.Module // Base module information
	User   *serverpb.User   // The user who loaded the module
	Logger *logrus.Entry    // Each module logs its ouput to the user's log file
}

// Option - Returns one of the module's options, by name
func (m *Module) Option(name string) *modulepb.Option {
	return m.Proto.Options[name]
}

// ParseMetadata - Looks for the module metadata in its respective source directory.
func (m *Module) ParseMetadata() error {
	return nil
}

// SetLogger - Initializes logging for the module
func (m *Module) SetLogger() {
	m.Logger = log.UserLogger(m.User.Name, m.User.ID, m.Proto.Path, "module")
}

// CheckRequiredOptions - Checks that all required options have a value
func (m *Module) CheckRequiredOptions() (ok bool, err error) {
	return
}

// Event - Pushes an event message (ex: for status) back to the console running the module.
// It also logs the event to the module user's log file.
func (m *Module) Event(event string, pending bool) {
}

// Asset - Find the path of an asset in the module source directory.
func (m *Module) Asset(path string) (filePath string, err error) {
	return
}

// ToProtobuf - When clients request a copy of the module, send it.
func (m *Module) ToProtobuf() (modpb *modulepb.Module) {
	return m.Proto
}
