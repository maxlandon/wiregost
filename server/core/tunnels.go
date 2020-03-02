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

package core

import (
	"crypto/rand"
	"encoding/binary"
	"sync"

	"github.com/gogo/protobuf/proto"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

var (
	// Tunnels - ALl programmatic channels between clients and implants
	Tunnels = tunnels{
		tunnels: &map[uint64]*tunnel{},
		mutex:   &sync.RWMutex{},
	}
)

// A tunnel is essentially a mapping between a specific client and server
// with an identifier: the server doesn't really care what data gets passed
// back and forth, it just facilitates the connection.
type tunnel struct {
	ID     uint64
	Ghost  *Ghost
	Client *Client
}

type tunnels struct {
	tunnels *map[uint64]*tunnel
	mutex   *sync.RWMutex
}

// CreateTunnel - Creates a tunnel between a Client and a Ghost objects
func (t *tunnels) CreateTunnel(client *Client, ghostID uint32) *tunnel {
	tunID := newTunnelID()
	ghost := Wire.Ghost(ghostID)

	tun := &tunnel{
		ID:     tunID,
		Ghost:  ghost,
		Client: client,
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()
	(*t.tunnels)[tun.ID] = tun

	return tun
}

// CloseTunnel - Terminates a tunnel set up between a Client and a Ghost ojects
func (t *tunnels) CloseTunnel(tunnelID uint64, reason string) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	tunnel := (*t.tunnels)[tunnelID]
	if tunnel != nil {
		tunnelClose, _ := proto.Marshal(&ghostpb.TunnelClose{
			TunnelID: tunnelID,
			Err:      reason,
		})

		tunnel.Client.Send <- &ghostpb.Envelope{
			Type: ghostpb.MsgTunnelClose,
			Data: tunnelClose,
		}

		tunnel.Ghost.Send <- &ghostpb.Envelope{
			Type: ghostpb.MsgTunnelClose,
			Data: tunnelClose,
		}

		delete(*t.tunnels, tunnelID)
		return true
	}
	return false
}

// Tunnel - Returns a Tunnel by ID
func (t *tunnels) Tunnel(tunnelID uint64) *tunnel {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return (*t.tunnels)[tunnelID]
}

func newTunnelID() uint64 {
	randBuf := make([]byte, 8)
	rand.Read(randBuf)
	return binary.LittleEndian.Uint64(randBuf)
}
