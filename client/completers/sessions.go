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

package completers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"
	"github.com/lmorg/readline"

	. "github.com/maxlandon/wiregost/client/commands"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func CompleteSessionIDs(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {
	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get last path
	splitLine := strings.Split(string(line), " ")
	last := splitLine[len(splitLine)-1]

	rpc := Context.Server.RPC

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
	}

	sessions := &clientpb.Sessions{}
	proto.Unmarshal(resp.Data, sessions)

	ghosts := map[uint32]*clientpb.Ghost{}
	for _, ghost := range sessions.Ghosts {
		ghosts[ghost.ID] = ghost
	}

	for _, g := range ghosts {
		if strings.HasPrefix(strconv.Itoa(int(g.ID)), string(last)) {
			suggestions = append(suggestions, strconv.Itoa(int(g.ID))[(len(last)):])

			var desc string
			desc += fmt.Sprintf("%s%s", tui.FOREWHITE, g.Name)
			desc += fmt.Sprintf("%s at %s%s", tui.DIM, tui.BLUE, g.RemoteAddress)
			desc += fmt.Sprintf("%s as %s%s%s@%s%s", tui.DIM, tui.FOREWHITE, g.Username, tui.BOLD, tui.RESET, g.Hostname)

			listSuggestions[strconv.Itoa(int(g.ID))[(len(last)):]] = tui.RESET + tui.DIM + desc + tui.RESET
		}
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayList
}
