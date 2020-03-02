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

const (
	tunDefaultTimeout = 30 * time.Second
)

func tunnelCreate(client *core.Client, req []byte, resp Response) {
	tunCreateReq := &clientpb.TunnelCreateReq{}
	proto.Unmarshal(req, tunCreateReq)

	tunnel := core.Tunnels.CreateTunnel(client, tunCreateReq.GhostID)

	data, err := proto.Marshal(&clientpb.TunnelCreate{
		GhostID:  tunnel.Ghost.ID,
		TunnelID: tunnel.ID,
	})

	resp(data, err)
}

func tunnelData(client *core.Client, req []byte, _ Response) {
	tunnelData := &ghostpb.TunnelData{}
	proto.Unmarshal(req, tunnelData)
	tunnel := core.Tunnels.Tunnel(tunnelData.TunnelID)
	if tunnel != nil && client.ID == tunnel.Client.ID {
		tunnel.Ghost.Request(ghostpb.MsgTunnelData, tunDefaultTimeout, req)
	} else {
		rpcLog.Warnf("Data sent on nil tunnel %d", tunnelData.TunnelID)
	}
}

func tunnelClose(client *core.Client, req []byte, resp Response) {
	tunCloseReq := &clientpb.TunnelCloseReq{}
	proto.Unmarshal(req, tunCloseReq)

	tunnel := core.Tunnels.Tunnel(tunCloseReq.TunnelID)

	if tunnel != nil && client.ID == tunnel.Client.ID {
		closed := core.Tunnels.CloseTunnel(tunCloseReq.TunnelID, "Client exit")
		closeResp := &ghostpb.TunnelClose{
			TunnelID: tunCloseReq.TunnelID,
		}
		if !closed {
			closeResp.Err = "Failed to close tunnel"
		}
		data, err := proto.Marshal(closeResp)
		resp(data, err)
	}
}
