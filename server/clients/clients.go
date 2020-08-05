package clients

import (
	"net"
	"sync"

	"google.golang.org/grpc"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
)

var (
	// Consoles - All client consoles connected to the Wiregost server
	Consoles = &consoles{
		Unauthenticated: &map[string]*clientpb.Client{},
		Connected:       &map[string]*clientpb.Client{},
		EventBrokers:    &map[string]*eventBroker{},
		ClientAttempts:  &map[string]int{},
		mutex:           &sync.Mutex{},
	}
)

type consoles struct {
	Unauthenticated *map[string]*clientpb.Client
	ClientAttempts  *map[string]int
	Connected       *map[string]*clientpb.Client
	EventBrokers    *map[string]*eventBroker
	mutex           *sync.Mutex
}

// GetClient - Find a client by UUID
func (c *consoles) GetClient(id string) (cli *clientpb.Client) {
	return (*c.Connected)[id]
}

// AddClient - Add a client (newly connected console) to the list
func (c *consoles) AddClient(cli clientpb.Client) {
	c.mutex.Lock()
	(*c.Unauthenticated)[cli.Token] = &cli
	c.mutex.Unlock()
}

func (c *consoles) ConfirmClient(cli clientpb.Client) {
	c.mutex.Lock()
	// Add client object
	(*c.Connected)[cli.Token] = &cli
	delete((*c.Unauthenticated), cli.Token)
	// Bind event broker
	(*c.EventBrokers)[cli.Token] = &eventBroker{}
	c.mutex.Unlock()
}

// RemoveClient - Remove a client from the list
func (c *consoles) RemoveClient(id string) {
	c.mutex.Lock()
	delete((*c.Connected), id)
	delete((*c.EventBrokers), id)
	c.mutex.Unlock()
}

func (c *consoles) IncrementClientAttempts(id string) {
	c.mutex.Lock()
	(*c.ClientAttempts)[id]++
	c.mutex.Unlock()
}

// GetUserClients - Get all clients owned by a user
func (c *consoles) GetUserClients(user *dbpb.User, username string) (clis []clientpb.Client) {
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

// Serve - Listen for incoming console connections
func Serve() (server *grpc.Server, ln net.Listener, err error) {

	// Logging

	// Checking for at least one user
	err = CreateDefaultUser()

	// Setup gRPC server
	server = SetupGRPC()

	// Register RPC Services
	RegisterServices(server)

	// Listen and serve gRPC
	ServeGRPC(server)

	return
}
