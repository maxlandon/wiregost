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
	"crypto/x509"
	"sync"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

var (
	Clients = &clientConns{
		Connections: &map[int]*Client{},
		mutex:       &sync.RWMutex{},
	}

	clientID = new(int)
)

// Client - Single client connection
type Client struct {
	ID          int
	User        string
	Certificate *x509.Certificate
	Send        chan *ghostpb.Envelope
	Resp        map[uint64]chan *ghostpb.Envelope
	mutex       *sync.Mutex

	// Added
	WorkspaceID uint
}

// ToProtobuf - Get the protobuf version of the Client object
func (c *Client) ToProtobuf() *clientpb.Client {
	return &clientpb.Client{
		ID:          int32(c.ID),
		User:        c.User,
		WorkspaceID: uint32(c.WorkspaceID),
	}
}

// Response - Drop an envelope into a response channel
func (c *Client) Response(envelope *ghostpb.Envelope) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if resp, ok := c.Resp[envelope.ID]; ok {
		resp <- envelope
	}
}

// clientConns - Manage client connections
type clientConns struct {
	mutex       *sync.RWMutex
	Connections *map[int]*Client
}

// AddClient - Add a client struct atomically
func (cc *clientConns) AddClient(client *Client) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	(*cc.Connections)[client.ID] = client
}

// RemoveClient - Remove a client atomically
func (cc *clientConns) RemoveClient(clientID int) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	delete((*cc.Connections), clientID)
}

// GetClientID - Get a client ID
func GetClientID() int {
	newID := (*clientID) + 1
	(*clientID)++
	return newID
}

// GetClient - Create a new client object
func GetClient(certificate *x509.Certificate) *Client {
	var user string
	if certificate != nil {
		user = certificate.Subject.CommonName
	} else {
		user = "server"
	}

	return &Client{
		ID:          GetClientID(),
		User:        user,
		Certificate: certificate,
		mutex:       &sync.Mutex{},
		Send:        make(chan *ghostpb.Envelope),
		Resp:        map[uint64]chan *ghostpb.Envelope{},
	}
}
