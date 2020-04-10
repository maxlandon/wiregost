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

package stack

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

type StackPopCmd struct {
	Positional struct {
		Path string `description:"Module path"`
	} `positional-args:"yes"`
}

var StackPop StackPopCmd

func RegisterStackPop() {
	stack := MainParser.Find(constants.Stack)

	stack.AddCommand(constants.StackPop, "", "", &StackPop)
}

func (sp *StackPopCmd) Execute(args []string) error {

	var all = false
	if sp.Positional.Path == "all" {
		all = true
	}

	var module []string
	if sp.Positional.Path == "" {
		module = Context.Module.Path
	} else {
		module = strings.Split(sp.Positional.Path, "/")
	}

	mod, _ := proto.Marshal(&clientpb.StackReq{
		Path:        module,
		All:         all,
		WorkspaceID: uint32(Context.Workspace.ID),
		User:        Context.Server.Config.User,
	})

	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgStackPop,
		Data: mod,
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	stack := &clientpb.Stack{}
	proto.Unmarshal(resp.Data, stack)
	if stack.Err != "" {
		fmt.Printf(Error+"%s \n", stack.Err)
		return nil
	}

	if all {
		Context.Module = &clientpb.Module{}
		return nil
	}

	if (stack.Path != nil) && (len(stack.Path) != 0) {
		*Context.Module = *stack.Modules[0]
		return nil
	}

	if len(stack.Path) == 0 {
		*Context.Module = clientpb.Module{}
	}

	return nil
}
