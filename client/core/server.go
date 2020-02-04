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
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/maxlandon/wiregost/client/assets"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// WiregostServer - Server info
type WiregostServer struct {
	Send      chan *ghostpb.Envelope
	recv      chan *ghostpb.Envelope
	responses *map[uint64]chan *ghostpb.Envelope
	mutex     *sync.RWMutex
	Config    *assets.ClientConfig
	Events    chan *clientpb.Event
	Tunnels   *tunnels
}

// CreateTunnel - Create a new tunnel on the server, returns tunnel metadata
func (ws *WiregostServer) CreateTunnel(ghostID uint32, defaultTimeout time.Duration) (*tunnel, error) {
	tunReq := &clientpb.TunnelCreateReq{GhostID: ghostID}
	tunReqData, _ := proto.Marshal(tunReq)

	tunResp := <-ws.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgTunnelCreate,
		Data: tunReqData,
	}, defaultTimeout)
	if tunResp.Err != "" {
		return nil, fmt.Errorf("Error: %s", tunResp.Err)
	}

	tunnelCreated := &clientpb.TunnelCreate{}
	proto.Unmarshal(tunResp.Data, tunnelCreated)

	tunnel := ws.Tunnels.bindTunnel(tunnelCreated.GhostID, tunnelCreated.TunnelID)

	log.Printf("Created new tunnel with ID %d", tunnel.ID)

	return tunnel, nil
}

// ResponseMapper - Maps recv'd envelopes to response channels
func (ws *WiregostServer) ResponseMapper() {
	for envelope := range ws.recv {
		if envelope.ID != 0 {
			ws.mutex.Lock()
			if resp, ok := (*ws.responses)[envelope.ID]; ok {
				resp <- envelope
			}
			ws.mutex.Unlock()
		} else {
			// If the mewsage does not have an envelope ID then we route it based on type
			switch envelope.Type {

			case clientpb.MsgEvent:
				event := &clientpb.Event{}
				err := proto.Unmarshal(envelope.Data, event)
				if err != nil {
					log.Printf("Failed to decode event envelope")
					continue
				}
				// log.Printf("[client] Routing event mewsage")
				ws.Events <- event

			case ghostpb.MsgTunnelData:
				tunnelData := &ghostpb.TunnelData{}
				err := proto.Unmarshal(envelope.Data, tunnelData)
				if err != nil {
					log.Printf("Failed to decode tunnel data envelope")
					continue
				}
				// log.Printf("[client] Routing tunnel data with id %d", tunnelData.TunnelID)
				ws.Tunnels.RecvTunnelData(tunnelData)

			case ghostpb.MsgTunnelClose:
				tunnelClose := &ghostpb.TunnelClose{}
				err := proto.Unmarshal(envelope.Data, tunnelClose)
				if err != nil {
					log.Printf("Failed to decode tunnel data envelope")
					continue
				}
				ws.Tunnels.Close(tunnelClose.TunnelID)

			}
		}
	}
}

// RPC - Send a request envelope and wait for a response (blocking)
func (ws *WiregostServer) RPC(envelope *ghostpb.Envelope, timeout time.Duration) chan *ghostpb.Envelope {
	reqID := EnvelopeID()
	envelope.ID = reqID
	envelope.Timeout = timeout.Nanoseconds()
	resp := make(chan *ghostpb.Envelope)
	ws.AddRespListener(reqID, resp)
	ws.Send <- envelope
	respCh := make(chan *ghostpb.Envelope)
	go func() {
		defer ws.RemoveRespListener(reqID)
		select {
		case respEnvelope := <-resp:
			respCh <- respEnvelope
		case <-time.After(timeout + time.Second):
			respCh <- &ghostpb.Envelope{Err: "Timeout"}
		}
	}()
	return respCh
}

// AddRespListener - Add a response listener
func (ws *WiregostServer) AddRespListener(envelopeID uint64, resp chan *ghostpb.Envelope) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	(*ws.responses)[envelopeID] = resp
}

// RemoveRespListener - Remove a listener
func (ws *WiregostServer) RemoveRespListener(envelopeID uint64) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	close((*ws.responses)[envelopeID])
	delete((*ws.responses), envelopeID)
}

// BindWiregostServer - Bind send/recv channels to a server
func BindWiregostServer(send, recv chan *ghostpb.Envelope) *WiregostServer {
	server := &WiregostServer{
		Send:      send,
		recv:      recv,
		responses: &map[uint64]chan *ghostpb.Envelope{},
		mutex:     &sync.RWMutex{},
		Events:    make(chan *clientpb.Event, 1),
	}
	server.Tunnels = &tunnels{
		server:  server,
		tunnels: &map[uint64]*tunnel{},
		mutex:   &sync.RWMutex{},
	}
	return server
}

// EnvelopeID - Generate random ID
func EnvelopeID() uint64 {
	randBuf := make([]byte, 8) // 64 bits of randomnews
	rand.Read(randBuf)
	return binary.LittleEndian.Uint64(randBuf)
}
