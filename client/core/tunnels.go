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
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

const (
	randomIDSize = 16 // 64bits
)

// tunnel - Duplex data tunnel
type tunnel struct {
	server  *WiregostServer
	GhostID uint32
	ID      uint64
	Recv    chan []byte
	isOpen  bool
}

type tunnels struct {
	server  *WiregostServer
	tunnels *map[uint64]*tunnel
	mutex   *sync.RWMutex
}

func (t *tunnels) bindTunnel(GhostID uint32, TunnelID uint64) *tunnel {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	(*t.tunnels)[TunnelID] = &tunnel{
		server:  t.server,
		GhostID: GhostID,
		ID:      TunnelID,
		Recv:    make(chan []byte),
		isOpen:  true,
	}

	return (*t.tunnels)[TunnelID]
}

// RecvTunnelData - Routes a TunnelData protobuf msg to the correct tunnel object
func (t *tunnels) RecvTunnelData(tunnelData *ghostpb.TunnelData) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	tunnel := (*t.tunnels)[tunnelData.TunnelID]
	if tunnel != nil {
		(*tunnel).Recv <- tunnelData.Data
	} else {
		log.Printf("No client tunnel with ID %d", tunnelData.TunnelID)
	}
}

func (t *tunnels) Close(ID uint64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	(*t.tunnels)[ID].isOpen = false
	close((*t.tunnels)[ID].Recv)
	delete(*t.tunnels, ID)
}

type tunnelAddr struct {
	network string
	addr    string
}

func (a *tunnelAddr) Network() string {
	return a.network
}

func (a *tunnelAddr) String() string {
	return fmt.Sprintf("%s://%s", a.network, a.addr)
}

func (t *tunnel) Write(data []byte) (n int, err error) {
	log.Printf("Sending %d bytes on tunnel %d (ghost %d)", len(data), t.ID, t.GhostID)
	if !t.isOpen {
		return 0, io.EOF
	}
	tunnelData := &ghostpb.TunnelData{
		GhostID:  t.GhostID,
		TunnelID: t.ID,
		Data:     data,
	}
	rawTunnelData, err := proto.Marshal(tunnelData)
	t.server.Send <- &ghostpb.Envelope{
		Type: ghostpb.MsgTunnelData,
		Data: rawTunnelData,
	}
	n = len(data)
	return
}

func (t *tunnel) Read(data []byte) (n int, err error) {
	var buff bytes.Buffer
	if !t.isOpen {
		return 0, io.EOF
	}
	select {
	case msg := <-t.Recv:
		buff.Write(msg)
	default:
		break
	}
	n = copy(data, buff.Bytes())
	return
}

func (t *tunnel) Close() error {
	tunnelClose, err := proto.Marshal(&ghostpb.ShellReq{
		TunnelID: t.ID,
	})
	t.server.RPC(&ghostpb.Envelope{
		Type: ghostpb.MsgTunnelClose,
		Data: tunnelClose,
	}, 30*time.Second)
	close(t.Recv)
	return err
}
