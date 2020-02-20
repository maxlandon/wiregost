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
	executeHelp = fmt.Sprintf(`%s%s Implant Execute Commands%s 

%s About:%s Execute programs/shellcode/assembly/MSF payloads in target
        Type 'help <command>' for command-specific help.

%s Commands:%s
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
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)

	msfInjectHelp = fmt.Sprintf(`%s%sCommand:%s msf-inject pid=<pid> lhost=<lhost> <options>%s 

%s About:%s Execute a metasploit payload in a remote process

%s Options:%s
    pid         %sTarget process ID%s
    lhost       %sListener host%s
    lport       %sListener port (default: 4444)%s
    payload     %sMSF payload (default: meterpreter_reverse_https)%s
    encoder     %sMSF encoder%s
    iterations  %siterations of the encoder (default: 1)%s`,
		tui.BLUE, tui.BOLD, tui.RESET, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)

	shellcodeHelp = fmt.Sprintf(`%s%sCommand:%s execute-shellcode <path-to-shellcode>%s 

%s About:%s Executes the given shellcode located at <path> in the implant's process 

%s Notes:%s
    %sShellcode files should be binary encoded. You can generate shellcode files by setting 
    the 'Format' option to 'shellcode' in payload modules.%s`,
		tui.BLUE, tui.BOLD, tui.RESET, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
	)

	assemblyHelp = fmt.Sprintf(`%s%sCommand:%s execute-shellcode <path-to-shellcode>%s 

%s About:%s Executes the given shellcode located at <path> in the implant's process 

%s Notes:%s
    %sShellcode files should be binary encoded. You can generate shellcode files by setting 
    the 'Format' option to 'shellcode' in payload modules.%s`,
		tui.BLUE, tui.BOLD, tui.RESET, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
	)

	sideloadHelp = fmt.Sprintf(`%s%sCommand:%s sideload <path> timeout=60%s 

%s About:%s (Windows only) Load and execute a DLL at <path> in a remote process

%s Options:%s
    timeout     %sOptional command timeout in seconds (default: 10)%s`,
		tui.BLUE, tui.BOLD, tui.RESET, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
	)

	spawndllHelp = fmt.Sprintf(`%s%sCommand:%s spawn_dll <path> <options>%s 

%s About:%s (Windows only) Load and execute a Reflective DLL locally at <path> in a remote process

%s Options:%s
    proc        %sPath to process to host the shellcode (default: c:\windows\system32\notepad.exe)%s
    export      %sEntrypoint of the Reflective DLL (default: ReflectiveLoader)%s
    timeout     %sOptional command timeout in seconds (default: 10)%s`,
		tui.BLUE, tui.BOLD, tui.RESET, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)
)
