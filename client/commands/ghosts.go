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

	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func registerGhostBuildsCommands() {

	ghosts := &Command{
		Name: "ghosts",
		Handle: func(r *Request) error {
			fmt.Println()
			listGhostBuilds(*r.context, r.context.Server.RPC)
			return nil
		},
	}

	AddCommand("main", ghosts)
	AddCommand("module", ghosts)

	canaries := &Command{
		Name: "canaries",
		Handle: func(r *Request) error {
			fmt.Println()
			listCanaries(*r.context, r.context.Server.RPC)
			return nil
		},
	}

	AddCommand("main", canaries)
	AddCommand("module", canaries)
}

func listGhostBuilds(ctx ShellContext, rpc RPCServer) {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgListGhostBuilds,
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}

	builds := &clientpb.GhostBuilds{}
	proto.Unmarshal(resp.Data, builds)
	if 0 < len(builds.Configs) {
		displayAllGhostBuilds(builds.Configs)
	} else {
		fmt.Printf(Info + "No ghost builds\n")
	}
}

func displayAllGhostBuilds(configs map[string]*clientpb.GhostConfig) {

	table := util.Table()
	table.SetHeader([]string{"WsID", "Name", "Platform", "Format", "Command & Control", "Limits", "Debug"})
	table.SetColWidth(40)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	for k, c := range configs {
		platform := fmt.Sprintf("%s/%s", c.GOOS, c.GOARCH)
		c2s := []string{}
		for _, c := range c.C2 {
			c2s = append(c2s, c.URL)
		}
		workspace := ""
		if c.WorkspaceID != 0 {
			workspace = strconv.Itoa(int(c.WorkspaceID))
		}
		var format string
		if c.Format == clientpb.GhostConfig_EXECUTABLE {
			format = "exe"
		}
		if c.Format == clientpb.GhostConfig_SHARED_LIB {
			format = "shared"
		}
		if c.Format == clientpb.GhostConfig_SHELLCODE {
			format = "shellcode"
		}

		limits := getLimitsString(c)
		table.Append([]string{workspace, k, platform, format,
			strings.Join(c2s, ","), limits, strconv.FormatBool(c.Debug)})
	}

	table.Render()
}

func listCanaries(ctx ShellContext, rpc RPCServer) {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgListCanaries,
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}

	canaries := &clientpb.Canaries{}
	proto.Unmarshal(resp.Data, canaries)
	if 0 < len(canaries.Canaries) {
		displayCanaries(canaries.Canaries)
	} else {
		fmt.Printf(Info + "No canaries in database\n")
	}
}

func displayCanaries(canaries []*clientpb.DNSCanary) {

	table := util.Table()
	table.SetHeader([]string{"Name", "Domain", "Triggered", "First trigger", "Latest trigger"})
	table.SetColWidth(40)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	for _, c := range canaries {
		triggered := ""
		if c.Triggered {
			triggered = tui.Red(tui.Bold("Yes"))
		} else {
			triggered = tui.Green("No")
		}

		table.Append([]string{c.GhostName, c.Domain, triggered, c.FirstTriggered, c.LatestTrigger})
	}

	table.Render()
}
