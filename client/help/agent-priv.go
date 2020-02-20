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
	privHelp = fmt.Sprintf(`%s%s Implant Privileges Commands%s 

%s About:%s Manage privileges for implant/target 
        Type 'help <command>' for command-specific help.

%s Commands:%s
    getsystem <proc>    %s(Windows Only) Spawns a new ghost session as the NT AUTHORITY\SYSTEM user,%s
                        %sinjected in process <proc>%s
    elevate             %s(Windows Only) Spawn a new sliver session as an elevated process (UAC bypass)%s
    run_as              %sRun a process as user, with optional arguments%s
    impersonate         %s(Windows only) Run a process as user, with optional arguments%s
    rev_to_self         %sTerminate any impersonation that you have actively enabled%s`,
		tui.BLUE, tui.BOLD, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)
)
