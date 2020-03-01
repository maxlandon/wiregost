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

package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/proto"

	consts "github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/data_service/remote"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/log"
)

var (
	handlerLog = log.ServerLogger("handlers", "ghosts")

	serverHandlers = map[uint32]interface{}{
		ghostpb.MsgRegister:    registerGhostHandler,
		ghostpb.MsgTunnelData:  tunnelDataHandler,
		ghostpb.MsgTunnelClose: tunnelCloseHandler,
	}
)

// GetGhostHandlers - Returns a map of server-side msg handlers
func GetGhostHandlers() map[uint32]interface{} {
	return serverHandlers
}

func registerGhostHandler(ghost *core.Ghost, data []byte) {
	register := &ghostpb.Register{}
	err := proto.Unmarshal(data, register)
	if err != nil {
		handlerLog.Warnf("error decoding message: %v", err)
		return
	}

	// If this is the first time we're getting reg info alert user(s)
	if ghost.Name == "" {
		defer func() {
			core.EventBroker.Publish(core.Event{
				EventType: consts.ConnectedEvent,
				Ghost:     ghost,
			})
		}()
	}

	ghost.Name = register.Name
	ghost.Hostname = register.Hostname
	ghost.Username = register.Username
	ghost.UID = register.Uid
	ghost.GID = register.Gid
	ghost.OS = register.Os
	ghost.Arch = register.Arch
	ghost.PID = register.Pid
	ghost.Filename = register.Filename
	ghost.ActiveC2 = register.ActiveC2
	ghost.Version = register.Version

	// Register workspace and host
	if register.WorkspaceID != 0 {
		ghost.WorkspaceID = uint(register.WorkspaceID)
	} else {
		ghost.WorkspaceID = uint(1) // Default workspace
	}
	id, err := registerHostToDB(ghost)
	if err != nil {
		fmt.Println(err)
	}
	if id != 0 {
		ghost.HostID = id
	}

	core.Wire.AddGhost(ghost)
}

func tunnelDataHandler(ghost *core.Ghost, data []byte) {
	tunnelData := &ghostpb.TunnelData{}
	proto.Unmarshal(data, tunnelData)
	tunnel := core.Tunnels.Tunnel(tunnelData.TunnelID)
	if tunnel != nil {
		if ghost.ID == tunnel.Ghost.ID {
			tunnel.Client.Send <- &ghostpb.Envelope{
				Type: ghostpb.MsgTunnelData,
				Data: data,
			}
		} else {
			handlerLog.Warnf("Warning: Ghost %d attempted to send data on tunnel it did not own", ghost.ID)
		}
	} else {
		handlerLog.Warnf("Data sent on nil tunnel %d", tunnelData.TunnelID)
	}
}

func tunnelCloseHandler(ghost *core.Ghost, data []byte) {
	tunnelClose := &ghostpb.TunnelClose{}
	proto.Unmarshal(data, tunnelClose)
	tunnel := core.Tunnels.Tunnel(tunnelClose.TunnelID)
	if tunnel.Ghost.ID == ghost.ID {
		handlerLog.Debugf("Ghost %d closed tunnel %d (reason: %s)", ghost.ID, tunnel.ID, tunnelClose.Err)
		core.Tunnels.CloseTunnel(tunnel.ID, tunnelClose.Err)
	} else {
		handlerLog.Warnf("Warning: Ghost %d attempted to close tunnel it did not own", ghost.ID)
	}
}

func registerHostToDB(ghost *core.Ghost) (hostID uint, err error) {

	opts := hostFilters(ghost)

	wsID := uint(1)
	if ghost.WorkspaceID != 0 {
		wsID = ghost.WorkspaceID
	}
	fmt.Println(wsID)
	rootCtx := context.Background()
	ctx := context.WithValue(rootCtx, "workspace_id", wsID)

	host, err := remote.ReportHost(ctx, opts)
	if err != nil {
		return 0, err
	} else {
		fmt.Println(host.ID)
		return host.ID, nil
	}
}

func hostFilters(ghost *core.Ghost) (opts map[string]interface{}) {
	opts = make(map[string]interface{}, 0)

	opts["os_name"] = ghost.OS
	opts["os_sp"] = ghost.Version
	opts["arch"] = ghost.Arch
	opts["addresses"] = []string{strings.Split(ghost.RemoteAddress, ":")[0]}
	opts["hostname"] = ghost.Hostname
	opts["usernames"] = ghost.Username
	opts["alive"] = true
	fmt.Println(opts)

	return opts
}
