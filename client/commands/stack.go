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

package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/olekukonko/tablewriter"

	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func registerStackCommands() {

	stack := &Command{
		Name: "stack",
		// SubCommands: []string{
		//         "use",
		//         "pop",
		// },
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				stackList(*r.context, r.context.Server.RPC)
			case length >= 1:
				switch r.Args[0] {
				case "use":
					if len(r.Args) == 1 {
						fmt.Printf("\n" + Error + "Provide a module path name\n")
					} else {
						stackUse(*r.context, r.Args[1], r.context.Server.RPC)
					}
				case "pop":
					if len(r.Args) == 1 {
						stackPop(*r.context, strings.Join(r.context.Module.Path, "/"), false, r.context.Server.RPC)
					}
					if len(r.Args) >= 2 {
						switch r.Args[1] {
						case "all":
							stackPop(*r.context, "", true, r.context.Server.RPC)
						default:
							for _, arg := range r.Args[1:] {
								stackPop(*r.context, arg, false, r.context.Server.RPC)
							}
						}
					}

				}
			}
			return nil
		},
	}

	// Add commands for each context
	AddCommand("main", stack)
	AddCommand("module", stack)
	AddCommand("ghost", stack)
}

func stackUse(ctx ShellContext, module string, rpc RPCServer) {
	mod, _ := proto.Marshal(&clientpb.StackReq{
		Path:         strings.Split(module, "/"),
		Action:       "use",
		WorkspaceID:  uint32(ctx.Workspace.ID),
		User:         ctx.Server.Config.User,
		ModuleUserID: ctx.UserID,
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgStackUse,
		Data: mod,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}

	stack := &clientpb.Stack{}
	proto.Unmarshal(resp.Data, stack)
	if stack.Err != "" {
		fmt.Printf("\n"+Error+"%s\n", stack.Err)
		return
	}

	*ctx.Module = *stack.Modules[0]
	// *ctx.CurrentModule = strings.Join(ctx.Module.Path, "/")
}

func stackList(ctx ShellContext, rpc RPCServer) {
	stack, _ := proto.Marshal(&clientpb.StackReq{
		Action:      "list",
		WorkspaceID: uint32(ctx.Workspace.ID),
		User:        ctx.Server.Config.User,
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgStackList,
		Data: stack,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}

	stackList := &clientpb.Stack{}
	proto.Unmarshal(resp.Data, stackList)
	if stackList.Err != "" {
		fmt.Println()
		fmt.Printf(Error+"%s", stackList.Err)
		fmt.Println()
		return
	}

	table := util.Table()
	table.SetHeader([]string{"Type", "Path"})
	table.SetColWidth(60)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	list := []string{}
	for _, m := range stackList.Modules {
		list = append(list, strings.Join(m.Path, "/"))
	}
	list = sort.StringSlice(list)

	for _, p := range list {
		for _, m := range stackList.Modules {
			if strings.Join(m.Path, "/") == p {
				if strings.Join(m.Path, "/") == strings.Join(ctx.Module.Path, "/") {
					table.Rich([]string{m.Type, strings.Join(m.Path[1:], "/")},
						[]tablewriter.Colors{tablewriter.Colors{tablewriter.Normal, tablewriter.FgBlueColor},
							tablewriter.Colors{tablewriter.Normal, tablewriter.FgBlueColor},
						})
				} else {
					table.Append([]string{m.Type, strings.Join(m.Path[1:], "/")})
				}
			}
		}
	}
	table.Render()
}

func stackPop(ctx ShellContext, module string, all bool, rpc RPCServer) {

	mod, _ := proto.Marshal(&clientpb.StackReq{
		Path:        strings.Split(module, "/"),
		All:         all,
		WorkspaceID: uint32(ctx.Workspace.ID),
		User:        ctx.Server.Config.User,
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgStackPop,
		Data: mod,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}

	stack := &clientpb.Stack{}
	proto.Unmarshal(resp.Data, stack)
	if stack.Err != "" {
		fmt.Printf(Error+"%s", stack.Err)
		return
	}

	if all {
		// *ctx.entModule = ""
		ctx.Module = nil
		return
	}

	if (stack.Path != nil) && (len(stack.Path) != 0) {
		// *ctx.CurrentModule = strings.Join(stack.Path, "/")
		*ctx.Module = *stack.Modules[0]
		return
	}

	if len(stack.Path) == 0 {
		// *ctx.CurrentModule = ""
		ctx.Module = nil
	}
}
