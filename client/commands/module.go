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
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"
	"github.com/olekukonko/tablewriter"

	"github.com/maxlandon/wiregost/client/help"
	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func RegisterModuleCommands() {

	moduleUse := &Command{
		Name: "use",
		Help: help.GetHelpFor("use"),
		Handle: func(r *Request) error {
			if len(r.Args) == 0 {
				fmt.Println()
				fmt.Printf("%s[!]%s Provide a module path name",
					tui.RED, tui.RESET)
				fmt.Println()
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
		Help: help.GetHelpFor("info"),
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
		Help: help.GetHelpFor("options"),
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
		Help: help.GetHelpFor("set"),
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
		Help: help.GetHelpFor("run"),
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
		Help: help.GetHelpFor("to_listener"),
		Handle: func(r *Request) error {
			fmt.Println()
			runModule("to_listener", *r.context, r.context.Server.RPC)
			fmt.Println()
			return nil
		},
	}

	AddCommand("module", listener)
}

func setOption(args []string, ctx ShellContext, rpc RPCServer) {

	if len(args) < 2 {
		fmt.Printf("%s[!]%s Option name/value pair not provided",
			tui.RED, tui.RESET)
		return
	}

	name := strings.TrimSpace(args[0])

	if _, found := ctx.Module.Options[name]; !found {
		fmt.Printf("%s[!]%s Invalid option: %s",
			tui.RED, tui.RESET, args[0])
		return
	}

	opt, _ := proto.Marshal(&clientpb.SetOptionReq{
		WorkspaceID: uint32(ctx.CurrentWorkspace.ID),
		Path:        strings.Split(*ctx.CurrentModule, "/"),
		Name:        name,
		Value:       args[1],
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgOptionReq,
		Data: opt,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf("%s[!]%s %s", tui.RED, tui.RESET, resp.Err)
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
	fmt.Println(tui.Dim(util.Wrap(m.Description, 140)))
	fmt.Println()

	// Listener Options
	fmt.Println(tui.Bold(tui.Blue(" Listener Options")))
	table := util.Table()
	table.SetHeader([]string{"Name", "Value", "Required", "Description"})
	table.SetColWidth(90)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)
	for _, v := range util.SortListenerOptionKeys(m.Options) {
		table.Append([]string{m.Options[v].Name, m.Options[v].Value, strconv.FormatBool(m.Options[v].Required), m.Options[v].Description})
	}
	table.Render()

	// Generate Options
	fmt.Println()
	fmt.Println(tui.Bold(tui.Blue(" Generate Options")))
	table = util.Table()
	table.SetHeader([]string{"Name", "Value", "Required", "Description"})
	table.SetColWidth(90)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)
	for _, v := range util.SortGenerateOptionKeys(m.Options) {
		table.Append([]string{m.Options[v].Name, m.Options[v].Value, strconv.FormatBool(m.Options[v].Required), m.Options[v].Description})
	}
	table.Render()

	// Notes
	if m.Notes != "" {
		fmt.Println()
		fmt.Printf("%sNotes:%s ", tui.YELLOW, tui.RESET)
		fmt.Println(tui.Dim(util.Wrap(m.Notes, 140)))
	}
}

func showOptions(ctx ShellContext) {
	m := ctx.Module

	// Listener Options
	fmt.Println(tui.Bold(tui.Blue(" Listener Options")))
	table := util.Table()
	table.SetHeader([]string{"Name", "Value", "Required", "Description"})
	table.SetColWidth(80)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)
	for _, v := range util.SortListenerOptionKeys(m.Options) {
		table.Append([]string{m.Options[v].Name, m.Options[v].Value, strconv.FormatBool(m.Options[v].Required), m.Options[v].Description})
	}
	table.Render()

	// Generate Options
	fmt.Println()
	fmt.Println(tui.Bold(tui.Blue(" Generate Options")))
	table = util.Table()
	table.SetHeader([]string{"Name", "Value", "Required", "Description"})
	table.SetColWidth(80)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)
	for _, v := range util.SortGenerateOptionKeys(m.Options) {
		table.Append([]string{m.Options[v].Name, m.Options[v].Value, strconv.FormatBool(m.Options[v].Required), m.Options[v].Description})
	}
	table.Render()
}

func runModule(action string, ctx ShellContext, rpc RPCServer) {
	m := ctx.Module

	run, _ := proto.Marshal(&clientpb.ModuleActionReq{
		WorkspaceID: uint32(ctx.CurrentWorkspace.ID),
		Path:        m.Path,
		Action:      action,
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgModuleReq,
		Data: run,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf("%s[!] RPC error:%s %s", tui.RED, tui.RESET, resp.Err)
		return
	}

	result := &clientpb.ModuleAction{}
	proto.Unmarshal(resp.Data, result)

	if result.Sucess == false {
		fmt.Printf("%s[!]%s %s", tui.RED, tui.RESET, result.Err)
	} else {
		fmt.Printf("%s[*]%s %s", tui.GREEN, tui.RESET, result.Result)
	}

}
