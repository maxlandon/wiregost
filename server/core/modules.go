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

package core

import (
	"errors"
	"sync"

	"github.com/maxlandon/wiregost/data_service/remote"
	pb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/certs"
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
// All Wiregost modules must implement this interface
type Module interface {
	Init(int32) error
	Run(string) (string, error)
	SetOption(string, string)
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

func AddModule(path string, mod Module) error {

	Modules.mutex.Lock()
	defer Modules.mutex.Unlock()

	(*Modules.Loaded)[path] = mod
	return nil
}

var (
	// Stacks - All module stacks (one per workspace),
	// which can be loaded/unloaded on demand, pulling from Modules
	Stacks = &map[uint]map[string]*stack{}
)

type stack struct {
	Loaded *map[string]Module
	mutex  *sync.RWMutex
}

// InitStacks - Creates a new stack for each workspace in Wiregost
func InitStacks() {
	clientCerts := certs.UserClientListCertificates()

	users := []string{}
	for _, c := range clientCerts {
		users = append(users, c.Subject.CommonName)
	}
	users = unique(users)

	workspaces, _ := remote.Workspaces(nil)
	for _, w := range workspaces {
		for _, user := range users {
			userStack := &map[string]*stack{}
			(*userStack)[user] = &stack{}
			(*Stacks)[w.ID] = (*userStack)
		}
	}

	for _, v := range *Stacks {
		for _, u := range users {
			v[u] = &stack{}
			v[u].Loaded = &map[string]Module{}
			v[u].mutex = &sync.RWMutex{}
		}
	}
}

// LoadModule - Load a module onto the stack, by fetching it into Modules
func (s *stack) LoadModule(userID int32, path string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	mod, err := GetModule(path)
	if err != nil {
		return err
	}

	// Init and load onto stack
	mod.Init(userID)
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
func (s *stack) Module(userID int32, path string) (Module, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if mod, ok := (*s.Loaded)[path]; !ok {
		s.LoadModule(userID, path)
		return (*s.Loaded)[path], nil
	} else {
		return mod, nil
	}

}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
