package endpoint

import (
	"fmt"
	"net"
	"strings"

	"github.com/maxlandon/wiregost/internal/logging"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/user"
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

func (e *Endpoint) AuthenticateConnection(msg messages.ClientRequest, id int) {
	switch id {
	case 0:
		fmt.Println("rejected")
		status := "rejected"
		connected := messages.EndpointResponse{
			Connected: false,
			Status:    status,
		}
		res := messages.Message{
			ClientId: msg.ClientId,
			Type:     "connection",
			Content:  connected,
		}
		for _, client := range e.clients {
			if client.id == res.ClientId {
				client.responses <- res
				client.disconnect <- true
			}
		}
	default:
		status := "authenticated"
		connected := messages.EndpointResponse{
			Connected: true,
			Status:    status,
		}
		res := messages.Message{
			ClientId: msg.ClientId,
			Type:     "connection",
			Content:  connected,
		}
		for _, client := range e.clients {
			if client.id == res.ClientId {
				// Send response back...
				client.responses <- res
				// And fill client information from message TEMPORARY WE NEED TO REWRITE THIS
				client.UserID = id
			}
		}

	}
}

// Listen listens for connections and messages to broadcast
func (e *Endpoint) Listen() {
	for {
		select {
		case conn := <-e.connect:
			e.Join(conn)
		case msg := <-e.requests:
			user.AuthReqs <- msg
			auth := <-user.AuthResp
			// If client is opening connection, send him confirmation
			if auth.Command[0] == "connect" {
				e.AuthenticateConnection(msg, auth.UserId)
				// Send current workspace of last shell to new shell
				if len(e.clients) > 1 {
					for i := 1; i < len(e.clients); i++ {
						if e.clients[i].UserID == auth.UserId {
							res := messages.Notification{
								Type:        "workspace",
								Action:      "switch",
								WorkspaceId: e.clients[i].CurrentWorkspaceId,
								Workspace:   e.clients[i].CurrentWorkspace,
							}
							msg := messages.Message{
								ClientId: auth.ClientId,
								Type:     "notification",
								Content:  res,
							}
							for _, client := range e.clients {
								if client.id == auth.ClientId {
									client.responses <- msg
								}
							}

						}
					}
				}
			}
			if strings.Join(auth.Command[:2], " ") == "log level" {
				for _, client := range e.clients {
					if client.id == auth.ClientId {
						client.Logger.SetLevel(auth)
					}
				}
			} else {
				// Else, authenticate anyway but forward requests to messageser
				switch auth.UserId {
				case 0:
					connected := messages.EndpointResponse{
						Connected: false,
						Status:    "rejected",
					}
					res := messages.Message{
						ClientId: msg.ClientId,
						Type:     "connection",
						Content:  connected,
					}
					for _, client := range e.clients {
						if client.id == res.ClientId {
							client.responses <- res
							client.disconnect <- true
						}
					}
				default:
					e.DispatchRequest(auth)
				}
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
		case res := <-messages.Responses:
			fmt.Println("Handled response from messages")
			for _, client := range e.clients {
				fmt.Println(client.id)
				if client.id == res.ClientId {
					client.responses <- res
				}
			}
		case res := <-messages.Notifications:
			if res.Type == "workspace" && res.Action == "delete" {
				for _, client := range e.clients {
					if client.CurrentWorkspaceId == res.WorkspaceId {
						msg := messages.Message{
							ClientId: client.id,
							Type:     "notification",
							Content:  res,
						}
						client.responses <- msg
					}
				}
			}
			if res.Type == "module" && res.Action == "pop" {
				for _, client := range e.clients {
					if client.CurrentWorkspaceId == res.WorkspaceId && client.id != res.NotConcerned {
						msg := messages.Message{
							ClientId: client.id,
							Type:     "notification",
							Content:  res,
						}
						client.responses <- msg
					}
				}
			}
		// Prepare message when its a log event
		case res := <-logging.ForwardLogs:
			fmt.Println("handled event from logger")
			for _, client := range e.clients {
				if client.CurrentWorkspaceId == res.Data["workspaceId"] {
					client.Logger.Forward(res)
				}
			}
		}
	}
}

func (e *Endpoint) DispatchRequest(req messages.ClientRequest) {
	// 1. Check commands: most of them can be run in either context
	// 2. For context-sensitive commands, check context
	fmt.Println(req.Command[0])
	switch req.Command[0] {
	// Server
	case "server":
		fmt.Println("launching handleServer")
		messages.ForwardServerManager <- req
	// Log
	case "log":
		messages.ForwardLogger <- req
		fmt.Println("Launching handleLog")
	// Stack
	case "stack":
		fmt.Println("Launching handleModule for stack")
		messages.ForwardModuleStack <- req
	// Workspace
	case "workspace":
		fmt.Println("Launching handleWorkspace")
		messages.ForwardWorkspace <- req
	// Module
	case "run", "show", "reload", "module":
		fmt.Println("launching handleModule")
		messages.ForwardModuleStack <- req
	// Compiler:
	case "list", "compile", "compiler":
		fmt.Println("Dispatched request to handleCompiler")
		messages.ForwardCompiler <- req
	// Agent
	case "agent", "interact", "cmd", "back", "download",
		"execute-shellcode", "kill", "main", "shell", "upload":
		fmt.Println("Launching handleAgent")
	// For both commands we need to check context
	case "use", "info", "set":
		switch req.Context {
		case "main":
			fmt.Println("Launching handleModule")
			messages.ForwardModuleStack <- req
		case "module":
			fmt.Println("Launching handleModule")
			messages.ForwardModuleStack <- req
		case "agent":
			fmt.Println("Launching handleAgent")
		case "compiler":
			messages.ForwardCompiler <- req
		}
	}
}
