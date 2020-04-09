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
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// SessionsKillCmd - Kill a ghost implant session
type SessionsKillCmd struct {
	Positional struct {
		SessionID int `description:"Session ID" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var SessionsKill SessionsKillCmd

func RegisterSessionsKill() {
	ses := CommandParser.Find(constants.Sessions)
	ses.AddCommand(constants.SessionsKill, "", "", &SessionsKill)

	kill := ses.Find(constants.SessionsKill)
	kill.ShortDescription = "Kill a ghost implant session"
}

// Execute - Kill a ghost implant session
func (s *SessionsKillCmd) Execute(args []string) error {

	ghost := getGhost(s.Positional.SessionID, Context.Server.RPC)
	if ghost != nil {
		data, _ := proto.Marshal(&ghostpb.KillReq{
			GhostID: ghost.ID,
			Force:   true,
		})
		Context.Server.RPC(&ghostpb.Envelope{
			Type: ghostpb.MsgKill,
			Data: data,
		}, 5*time.Second)

		fmt.Printf(Info+"Killed agent %s (Session %d)\n", ghost.Name, ghost.ID)
	} else {
		fmt.Printf(Error+"Invalid ghost session ID: %d", s.Positional.SessionID)
	}

	return nil
}
