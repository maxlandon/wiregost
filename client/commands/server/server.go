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

package server

import (
	"fmt"
	"strconv"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/assets"
	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/util"
)

// ServerCmd - List available Wiregost servers
type ServerCmd struct {
}

var Server ServerCmd

func RegisterServer() {
	MainParser.AddCommand(constants.Server, "", "", &Server)

	s := MainParser.Find(constants.Server)
	CommandMap[MAIN_CONTEXT] = append(CommandMap[MAIN_CONTEXT], s)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], s)
	s.ShortDescription = "List or connect to one of the available Wiregost servers"
	s.SubcommandsOptional = true
}

// Execute - Command
func (r *ServerCmd) Execute(args []string) error {

	configs := assets.GetConfigs()
	if len(configs) == 0 {
		fmt.Printf(Warnf+"No config files found at %s or -config\n", assets.GetConfigDir())
		return nil
	}

	table := util.NewTable()
	headers := []string{"User Name", "LHost", "LPort", "Default", "Connected"}
	widths := []int{15, 15, 5, 7, 10}
	table.SetColumns(headers, widths)
	table.SetColWidth(40)

	for _, c := range configs {
		def := ""
		if c.IsDefault {
			def = "default"
		}
		port := strconv.Itoa(c.LPort)

		connected := ""
		if c.LHost == Context.Server.Config.LHost && c.LPort == Context.Server.Config.LPort {
			connected = fmt.Sprintf("%sconnected%s", tui.GREEN, tui.RESET)
		}

		table.Append([]string{c.User, c.LHost, port, def, connected})
	}
	table.Output()

	return nil
}
