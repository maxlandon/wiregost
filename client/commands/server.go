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

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/assets"
	"github.com/maxlandon/wiregost/client/help"
	"github.com/maxlandon/wiregost/client/util"
	"github.com/olekukonko/tablewriter"
)

func RegisterServerCommands() {

	server := &Command{
		Name: "server",
		Help: help.GetHelpFor("server"),
		SubCommands: []string{
			"connect",
		},
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				listServers(*r.context)
			case length >= 1:
				switch r.Args[0] {
				case "connect":
				}
			}
			return nil
		},
	}

	AddCommand("main", server)
	AddCommand("module", server)
}

func listServers(ctx ShellContext) {

	configs := assets.GetConfigs()
	if len(configs) == 0 {
		fmt.Printf("%s[!] No config files found at %s or -config\n", tui.YELLOW, assets.GetConfigDir())
		return
	}

	table := util.Table()
	table.SetHeader([]string{"User Name", "LHost", "LPort", "Default", "Connected"})
	table.SetColWidth(40)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	for _, c := range configs {
		def := ""
		if c.IsDefault {
			def = "default"
		}
		port := strconv.Itoa(c.LPort)

		connected := ""
		if c.LHost == ctx.CurrentServer.LHost && c.LPort == ctx.CurrentServer.LPort {
			connected = fmt.Sprintf("%sconnected%s", tui.GREEN, tui.RESET)
		}

		table.Append([]string{c.User, c.LHost, port, def, connected})
	}
	table.Render()
}
