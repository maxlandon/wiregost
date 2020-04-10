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

package ghosts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// GhostsCmd - List ghost implant builds
type GhostsCmd struct{}

var Ghosts GhostsCmd

func RegisterGhosts() {
	MainParser.AddCommand(constants.Ghosts, "", "", &Ghosts)

	ghosts := MainParser.Find(constants.Ghosts)
	CommandMap[MAIN_CONTEXT] = append(CommandMap[MAIN_CONTEXT], ghosts)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], ghosts)
	ghosts.ShortDescription = "List previously compiled ghost implant builds"
	ghosts.SubcommandsOptional = true
}

// Execute - List ghost implant builds
func (g *GhostsCmd) Execute(args []string) error {
	listGhostBuilds(Context)
	return nil
}

func listGhostBuilds(ctx ShellContext) {
	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgListGhostBuilds,
	}, DefaultTimeout)
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

	table := util.NewTable()
	headers := []string{"WsID", "Name", "Platform", "Format", "Command & Control", "Limits", "Debug"}
	widths := []int{4, 15, 15, 7, 20, 20, 5}
	table.SetColumns(headers, widths)
	table.SetColWidth(40)

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

func getLimitsString(config *clientpb.GhostConfig) string {
	limits := []string{}
	if config.LimitDatetime != "" {
		limits = append(limits, fmt.Sprintf("datetime=%s", config.LimitDatetime))
	}
	if config.LimitDomainJoined {
		limits = append(limits, fmt.Sprintf("domainjoined=%v", config.LimitDomainJoined))
	}
	if config.LimitUsername != "" {
		limits = append(limits, fmt.Sprintf("username=%s", config.LimitUsername))
	}
	if config.LimitHostname != "" {
		limits = append(limits, fmt.Sprintf("hostname=%s", config.LimitHostname))
	}
	return strings.Join(limits, "; ")
}
