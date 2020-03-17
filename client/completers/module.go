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
	"sort"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/help"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

type moduleCompleter struct {
	Command *commands.Command
}

// Do is the completion function triggered at each line
func (mc *moduleCompleter) Do(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {

	splitLine := strings.Split(string(line), " ")
	line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))

	line = append(line, []rune(help.GetHelpFor("module"))...)
	stack, _ := proto.Marshal(&clientpb.ModuleActionReq{
		Action: "list",
	})

	rpc := ctx.Server.RPC

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgModuleList,
		Data: stack,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError, "%s\n", resp.Err)
		return
	}

	modList := &clientpb.ModuleAction{}
	proto.Unmarshal(resp.Data, modList)

	sort.Strings(modList.Modules)
	for _, mod := range modList.Modules {
		search := mod
		if !hasPrefix(line, []rune(search)) {
			sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
			options = append(options, sLine...)
			offset = sOffset
		}
	}

	return options, offset
}
