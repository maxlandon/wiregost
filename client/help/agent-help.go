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
	completeAgentHelp = fmt.Sprintf(`%s%s Implant Commands%s 

%s About:%s These commands are available when using a ghost implant
        Type 'help <command>' for command-specific help.

%s Shell:%s
    shell --no-puty       %sStart an interactive shell, with optional putty disable on MacOS/Linux%s

%s File System:%s
    ls                          %sList files/directories in current working directory%s
    cd                          %sChange the implant's current working directory (works with relative paths)%s
    pwd                         %sPring the current working directory%s
    rm                          %sRemove a file or directory from the target%s
    mkdir <dir>                 %sCreate a directory with path <dir>%s
    download <remote> <local>   %sDownload file <remote> and save as file <local>%s
    upload <local> <remote>     %sUpload file <local> and save as file <remote>%s

%s Information:%s
    ifconfig        %sPrint network/interfaces information of target%s
    info            %sPrint implant/target information%s
    whoami          %sShow username/hostname of target%s
    ping            %sPing the implant. (This does NOT send an ICMP packet, but an empty C2 message)%s
    getuid          %sGet implant process User ID%s
    getpid          %sGet implant Process ID%s
    getgid          %sGet implant process Group ID%s

%s Privileges:%s
    getsystem <proc>    %s(Windows Only) Spawns a new ghost session as the NT AUTHORITY\SYSTEM user,%s
                        %sinjected in process <proc>%s
    elevate             %s(Windows Only) Spawn a new sliver session as an elevated process (UAC bypass)%s
    run_as              %sRun a process as user, with optional arguments%s
    impersonate         %s(Windows only) Run a process as user, with optional arguments%s
    rev_to_self         %sTerminate any impersonation that you have actively enabled%s

%s Processes:%s
    ps <filters>            %sList processes on remote system%s
    procdump <pid>|<name>   %sDumps the process memory given a process identifier <pid> or process name <name>%s
    migrate <pid>           %sMigrate implant into a host process with identifier <pid>%s
    terminate <pid>         %sTerminates a process given a process identifier <pid>%s

%s Execution/Injection:%s
    execute <path>      %sExecute a process located at <path> in target%s
    msf-inject          %sExecute a metasploit payload in a remote process%s
    execute-shellcode   %sExecutes the given shellcode in the implant's process%s
    execute-assembly    %s(Windows only) Executes the .NET assembly in a child process%s
    sideload            %s(Windows only) Load and execute a DLL  in a remote process%s
    spawn_dll           %s(Windows only) Load and execute a Reflective DLL in a remote process%s`,
		tui.BLUE, tui.BOLD, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)
)
