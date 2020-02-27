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

// CHANGE THE NAME OF THE PACKAGE WITH THE NAME OF YOUR MODULE/DIRECTORY
package main

import (
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module"
)

// [ Base Methods ] ------------------------------------------------------------------------//

// Auxiliary - An Auxiliary Module  (Change "Auxiliary")
type Auxiliary struct {
	*module.Module
}

// New - Instantiates an Auxiliary module, loading its metadata.
// - Change the field "path/to/module/directory" by something like "scanner/nmap/moduleDirName"
func New() *Auxiliary {
	mod := &Auxiliary{&module.Module{}}
	mod.Path = []string{"auxiliary", "path/to/module/directory"}
	return mod
}

var modLog = log.ServerLogger("path/to/module/directory", "module")

// [ Module Methods ] ------------------------------------------------------------------------//

// Run - Module entrypoint. ** DO NOT ERASE **
func (s *Auxiliary) Run(requestID int32, command string) (result string, err error) {

	return "", nil
}
