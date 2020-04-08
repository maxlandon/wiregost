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
	"sort"
	"strings"

	"github.com/gogo/protobuf/proto"
	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/olekukonko/tablewriter"
)

// StackCmd - List modules on the stack
type StackCmd struct{}

var Stack StackCmd

func RegisterStack() {
	CommandParser.AddCommand(constants.Stack, "", "", &Stack)

	stack := CommandParser.Find(constants.Stack)
	CommandMap[MAIN_CONTEXT] = append(CommandMap[MAIN_CONTEXT], stack)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], stack)
	stack.ShortDescription = "List modules currently loaded on the stack"
	stack.SubcommandsOptional = true
}

// Execute - List modules on the stack
func (s *StackCmd) Execute(args []string) error {

	stack, _ := proto.Marshal(&clientpb.StackReq{
		Action:      constants.StackList,
		WorkspaceID: uint32(Context.Workspace.ID),
		User:        Context.Server.Config.User,
	})

	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgStackList,
		Data: stack,
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	stackList := &clientpb.Stack{}
	proto.Unmarshal(resp.Data, stackList)
	if stackList.Err != "" {
		fmt.Printf(Error+"%s \n", stackList.Err)
		return nil
	}

	tab := util.NewTable()
	headers := []string{"Type", "Path"}
	widths := []int{15, 30}
	tab.SetColumns(headers, widths)

	tab.SetColWidth(60)

	list := []string{}
	for _, m := range stackList.Modules {
		list = append(list, strings.Join(m.Path, "/"))
	}
	list = sort.StringSlice(list)

	for _, p := range list {
		for _, m := range stackList.Modules {
			if strings.Join(m.Path, "/") == p {
				if strings.Join(m.Path, "/") == strings.Join(Context.Module.Path, "/") {
					tab.Rich([]string{m.Type, strings.Join(m.Path[1:], "/")},
						[]tablewriter.Colors{tablewriter.Colors{tablewriter.Normal, tablewriter.FgBlueColor},
							tablewriter.Colors{tablewriter.Normal, tablewriter.FgBlueColor},
						})
				} else {
					tab.Append([]string{m.Type, strings.Join(m.Path[1:], "/")})
				}
			}
		}
	}
	tab.Output()

	return nil
}
