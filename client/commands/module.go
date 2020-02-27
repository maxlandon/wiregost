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

	"github.com/maxlandon/wiregost/client/util"
	. "github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func RegisterModuleCommands() {

	moduleUse := &Command{
		Name: "use",
		Handle: func(r *Request) error {
			if len(r.Args) == 0 {
				fmt.Println()
				fmt.Printf("\n" + Error + "Provide a module path name\n")
			} else {
				stackUse(*r.context, r.Args[0], r.context.Server.RPC)
			}
			return nil
		},
	}

	AddCommand("main", moduleUse)
	AddCommand("module", moduleUse)

	info := &Command{
		Name: "info",
		Handle: func(r *Request) error {
			fmt.Println()
			showInfo(*r.context)
			fmt.Println()
			return nil
		},
	}

	AddCommand("module", info)

	options := &Command{
		Name: "options",
		Handle: func(r *Request) error {
			fmt.Println()
			showOptions(*r.context)
			fmt.Println()
			return nil
		},
	}

	AddCommand("module", options)

	setOptions := &Command{
		Name: "set",
		Handle: func(r *Request) error {
			fmt.Println()
			setOption(r.Args, *r.context, r.context.Server.RPC)
			fmt.Println()
			return nil
		},
	}

	AddCommand("module", setOptions)

	run := &Command{
		Name: "run",
		Handle: func(r *Request) error {
			fmt.Println()
			runModule("run", *r.context, r.context.Server.RPC)
			fmt.Println()
			return nil
		},
	}

	AddCommand("module", run)

	listener := &Command{
		Name: "to_listener",
		Handle: func(r *Request) error {
			fmt.Println()
			runModule("to_listener", *r.context, r.context.Server.RPC)
			fmt.Println()
			return nil
		},
	}

	AddCommand("module", listener)

	back := &Command{
		Name: "back",
		Handle: func(r *Request) error {
			backToMainMenu(*r.context)
			return nil
		},
	}

	AddCommand("module", back)

	parseProfile := &Command{
		Name: "parse_profile",
		Handle: func(r *Request) error {
			fmt.Println()
			if len(r.Args) == 0 {
				fmt.Printf(Error + "Provide a Ghost profile name")
				return nil
			}
			parseProfile(r.Args[0], *r.context, r.context.Server.RPC)
			fmt.Println()
			return nil
		},
	}

	AddCommand("module", parseProfile)

	toProfile := &Command{
		Name: "to_profile",
		Handle: func(r *Request) error {
			fmt.Println()
			if len(r.Args) == 0 {
				fmt.Printf(Error + "Provide a profile name (to_profile <name>)\n")
				return nil
			}
			toProfile(r.Args[0], *r.context, r.context.Server.RPC)
			fmt.Println()
			return nil
		},
	}

	AddCommand("module", toProfile)
}

func setOption(args []string, ctx ShellContext, rpc RPCServer) {

	if len(args) < 2 {
		fmt.Printf(Error + "Option name/value pair not provided")
		return
	}

	name := strings.TrimSpace(args[0])

	if _, found := ctx.Module.Options[name]; !found {
		fmt.Printf(Error+"Invalid option: %s", args[0])
		return
	}

	opt, _ := proto.Marshal(&clientpb.SetOptionReq{
		WorkspaceID: uint32(ctx.CurrentWorkspace.ID),
		User:        ctx.Server.Config.User,
		Path:        strings.Split(*ctx.CurrentModule, "/"),
		Name:        name,
		Value:       args[1],
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgOptionReq,
		Data: opt,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError, "%s\n", resp.Err)
		return
	}

	changed := ctx.Module.Options[args[0]]
	changed.Value = args[1]

	fmt.Printf("%s*%s %s => %s",
		tui.BLUE, tui.RESET, name, args[1])
}

func showInfo(ctx ShellContext) {
	m := ctx.Module

	// Info
	fmt.Printf("%sModule:%s\r\t\t%s\r\n", tui.YELLOW, tui.RESET, m.Name)
	fmt.Printf("%sPlatform:%s \t%s (%s)\r\n", tui.YELLOW, tui.RESET, m.Platform, m.Targets)
	fmt.Printf("%sModule Authors:%s ", tui.YELLOW, tui.RESET)
	for a := range m.Author {
		fmt.Printf("%s ", m.Author[a])
	}
	fmt.Println()
	fmt.Printf("%sCredits:%s \t", tui.YELLOW, tui.RESET)
	for c := range m.Credits {
		fmt.Printf("%s ", m.Credits[c])
	}
	fmt.Println()
	fmt.Printf("%sLanguage:%s\r\t\t%s\n", tui.YELLOW, tui.RESET, m.Lang)
	fmt.Printf("%sPriviledged:%s \t%t\n", tui.YELLOW, tui.RESET, m.Priviledged)
	fmt.Println()
	fmt.Printf("%sDescription:%s\r\n", tui.YELLOW, tui.RESET)
	fmt.Println(tui.Dim(util.Wrap(m.Description, 100)))
	fmt.Println()

	// Options
	util.PrintOptions(m)

	// Notes
	if m.Notes != "" {
		fmt.Println()
		fmt.Printf("%sNotes:%s ", tui.YELLOW, tui.RESET)
		fmt.Println(tui.Dim(util.Wrap(m.Notes, 100)))
	}
}

func showOptions(ctx ShellContext) {
	m := ctx.Module

	util.PrintOptions(m)
}

func runModule(action string, ctx ShellContext, rpc RPCServer) {
	m := ctx.Module

	run, _ := proto.Marshal(&clientpb.ModuleActionReq{
		WorkspaceID:     uint32(ctx.CurrentWorkspace.ID),
		User:            ctx.Server.Config.User,
		Path:            m.Path,
		Action:          action,
		ModuleRequestID: int32(*ctx.ModuleRequestID),
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgModuleReq,
		Data: run,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s", resp.Err)
		return
	}

	result := &clientpb.ModuleAction{}
	proto.Unmarshal(resp.Data, result)

	if result.Success == false {
		fmt.Printf(Error+"%s", result.Err)
	} else {
		fmt.Printf(Success+"%s", result.Result)
	}

}

func parseProfile(profile string, ctx ShellContext, rpc RPCServer) {
	m := ctx.Module

	run, _ := proto.Marshal(&clientpb.ModuleActionReq{
		WorkspaceID: uint32(ctx.CurrentWorkspace.ID),
		User:        ctx.Server.Config.User,
		Path:        m.Path,
		Action:      "parse_profile",
		Profile:     profile,
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgModuleReq,
		Data: run,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s", resp.Err)
		return
	}

	result := &clientpb.ModuleAction{}
	proto.Unmarshal(resp.Data, result)

	if result.Success == false {
		fmt.Printf(Error+"%s", result.Err)
	} else {
		*m = *result.Updated
		// m.ParseProto(result.Updated)
		fmt.Printf(Info+"%s", result.Result)
	}
}

func toProfile(profile string, ctx ShellContext, rpc RPCServer) {
	m := ctx.Module

	run, _ := proto.Marshal(&clientpb.ModuleActionReq{
		WorkspaceID: uint32(ctx.CurrentWorkspace.ID),
		User:        ctx.Server.Config.User,
		Path:        m.Path,
		Action:      "to_profile",
		Profile:     profile,
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgModuleReq,
		Data: run,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s", resp.Err)
		return
	}

	result := &clientpb.ModuleAction{}
	proto.Unmarshal(resp.Data, result)

	if result.Success == false {
		fmt.Printf(Error+"%s", result.Err)
	} else {
		fmt.Printf(Info+"%s", result.Result)
	}
}

func backToMainMenu(ctx ShellContext) {
	*ctx.CurrentModule = ""
	ctx.Module = nil
}
