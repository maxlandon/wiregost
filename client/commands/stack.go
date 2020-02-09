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
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"

	"github.com/maxlandon/wiregost/client/help"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func RegisterStackCommands() {

	stack := &Command{
		Name: "stack",
		Help: help.GetHelpFor("stack"),
		SubCommands: []string{
			"use",
			"pop",
		},
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				stackList(r.context.CurrentWorkspace.ID)
				fmt.Println()
			case length >= 1:
				switch r.Args[0] {
				case "use":
					if len(r.Args) == 1 {
						fmt.Println()
						fmt.Printf("%s[!]%s Provide a module path name",
							tui.RED, tui.RESET)
						fmt.Println()
					} else {
						stackUse(*r.context, r.Args[1], r.context.Server.RPC)
					}
				case "pop":
					if len(r.Args) == 1 {
						fmt.Println()
						stackPop(r.context.CurrentWorkspace.ID,
							*r.context.CurrentModule, false, r.context.CurrentModule)
						fmt.Println()
					}
					if len(r.Args) >= 2 {
						switch r.Args[1] {
						case "all":
							fmt.Println()
							stackPop(r.context.CurrentWorkspace.ID, "", true, r.context.CurrentModule)
							fmt.Println()
						default:
							fmt.Println()
							stackPop(r.context.CurrentWorkspace.ID, r.Args[1], false, r.context.CurrentModule)
							fmt.Println()
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
		Path:        strings.Split(module, "/"),
		Action:      "use",
		WorkspaceID: uint32(ctx.CurrentWorkspace.ID),
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgStackUse,
		Data: mod,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf("%s[!]%s %s", tui.RED, tui.RESET, resp.Err)
		return
	}

	stack := &clientpb.Stack{}
	proto.Unmarshal(resp.Data, stack)
	if stack.Err != "" {
		fmt.Println()
		fmt.Printf("%s[!]%s %s", tui.RED, tui.RESET, stack.Err)
		fmt.Println()
		return
	}

	currentMod := stack.Modules[0]
	*ctx.CurrentModule = strings.Join(currentMod.Path, "/")
	ctx.Module.ParseProto(currentMod)
}

func stackList(workspaceID uint) {

}

func stackPop(workspaceID uint, module string, all bool, current *string) {

}
