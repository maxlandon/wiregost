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
	"sync"

	"github.com/maxlandon/wiregost/data_service/remote"
)

var (
	// Stacks - All module stacks (one per workspace),
	// which can be loaded/unloaded on demand, pulling from Modules
	Stacks = &map[uint]*stack{}
)

type stack struct {
	Loaded *map[string]Module
	mutex  *sync.RWMutex
}

// InitStacks - Creates a new stack for each workspace in Wiregost
func InitStacks() {
	workspaces, _ := remote.Workspaces(nil)
	for _, w := range workspaces {
		(*Stacks)[w.ID] = &stack{}
	}

	for _, v := range *Stacks {
		v.Loaded = &map[string]Module{}
		v.mutex = &sync.RWMutex{}
	}
}

// LoadModule - Load a module onto the stack, by fetching it into Modules
func (s *stack) LoadModule(path string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	mod, err := GetModule(path)
	if err != nil {
		return err
	}

	// Init and load onto stack
	mod.Init()
	(*s.Loaded)[path] = mod

	return nil
}

// PopModule - Unload a module from stack
func (s *stack) PopModule(path string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(*s.Loaded, path)

	return nil
}

// Module - Get a module by path, (load it onto the stack if needed)
func (s *stack) Module(path string) (Module, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if mod, ok := (*s.Loaded)[path]; !ok {
		s.LoadModule(path)
		return (*s.Loaded)[path], nil
	} else {
		return mod, nil
	}
}
