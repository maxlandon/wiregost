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

package sessions

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// SessionsCmd - List connected ghost implants
type SessionsCmd struct{}

var Sessions SessionsCmd

func RegisterSessions() {
	CommandParser.AddCommand(constants.Sessions, "", "", &Sessions)

	ses := CommandParser.Find(constants.Sessions)
	CommandMap[MAIN_CONTEXT] = append(CommandMap[MAIN_CONTEXT], ses)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], ses)
	ses.ShortDescription = "List/interact/kill currently connected ghost implants"
	ses.SubcommandsOptional = true
}

// Execute - List connected ghost implants
func (s *SessionsCmd) Execute(args []string) error {
	listSessions()
	return nil
}

func listSessions() {

	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, DefaultTimeout)
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
		printGhosts(ghosts, Context.Server.RPC)
	} else {
		fmt.Printf(Info + "No ghosts connected\n")
	}
}

func printGhosts(sessions map[uint32]*clientpb.Ghost, rpc RPCServer) {

	table := util.NewTable()
	headers := []string{"WsID", "ID", "Name", "Proto", "Remote Address", "user@host", "Platform", "Status"}
	widths := []int{4, 2, 15, 5, 15, 30, 15, 8}
	table.SetColumns(headers, widths)
	table.SetColWidth(40)

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
		userHost := fmt.Sprintf("%s@%s", ghost.Username, ghost.Hostname)

		status := getSessionStatus(ghost, rpc)
		table.Append([]string{workspace, strconv.Itoa(int(ghost.ID)), ghost.Name, ghost.Transport,
			ghost.RemoteAddress, userHost, os, status})
	}

	table.Render()
}

func getSessionStatus(ghost *clientpb.Ghost, rpc RPCServer) string {

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgListGhostBuilds,
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return ""
	}

	config := &clientpb.GhostConfig{}

	builds := &clientpb.GhostBuilds{}
	proto.Unmarshal(resp.Data, builds)
	for _, b := range builds.Configs {
		if ghost.Name == b.Name {
			config = b
		}
	}

	dur, errDur := time.ParseDuration(fmt.Sprintf("%ds", config.ReconnectInterval))
	if errDur != nil {
		fmt.Println(errDur)
	}

	lastCheckin, err := time.Parse(time.RFC1123, ghost.LastCheckin)
	if err != nil {
		fmt.Println(err)
	}

	if lastCheckin.Add(dur).After(time.Now()) {
		return tui.Green("Alive")
	} else if lastCheckin.Add(dur * time.Duration(config.MaxConnectionErrors+1)).After(time.Now()) {
		return tui.Yellow("Delayed")
	} else {
		return tui.Red("Dead")
	}
}
