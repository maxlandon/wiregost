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

// StackUse - Use a module currently on the stack
type StackUseCmd struct {
	Positional struct {
		Path string `description:"Module path" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var StackUse StackUseCmd

func RegisterStackUse() {
	stack := MainParser.Find(constants.Stack)

	stack.AddCommand(constants.StackUse, "", "", &StackUse)
	use := stack.Find(constants.StackUse)
	use.ShortDescription = "Use a module currently loaded onto the stack"
	use.Args()[0].RequiredMaximum = 1
}

// Execute - Use a module currently on the stack
func (su *StackUseCmd) Execute(args []string) error {

	mod, _ := proto.Marshal(&clientpb.StackReq{
		Path:         strings.Split(su.Positional.Path, "/"),
		Action:       constants.StackUse,
		WorkspaceID:  uint32(Context.Workspace.ID),
		User:         Context.Server.Config.User,
		ModuleUserID: Context.UserID,
	})

	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgStackUse,
		Data: mod,
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	stack := &clientpb.Stack{}
	proto.Unmarshal(resp.Data, stack)
	if stack.Err != "" {
		fmt.Printf("\n"+Error+"%s\n", stack.Err)
		return nil
	}

	*Context.Module = *stack.Modules[0]

	return nil
}
