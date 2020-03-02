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

package rpc

import (
	"time"

	"github.com/gogo/protobuf/proto"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
)

func rpcSessions(_ []byte, timeout time.Duration, resp Response) {

	sessions := &clientpb.Sessions{}
	if 0 < len(*core.Wire.Ghosts) {
		for _, ghost := range *core.Wire.Ghosts {
			sessions.Ghosts = append(sessions.Ghosts, ghost.ToProtobuf())
		}
	}
	data, err := proto.Marshal(sessions)
	if err != nil {
		rpcLog.Errorf("Error encoding rpc response %v", err)
	}
	resp(data, err)
}

func rpcKill(data []byte, timeout time.Duration, resp Response) {

	killReq := &ghostpb.KillReq{}
	err := proto.Unmarshal(data, killReq)
	if err != nil {
		resp([]byte{}, err)
	}
	ghost := core.Wire.Ghost(killReq.GhostID)
	core.Wire.RemoveGhost(ghost)

	data, err = ghost.Request(ghostpb.MsgKill, timeout, data)
	resp(data, err)
}
