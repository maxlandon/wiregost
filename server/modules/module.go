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

package modules

import (
	"sync"

	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

// Module - All modules in Wiregost must implement this interface.
type Module interface {
	ParseMetadata() error                         // Parse module metadata
	ToProtobuf() *modulepb.Module                 // When consoles request a copy of the module
	Run(action string) (result string, err error) // Run one of the module's functions
	Option(name string) (opt *modulepb.Option)    // Get an option of this module
	CheckRequiredOptions() (ok bool, err error)   // Check all required options have a value
	Event(event string, pending bool)             // Send an event/message back to the console running the module
	Asset(string) (filePath string, err error)    // Find the path of an asset in the module directory.
}

// Modules - Map of all modules available in Wiregost (map["path/to/module"] = Module)
var Modules = &modules{
	Loaded: &map[string]Module{},
	mutex:  &sync.Mutex{},
}

// modules - A struct to handle all registered modules
type modules struct {
	Loaded *map[string]Module
	mutex  *sync.Mutex
}
