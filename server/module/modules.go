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

package module

import (
	"errors"
	"sync"

	pb "github.com/maxlandon/wiregost/protobuf/client"
)

var (
	// Modules - Map of all modules available in Wiregost (map["path/to/module"] = Module)
	Modules = &modules{
		Loaded: &map[string]Module{},
		mutex:  &sync.RWMutex{},
	}
)

type modules struct {
	Loaded *map[string]Module
	mutex  *sync.RWMutex
}

// Module - Represents a module, providing access to its methods
type Module interface {
	Init() error
	Run(string) error
	ToProtobuf() *pb.Module
}

// Module - Get module by path, (load it if needed)
func GetModule(path string) (Module, error) {

	Modules.mutex.Lock()
	defer Modules.mutex.Unlock()

	if mod, ok := (*Modules.Loaded)[path]; !ok {
		return nil, errors.New("No module for given path")
	} else {
		return mod, nil
	}
}
