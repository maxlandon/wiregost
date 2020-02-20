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
	agentHelp = fmt.Sprintf(`%s%s Implant Command Categories%s 

%s About:%s These commands are available when using a ghost implant
        Type 'help <category>' for category-specific help.

%s Categories:%s
    filesystem      %sFile system management (ls,mkdir,download, upload, etc....)%s
    info            %sInformation commands (network, implant...)%s
    priv            %sManage target privileges (impersonate, getsystem, run as...)%s
    proc            %sManage target processes (ps, procdump, migrate, terminate...)%s
    shell           %sHelp for target shell usage%s
    execute         %sHelp for executing MSF/shellcode/assembly payloads on target%s`,
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

	filesystemHelp = fmt.Sprintf(`%s%s Implant file-system Commands%s 

%s About:%s Manage target file-system.

%s Commands:%s
    ls                          %sList files/directories in current working directory%s
    cd                          %sChange the implant's current working directory (works with relative paths)%s
    pwd                         %sPring the current working directory%s
    rm                          %sRemove a file or directory from the target%s
    mkdir <dir>                 %sCreate a directory with path <dir>%s
    download <remote> <local>   %sDownload file <remote> and save as file <local>%s
    upload <local> <remote>     %sUpload file <local> and save as file <remote>%s`,
		tui.BLUE, tui.BOLD, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)
)
