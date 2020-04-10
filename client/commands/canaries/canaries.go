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

package canaries

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// CanariesCmd - List DNS Canaries
type CanariesCmd struct{}

var Canaries CanariesCmd

func RegisterCanaries() {
	MainParser.AddCommand(constants.Canaries, "", "", &Canaries)

	can := MainParser.Find(constants.Canaries)
	CommandMap[MAIN_CONTEXT] = append(CommandMap[MAIN_CONTEXT], can)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], can)
	can.ShortDescription = "List DNS canaries"
	can.SubcommandsOptional = true
}

// Execute - List DNS Canaries
func (g *CanariesCmd) Execute(args []string) error {
	listCanaries(Context)
	return nil
}

func listCanaries(ctx ShellContext) {
	resp := <-ctx.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgListCanaries,
	}, DefaultTimeout)
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

	table := util.NewTable()
	headers := []string{"Name", "Domain", "Triggered", "First trigger", "Latest trigger"}
	widths := []int{15, 20, 8, 20, 20}
	table.SetColumns(headers, widths)
	table.SetColWidth(40)

	for _, c := range canaries {
		triggered := ""
		if c.Triggered {
			triggered = tui.Red(tui.Bold("Yes"))
		} else {
			triggered = tui.Green("No")
		}

		table.Append([]string{c.GhostName, c.Domain, triggered, c.FirstTriggered, c.LatestTrigger})
	}

	table.Output()
}
