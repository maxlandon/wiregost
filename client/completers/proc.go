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

package completers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"
	"github.com/lmorg/readline"

	. "github.com/maxlandon/wiregost/client/commands"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func CompleteProcesses(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get last path
	splitLine := strings.Split(string(line), " ")
	last := splitLine[len(splitLine)-1]

	data, _ := proto.Marshal(&ghostpb.PsReq{GhostID: Context.Ghost.ID})
	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: ghostpb.MsgPsReq,
		Data: data,
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return string(last), suggestions, listSuggestions, readline.TabDisplayList
	}
	ps := &ghostpb.Ps{}
	err := proto.Unmarshal(resp.Data, ps)
	if err != nil {
		fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
		return string(last), suggestions, listSuggestions, readline.TabDisplayList
	}

	for _, proc := range ps.Processes {
		pid := strconv.Itoa(int(proc.Pid))

		// Check if input is a string or an int:
		_, err := strconv.Atoi(string(last))

		// If there is an error, user searches based on name
		// if err != nil && string(last) != "" {
		//         if strings.HasPrefix(proc.Executable, string(last)) {
		//                 suggestions = append(suggestions, pid)
		//                 listSuggestions[pid] = proc.Executable
		//         }
		//
		// } else
		// If no error, user searches by PID, or if nothing in input yet.
		if (err == nil || err != nil) && (string(last) == "" || string(last) != "") {
			if strings.HasPrefix(pid, string(last)) {
				suggestions = append(suggestions, pid[(len(last)):])

				// Get parent proc
				var parent *ghostpb.Process
				for _, p := range ps.Processes {
					if p.Pid == proc.Pid {
						parent = proc
					}
				}

				exe := exePad(proc, ps.Processes)
				ownerDesc := ownerPad(proc, ps.Processes)
				parentDesc := parentPad(parent, ps.Processes)

				desc := proc.Executable + tui.Dim(exe+"<--    "+parentDesc+"<--   "+ownerDesc)
				listSuggestions[pid[(len(last)):]] = desc
			}
		}
	}

	// return "", suggestions, listSuggestions, readline.TabDisplayList
	return string(last), suggestions, listSuggestions, readline.TabDisplayList
}

func exePad(child *ghostpb.Process, procs []*ghostpb.Process) string {
	var max int
	for _, proc := range procs {
		if len(proc.Executable) > max {
			max = len(proc.Executable)
		}
	}
	var pad string
	for i := 0; i < max-len(child.Executable); i++ {
		pad += " "
	}

	return pad
}

func ownerPad(child *ghostpb.Process, procs []*ghostpb.Process) string {
	var max int
	for _, proc := range procs {
		if len(proc.Owner) > max {
			max = len(proc.Owner)
		}
	}
	var pad string
	for i := 0; i < max-len(child.Owner); i++ {
		pad += " "
	}

	return child.Owner + pad
}

func parentPad(parent *ghostpb.Process, procs []*ghostpb.Process) string {
	var max int
	var str string
	for _, proc := range procs {
		str = fmt.Sprintf(strconv.Itoa(int(proc.Pid)) + "/" + proc.Executable)
		if len(str) > max {
			max = len(str)
		}
	}
	var pad string
	parStr := strconv.Itoa(int(parent.Ppid)) + "/" + parent.Executable
	for i := 0; i < max-len(parStr); i++ {
		pad += " "
	}

	return parStr + pad
}
