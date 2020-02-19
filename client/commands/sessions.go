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
	"sort"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/olekukonko/tablewriter"

	. "github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
)

func RegisterSessionCommands() {

	sessions := &Command{
		Name: "sessions",
		SubCommands: []string{
			"interact",
			"kill",
			"kill-all",
		},
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				listSessions(*r.context, r.context.Server.RPC)
			case length >= 1:
				switch r.Args[0] {
				case "interact":
					interactGhost(r.Args, *r.context, r.context.Server.RPC)
				case "kill":
				case "kill-all":
				}
			}

			return nil
		},
	}

	AddCommand("main", sessions)
	AddCommand("module", sessions)
}

func listSessions(ctx ShellContext, rpc RPCServer) {

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}
	sessions := &clientpb.Sessions{}
	proto.Unmarshal(resp.Data, sessions)

	ghosts := map[uint32]*clientpb.Ghost{}
	for _, ghost := range sessions.Ghosts {
		ghosts[ghost.ID] = ghost
	}
	if 0 < len(ghosts) {
		printGhosts(ghosts)
	} else {
		fmt.Printf(Info + "No ghosts connected\n")
	}
}

func killSession(ctx ShellContext, rpc RPCServer) {

}

func killAllSessions(ctx ShellContext, rpc RPCServer) {

}

func printGhosts(sessions map[uint32]*clientpb.Ghost) {
	table := Table()
	table.SetHeader([]string{"WsID", "ID", "Name", "Transport", "Remote Address", "Username", "Hostname", "Operating System", "Last Check-in"})
	table.SetColWidth(40)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	// Sort the keys because maps have a randomized order
	var keys []int
	for _, sliver := range sessions {
		keys = append(keys, int(sliver.ID))
	}
	sort.Ints(keys) // Fucking Go can't sort int32's, so we convert to/from int's

	for _, key := range keys {
		ghost := sessions[uint32(key)]
		workspace := ""
		if ghost.WorkspaceID != 0 {
			workspace = strconv.Itoa(int(ghost.WorkspaceID))
		}
		os := fmt.Sprintf("%s/%s", ghost.OS, ghost.Arch)

		table.Append([]string{workspace, strconv.Itoa(int(ghost.ID)), ghost.Name, ghost.Transport,
			ghost.RemoteAddress, ghost.Username, ghost.Hostname, os, ghost.LastCheckin})
	}

	table.Render()
}

func interactGhost(args []string, ctx ShellContext, rpc RPCServer) {

	name := ""
	if len(args) < 2 {
		fmt.Printf("\n" + Error + "Provide a ghost name or session number\n")
		return
	} else {
		name = args[1]
	}

	ghost := getGhost(name, rpc)
	if ghost != nil {
		ctx.CurrentAgent = parseGhost(ghost)
		// ctx.CurrentAgent = ghost
		fmt.Println(ctx.CurrentAgent)
	} else {
		fmt.Printf(Error+"Invalid ghost name or session number: %s", name)
	}
}

// Get Ghost by session ID or name
func getGhost(arg string, rpc RPCServer) *clientpb.Ghost {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, defaultTimeout)
	sessions := &clientpb.Sessions{}
	proto.Unmarshal((resp).Data, sessions)

	for _, ghost := range sessions.Ghosts {
		if strconv.Itoa(int(ghost.ID)) == arg || ghost.Name == arg {
			return ghost
		}
	}
	return nil
}

func GetGhosts(rpc RPCServer) *clientpb.Sessions {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, defaultTimeout)
	sessions := &clientpb.Sessions{}
	proto.Unmarshal((resp).Data, sessions)

	return sessions
}

func parseGhost(ghost *clientpb.Ghost) *core.Ghost {

	g := &core.Ghost{}
	g.ID = ghost.ID
	g.Name = ghost.Name
	g.Hostname = ghost.Hostname
	g.Username = ghost.Username
	g.UID = ghost.UID
	g.GID = ghost.GID
	g.OS = ghost.OS
	g.Version = ghost.Version
	g.Arch = ghost.Arch
	g.Transport = ghost.Transport
	g.RemoteAddress = ghost.RemoteAddress
	g.PID = ghost.PID
	g.Filename = ghost.Filename

	layout := "2006-01-02T15:04:05"
	last, _ := time.Parse(layout, ghost.LastCheckin)
	g.LastCheckin = &last
	// g.Send =
	// g.Resp
	// g.ActiveC2

	// Added
	// WorkspaceID
	// Host        *models.Host
	return g
}
