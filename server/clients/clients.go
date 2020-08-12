package clients

import (
	"sync"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	"github.com/maxlandon/wiregost/server/module/stack"
)

var (
	// Consoles - All client consoles connected to the Wiregost server
	Consoles = &clients{
		Unauthenticated: &map[string]*clientpb.Client{},
		Connected:       &map[string]*clientpb.Client{},
		ClientAttempts:  &map[string]int{},
		mutex:           &sync.Mutex{},
	}
)

type clients struct {
	Unauthenticated *map[string]*clientpb.Client
	ClientAttempts  *map[string]int
	Connected       *map[string]*clientpb.Client
	mutex           *sync.Mutex
}

// GetClient - Find a client by UUID
func (c *clients) GetClient(id string) (cli *clientpb.Client) {
	return (*c.Connected)[id]
}

// AddClient - Add a client (newly connected console) to the list
func (c *clients) AddClient(cli clientpb.Client) {
	c.mutex.Lock()
	(*c.Unauthenticated)[cli.ID] = &cli
	c.mutex.Unlock()
}

func (c *clients) ConfirmClient(cli clientpb.Client) {
	c.mutex.Lock()
	(*c.Connected)[cli.ID] = &cli
	delete((*c.Unauthenticated), cli.ID)
	c.mutex.Unlock()

	// We register the client to a module stack
	stack.AssignStack(&cli)
}

// RemoveClient - Remove a client from the list
func (c *clients) RemoveClient(id string) {
	c.mutex.Lock()
	delete((*c.Connected), id)
	c.mutex.Unlock()
}

func (c *clients) IncrementClientAttempts(id string) {
	c.mutex.Lock()
	(*c.ClientAttempts)[id]++
	c.mutex.Unlock()
}

// GetUserClients - Get all clients owned by a user
func (c *clients) GetUserClients(user *dbpb.User, username string) (clis []clientpb.Client) {
	// If full user is given
	if username == "" {
		for _, cli := range *c.Connected {
			if (user.Name == cli.User.Name) && (user.ID == cli.User.ID) {
				clis = append(clis, *cli)
			}
			return
		}

	}
	// If only name is given
	for _, cli := range *c.Connected {
		if cli.User.Name == username {
			clis = append(clis, *cli)
		}
	}

	return
}
