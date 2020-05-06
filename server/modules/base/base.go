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
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

// Module - The base module, embedding a protobuf object
type Module struct {
	Proto *modulepb.Module
}

// Option - Returns one of the module's options, by name
func (m *Module) Option(name string) *modulepb.Option {
	return m.Proto.Options[name]
}

// ParseMetadata - Looks for the module metadata in its respective source directory.
func (m *Module) ParseMetadata() error {
	return nil
}

// CheckRequiredOptions - Checks that all required options have a value
func (m *Module) CheckRequiredOptions() (ok bool, err error) {
	return
}

// Event - Pushes an event message (ex: for status) back to the console running the module.
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
