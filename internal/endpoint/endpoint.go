package endpoint

import (
	"fmt"
	"net"

	"github.com/maxlandon/wiregost/internal/dispatch"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/user"
	"github.com/maxlandon/wiregost/internal/workspace"
)

type Endpoint struct { // PROPOSED CHANGES
	clients  []*Client
	connect  chan net.Conn
	requests chan messages.ClientRequest
}

func NewEndpoint() *Endpoint {
	e := &Endpoint{
		clients:  make([]*Client, 0),
		connect:  make(chan net.Conn),
		requests: make(chan messages.ClientRequest),
	}

	go e.Listen()
	go e.ForwardResponses()

	return e
}

// Listen listens for connections and messages to broadcast
func (e *Endpoint) Listen() {
	for {
		select {
		case conn := <-e.connect:
			e.Join(conn)
		case msg := <-e.requests:
			fmt.Println("Received request")
			user.AuthReqs <- msg
			fmt.Println(msg.UserId)
			auth := <-user.AuthResp
			switch auth.UserId {
			case 0:
				fmt.Println("Here should be sent a ConnRefused message to the client")
			default:
				fmt.Println("dispatching")
				dispatch.DispatchRequest(auth)
			}
		}
	}
}

// Connect passing connection to the server
func (e *Endpoint) Connect(conn net.Conn) {
	e.connect <- conn
}

// Join creates new client and starts listening for client messages
func (e *Endpoint) Join(conn net.Conn) {
	client := CreateClient(conn)
	e.clients = append(e.clients, client)
	fmt.Println(e.clients)
	for _, client := range e.clients {
		fmt.Println(client.status)
	}
	go func() {
		for {
			e.requests <- <-client.requests
		}
	}()
}

// Remove disconnected client from list
func (e *Endpoint) Remove(i int) {
	e.clients = append(e.clients[:i], e.clients[i+1:]...)
}

func (e *Endpoint) ForwardResponses() {
	for {
		// Remove disconnected clients
		for i, client := range e.clients {
			if client.status == 0 {
				e.Remove(i)
			}
		}
		select {
		case res := <-dispatch.Responses:
			fmt.Println("Handled response from dispatch")
			for _, client := range e.clients {
				if client.id == res.ClientId {
					client.responses <- res
				}
			}
		case res := <-workspace.Responses:
			fmt.Println("handled response from workspace")
			for _, client := range e.clients {
				if client.id == res.ClientId {
					client.responses <- res
				}
			}
		}
	}
}
