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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"

	"github.com/maxlandon/wiregost/client/assets"
	"github.com/maxlandon/wiregost/client/core"
	"github.com/maxlandon/wiregost/client/transport"
	"github.com/maxlandon/wiregost/client/util"
)

func registerServerCommands() {

	server := &Command{
		Name: "server",
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
					fmt.Println()
					connectServer(r.Args[1:], *r.context)
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
		fmt.Printf(Warnf+"No config files found at %s or -config\n", assets.GetConfigDir())
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
		if c.LHost == ctx.Server.Config.LHost && c.LPort == ctx.Server.Config.LPort {
			connected = fmt.Sprintf("%sconnected%s", tui.GREEN, tui.RESET)
		}

		table.Append([]string{c.User, c.LHost, port, def, connected})
	}
	table.Render()
}

func connectServer(args []string, ctx ShellContext) error {
	if len(args) == 0 {
		fmt.Printf("\n" + Warn + "Provide a server address \n")
		return nil
	}

	lhost := strings.Split(args[0], ":")[0]
	port := strings.Split(args[0], ":")[1]
	lport, _ := strconv.Atoi(port)
	user := args[2]

	configs := assets.GetConfigs()
	var config *assets.ClientConfig
	for _, conf := range configs {
		if (lhost == conf.LHost) && (lport == conf.LPort) && (user == conf.User) {
			config = conf
		}
	}

	fmt.Printf(Warn+"Disconnecting from current server %s:%d \n...\n", ctx.Server.Config.LHost, ctx.Server.Config.LPort)

	fmt.Printf(Info+"Connecting to %s:%d ...\n", config.LHost, config.LPort)
	send, recv, err := transport.MTLSConnect(config)
	if err != nil {
		errString := fmt.Sprintf(Errorf+"Connection to server failed: %v", err)
		return errors.New(errString)

	}
	fmt.Printf(Success+"Connected to Wiregost server at %s:%d, as user %s%s%s",
		config.LHost, config.LPort, tui.YELLOW, config.User, tui.RESET)
	fmt.Println()

	// Bind connection to server object in console
	ctx.Server = core.BindWiregostServer(send, recv)
	ctx.Server.Config = config
	go ctx.Server.ResponseMapper()

	time.Sleep(time.Millisecond * 50)

	return nil
}
