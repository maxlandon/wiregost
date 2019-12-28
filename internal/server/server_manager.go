package server

import (
	"fmt"
	"strconv"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/dispatch"
	"github.com/maxlandon/wiregost/internal/messages"
)

// var ServerReqs = make(chan messages.ClientRequest)

type ServerManager struct {
	// Servers
	Servers map[int]Server
}

// Each time the Manager has to spawn a Server, it should load its required parameters.
// Ideally, it should load them only before spawn, and not keep them as state for too long

func NewServerManager() *ServerManager {
	manager := &ServerManager{
		Servers: make(map[int]Server),
	}
	// Handle requests
	go manager.handleClientRequests()

	return manager
}

func (sm *ServerManager) handleWorkspaceRequests() {
	for {
		// request := <- workspace.ServerRequests
		// Here spawn a server based on parameters of the file loaded from path
	}
}

func (sm *ServerManager) handleClientRequests() {
	for {
		request := <-dispatch.ForwardServerManager
		sm.ReloadServer(request)
	}
}

// This function instantiates a new Server object when starting Wiregost and all saved workspaces
func (sm *ServerManager) SpawnServer() {
	// 1. Load configuration from file
}

// This function instantiates a new Server object upon request of a client
func (sm *ServerManager) ReloadServer(request messages.ClientRequest) {
	params := request.ServerParams
	port, _ := strconv.Atoi(params["server.port"])
	server, err := New(params["server.address"], port,
		params["server.protocol"], params["server.key"],
		params["server.certificate"], params["server.psk"])
	server.Running = false

	status := ""
	if err != nil {
		status = fmt.Sprintf("%s[!]%s There was an error creating a new server instance:\r\n%s", tui.RED, tui.RESET, err.Error())
	} else {
		status = fmt.Sprintf("%s[-]%s HTTP2 Server (%s) ready to listen on %s:%s",
			tui.GREEN, tui.RESET, server.ID, server.Interface, strconv.Itoa(server.Port))
		// Here spawn a server based on parameters of the client request
		sm.Servers[request.CurrentWorkspaceId] = server
	}
	response := messages.ServerResponse{
		Status: status,
	}
	msg := messages.Message{
		ClientId: request.ClientId,
		Type:     "server",
		Content:  response,
	}
	dispatch.Responses <- msg

}
