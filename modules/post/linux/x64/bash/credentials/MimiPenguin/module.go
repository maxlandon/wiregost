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

package MimiPenguin

import (
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module/templates"
)

// [ Base Methods ] ------------------------------------------------------------------------//

// MimiPenguin - A single stage DNS implant
type MimiPenguin struct {
	*templates.Module
}

// New - Instantiates a reverse DNS module, empty.
func New() *MimiPenguin {
	mod := &MimiPenguin{&templates.Module{}}
	mod.Path = []string{"post/linux/x64/bash/credentials/MimiPenguin"}
	return mod
}

var modLog = log.ServerLogger("post/linux/x64/bash/credentials/MimiPenguin", "module")

// [ Module Methods ] ------------------------------------------------------------------------//
func (s *MimiPenguin) Run(command string) (result string, err error) {

	// Check options
	if ok, err := s.CheckRequiredOptions(); !ok {
		return "", err
	}

	// Check session
	sess, err := s.CheckSession()
	if sess == nil {
		return "", err
	}

	return result, nil
}
