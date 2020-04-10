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

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// SessionsInteractCmd - Interact with a ghost implant session
type SessionsInteractCmd struct {
	Positional struct {
		SessionID int `description:"Session ID" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var SessionsInteract SessionsInteractCmd

func RegisterSessionsInteract() {
	ses := MainParser.Find(constants.Sessions)
	ses.AddCommand(constants.SessionsInteract, "", "", &SessionsInteract)

	sesInt := ses.Find(constants.SessionsInteract)
	sesInt.ShortDescription = "Interact with a ghost implant session"
}

// Execute - Interact with a ghost implant session
func (s *SessionsInteractCmd) Execute(args []string) error {
	interactGhost(s.Positional.SessionID, Context)
	return nil
}

func interactGhost(id int, ctx ShellContext) {

	ghost := &clientpb.Ghost{}

	resp := <-ctx.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(Error+"Impossible to establish communication with session %d\n", id)
		return
	}

	sessions := &clientpb.Sessions{}
	proto.Unmarshal((resp).Data, sessions)

	for _, g := range sessions.Ghosts {
		if int(g.ID) == id {
			ghost = g
		}
	}

	if ghost != nil {
		// Get cwd, and check that session is connected by the same way
		data, _ := proto.Marshal(&ghostpb.PwdReq{
			GhostID: ghost.ID,
		})
		resp := <-Context.Server.RPC(&ghostpb.Envelope{
			Type: ghostpb.MsgPwdReq,
			Data: data,
		}, DefaultTimeout)
		if resp.Err != "" {
			fmt.Printf("\n"+Error+"Impossible to establish communication with session %d\n", id)
			return
		}

		pwd := &ghostpb.Pwd{}
		err := proto.Unmarshal(resp.Data, pwd)
		if err != nil {
			fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
			return
		}
		*ctx.Ghost = *ghost
		*ctx.GhostPwd = pwd.Path

	} else {
		fmt.Printf(Error+"Invalid ghost name or session number: %d\n", id)
	}
}

// Get Ghost by session ID or name
func getGhost(id int, rpc RPCServer) *clientpb.Ghost {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, DefaultTimeout)
	sessions := &clientpb.Sessions{}
	proto.Unmarshal((resp).Data, sessions)

	for _, ghost := range sessions.Ghosts {
		if int(ghost.ID) == id {
			return ghost
		}
	}
	return nil
}

// // GhostSessionsByName - Get a session by name
// func GhostSessionsByName(name string, rpc RPCServer) []*clientpb.Ghost {
//         resp := <-rpc(&ghostpb.Envelope{
//                 Type: clientpb.MsgSessions,
//                 Data: []byte{},
//         }, defaultTimeout)
//         allSessions := &clientpb.Sessions{}
//         proto.Unmarshal((resp).Data, allSessions)
//
//         sessions := []*clientpb.Ghost{}
//         for _, ghost := range allSessions.Ghosts {
//                 if ghost.Name == name {
//                         sessions = append(sessions, ghost)
//                 }
//         }
//         return sessions
// }
//
