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

	"github.com/golang/protobuf/proto"

	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
)

func rpcShell(req []byte, timeout time.Duration, resp Response) {
	shellReq := &ghostpb.ShellReq{}
	proto.Unmarshal(req, shellReq)

	ghost := core.Wire.Ghost(shellReq.GhostID)
	tunnel := core.Tunnels.Tunnel(shellReq.TunnelID)

	startShellReq, err := proto.Marshal(&ghostpb.ShellReq{
		EnablePTY: shellReq.EnablePTY,
		TunnelID:  tunnel.ID,
		Path:      shellReq.Path,
	})
	if err != nil {
		resp([]byte{}, err)
		return
	}
	rpcLog.Infof("Requesting Ghost %d to start shell", ghost.ID)
	data, err := ghost.Request(ghostpb.MsgShellReq, timeout, startShellReq)
	rpcLog.Infof("Ghost %d responded to shell start request", ghost.ID)
	resp(data, err)
}
