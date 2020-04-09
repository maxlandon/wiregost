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
	"time"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// SessionsKillAllCmd - Kill all ghost implant sessions
type SessionsKillAllCmd struct {
}

var SessionsKillAll SessionsKillAllCmd

func RegisterSessionsKillAll() {
	ses := CommandParser.Find(constants.Sessions)
	ses.AddCommand(constants.SessionsKillAll, "", "", &SessionsKillAll)

	kill := ses.Find(constants.SessionsKillAll)
	kill.ShortDescription = "Kill all ghost implant sessions"
}

// Execute - Kill all ghost implant sessions
func (s *SessionsKillAllCmd) Execute(args []string) error {

	sessions := GetGhosts(Context.Server.RPC)
	for _, session := range sessions.Ghosts {
		data, _ := proto.Marshal(&ghostpb.KillReq{
			GhostID: session.ID,
			Force:   true,
		})
		Context.Server.RPC(&ghostpb.Envelope{
			Type: ghostpb.MsgKill,
			Data: data,
		}, 5*time.Second)

		fmt.Printf(Info+"Killed %s (%d)\n", Context.Ghost.Name, Context.Ghost.ID)
	}

	return nil
}

// GetGhosts - Get all connected sessions
func GetGhosts(rpc RPCServer) *clientpb.Sessions {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, DefaultTimeout)
	sessions := &clientpb.Sessions{}
	proto.Unmarshal((resp).Data, sessions)

	return sessions
}
