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

package console

import (
	"bufio"
	"fmt"
	"os"

	"github.com/evilsocket/islazy/tui"

	consts "github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/core"
)

func (c *Console) eventLoop(server *core.WiregostServer) {
	stdout := bufio.NewWriter(os.Stdout)
	for event := range server.Events {

		switch event.EventType {

		case consts.CanaryEvent:
			fmt.Printf("%s[WARNING]%s %s has been burned (DNS Canary) \n", tui.YELLOW, tui.RESET, event.Ghost.Name)
			// sessions := cmd.GhostSessionsByName(event.Ghost.Name, server.RPC)
			// for _, ghost := range sessions {
			//         fmt.Printf("%s[!]%s \tSession #%d is compromised\n", tui.YELLOW, tui.RESET, ghost.ID)
			// }
			fmt.Println()

		case consts.ServerErrorStr:
			fmt.Printf("%s[!]%s Server connection error! \n\n", tui.RED, tui.RESET)
			os.Exit(4)

		case consts.JoinedEvent:
			fmt.Println()
			fmt.Printf("%s*%s %s connected to the server \n", tui.BLUE, tui.RESET, event.Client.User)
			c.hardRefresh()
		case consts.LeftEvent:
			fmt.Println()
			fmt.Printf("%s*%s %s disconnected from the server \n", tui.BLUE, tui.RESET, event.Client.User)
			c.hardRefresh()

		case consts.StoppedEvent:
			job := event.Job
			fmt.Printf("%s[*]%s Job #%d stopped (%s/%s) \n", tui.BLUE, tui.RESET, job.ID, job.Protocol, job.Name)

		case consts.ConnectedEvent:
			ghost := event.Ghost
			fmt.Printf("%s[*]%s Session #%d %s - %s (%s) - %s/%s \n\n",
				tui.GREEN, tui.RESET, ghost.ID, ghost.Name, ghost.RemoteAddress, ghost.Hostname, ghost.OS, ghost.Arch)

		case consts.DisconnectedEvent:
			ghost := event.Ghost
			fmt.Printf("%s[!]%s Lost session #%d %s - %s (%s) - %s/%s\n",
				tui.RED, tui.RESET, ghost.ID, ghost.Name, ghost.RemoteAddress, ghost.Hostname, ghost.OS, ghost.Arch)
			activeGhost := c.currentAgent
			if activeGhost != nil && ghost.ID == activeGhost.ID {
				c.currentAgent = nil
				// app.SetPrompt(getPrompt())
				fmt.Printf("%s[!]%s Active sliver diconnected\n", tui.RED, tui.RESET)
			}
			fmt.Println()

		}

		// fmt.Printf(getPrompt())
		stdout.Flush()
	}
}
