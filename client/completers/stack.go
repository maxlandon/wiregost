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
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/maxlandon/wiregost/client/commands"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

var defaultTimeout = 30 * time.Second

type stackCompleter struct {
	Command *commands.Command
}

// Do is the completion function triggered at each line
func (mc *stackCompleter) Do(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {

	splitLine := strings.Split(string(line), " ")
	line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))

	switch splitLine[0] {
	case "use", "pop":
		// Get stack modules
		stack, _ := proto.Marshal(&clientpb.StackReq{
			Action:      "list",
			WorkspaceID: uint32(ctx.Workspace.ID),
			User:        ctx.Server.Config.User,
		})

		rpc := ctx.Server.RPC

		resp := <-rpc(&ghostpb.Envelope{
			Type: clientpb.MsgStackList,
			Data: stack,
		}, defaultTimeout)

		if resp.Err != "" {
			fmt.Printf(RPCError, "%s\n", resp.Err)
			return
		}

		stackList := &clientpb.Stack{}
		proto.Unmarshal(resp.Data, stackList)
		if stackList.Err != "" {
			fmt.Println()
			fmt.Printf(Error, "%s", stackList.Err)
			fmt.Println()
			return
		}

		for _, mod := range stackList.Modules {
			search := strings.Join(mod.Path, "/")
			if !hasPrefix(line, []rune(search)) {
				sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}
	}

	return options, offset
}
