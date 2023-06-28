package core

// Wiregost - Post-Exploitation & Implant Framework
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

import (
	"sync"

	consts "github.com/maxlandon/wiregost/internal/client/constants"
	"github.com/maxlandon/wiregost/internal/proto/clientpb"
)

var (
	// Clients - Manages client active
	Clients = &clients{
		active: map[int]*Client{},
		mutex:  &sync.Mutex{},
	}

	clientID = 0
)

// Client - Single client connection
type Client struct {
	ID       int
	Operator *clientpb.Operator
}

// ToProtobuf - Get the protobuf version of the object
func (c *Client) ToProtobuf() *clientpb.Client {
	return &clientpb.Client{
		ID:       uint32(c.ID),
		Operator: c.Operator,
	}
}

// clients - Manage active clients
type clients struct {
	mutex  *sync.Mutex
	active map[int]*Client
}

// AddClient - Add a client struct atomically
func (cc *clients) Add(client *Client) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	cc.active[client.ID] = client
	EventBroker.Publish(Event{
		EventType: consts.JoinedEvent,
		Client:    client,
	})
}

// AddClient - Add a client struct atomically
func (cc *clients) ActiveOperators() []string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	operators := []string{}
	for _, client := range cc.active {
		operators = append(operators, client.Operator.Name)
	}
	return operators
}

// RemoveClient - Remove a client struct atomically
func (cc *clients) Remove(clientID int) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	client := cc.active[clientID]
	delete(cc.active, clientID)
	EventBroker.Publish(Event{
		EventType: consts.LeftEvent,
		Client:    client,
	})
}

// nextClientID - Get a client ID
func nextClientID() int {
	newID := clientID + 1
	clientID++
	return newID
}

// NewClient - Create a new client object
func NewClient(operatorName string) *Client {
	return &Client{
		ID: nextClientID(),
		Operator: &clientpb.Operator{
			Name: operatorName,
		},
		// mutex: &sync.RWMutex{},
	}
}
