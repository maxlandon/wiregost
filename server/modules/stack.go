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
)

var (
	// Stacks - All modules stacks (one per workspace & per user)
	Stacks = &map[uint32]map[string]*stack{}
)

// stack - A single module stack, owned by a user for a given workspace
type stack struct {
	Loaded *map[string]Module
	mutex  *sync.Mutex
}

// LoadModule - Loads a modules onto a user's stack. Should take into account the
// console ID, so that module events are redirected to the good console.
func (s *stack) LoadModule() error {
	return nil
}

// PopModule - Unload a module from the stack
func (s *stack) PopModule(path string) error {
	return nil
}

// InitStacks - Instantiates empty stacks for each registered user, in each workspace.
func InitStacks() error {
	return nil
}
