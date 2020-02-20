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

package help

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
)

var (
	procHelp = fmt.Sprintf(`%s%s Implant Process Commands%s 

%s About:%s Manage processes in target
        Type 'help <command>' for command-specific help.

%s Commands:%s
    ps                      %sList processes on remote system%s
    procdump <pid>|<name>   %sDumps the process memory given a process identifier <pid> or process name <name>%s
    migrate <pid>           %sMigrate implant into a host process with identifier <pid>%s
    terminate <pid>         %sTerminates a process given a process identifier <pid>%s`,
		tui.BLUE, tui.BOLD, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)

	procdumpHelp = fmt.Sprintf(`%s%sCommand:%s procdump pid=<pid> name=<name> timeout=60%s 

%s About:%s Dumps the process memory given a process identifier (pid) 

%s Filters:%s
    name        %sTarget process name%s
    pid         %sTarget process ID%s
    timeout     %sOptional command timeout in seconds (default:360)%s`,
		tui.BLUE, tui.BOLD, tui.RESET, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)
)
