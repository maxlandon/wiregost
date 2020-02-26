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
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/util"
	. "github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// AutoCompleter is the autocompletion engine
type OptionCompleter struct {
	Command *commands.Command
}

// Do is the completion function triggered at each line
func (oc *OptionCompleter) Do(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {

	switch *ctx.MenuContext {
	case "module":
		for _, v := range util.SortListenerOptionKeys(ctx.Module.Options) {
			search := v + " "
			if !hasPrefix(line, []rune(search)) {
				sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}

		for _, v := range util.SortGenerateOptionKeys(ctx.Module.Options) {
			search := v + " "
			if !hasPrefix(line, []rune(search)) {
				sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}
	}

	// Else, provide some specific option values:
	words := strings.Split(string(line), " ")
	argInput := lastString(words)

	// For some arguments, the split results in a last empty item.
	if words[len(words)-1] == "" {
		argInput = words[0]
	}
	if argInput == "StageImplant" {
	}

	switch argInput {
	case "StageImplant", "StageConfig":
		// Get ghost builds
		rpc := ctx.Server.RPC
		resp := <-rpc(&ghostpb.Envelope{
			Type: clientpb.MsgListGhostBuilds,
		}, defaultTimeout)
		if resp.Err != "" {
			fmt.Printf(RPCError+"%s\n", resp.Err)
			return
		}

		builds := &clientpb.GhostBuilds{}
		proto.Unmarshal(resp.Data, builds)
		shellcodeBuilds := []*clientpb.GhostConfig{}
		for _, c := range builds.Configs {
			if (c.Format == clientpb.GhostConfig_SHARED_LIB) || (c.Format == clientpb.GhostConfig_SHELLCODE) {
				shellcodeBuilds = append(shellcodeBuilds, c)
			}
		}

		for _, c := range shellcodeBuilds {
			options = append(options, []rune(c.Name))
			offset = len(argInput + " ")
			// search := c.Name
			// if !hasPrefix(line, []rune(search)) {
			//         sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
			//         options = append(options, sLine...)
			//         offset = sOffset
			// } else {
			//
			// }
		}
		return
	}

	return options, offset
}
