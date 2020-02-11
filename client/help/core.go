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
	helpHelp = fmt.Sprintf(`%s%s Command Categories%s 

%s About:%s Type 'help <category>' for category-specific help.

%s Categories:%s
    core            %sCore shell commands%s
    workspace       %sManage Wiregost workspaces%s
    hosts           %sManage database hosts%s`,
		tui.BLUE, tui.BOLD, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)

	coreHelp = fmt.Sprintf(`%s%s Core Commands%s 

%s About:%s Core shell commands.

%s Commands:%s
    exit                        %sExit console%s
    ! <args>                    %sExecute a shell command through the console (bin/sh is used)%s
    cd                          %sChange the shell's current working directory%s
    resource make|load <file>   %sMake a resource file with commands, or load and execute one%s`,
		tui.BLUE, tui.BOLD, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)

	shellHelp = fmt.Sprintf(`%s%s Command:%s ! <args>%s

%s About:%s Execute a shell command locally and print its output.

%s Examples:%s
    ! ls -al            %sPrint all current working directory contents with file infos.%s
    ! cat file.txt      %sPrint file.txt%s`,
		tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)

	cdHelp = fmt.Sprintf(`%s%s Command:%s cd <dir>%s

%s About:%s Change the client shell's working directory.`,
		tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
		tui.YELLOW, tui.RESET,
	)

	resourceHelp = fmt.Sprintf(`%s%s Command:%s resource <verb> <args>%s

%s About:%s Make a resource file with past <int> commands in history, or load and execute commands of a file.

%s Examples:%s
    resource make filename=resource.rc length=10    %sMake a resouce file named resource.rc (can omit extension)
                                                    with the last ten commands in shell history.%s
    resource load resource.rc                       %sLoads the resource.rc file and executes each command in it
                                                    (resources are tab-completed)%s`,
		tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)
)
