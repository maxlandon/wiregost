package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/util"
	"github.com/maxlandon/wiregost/internal/workspace"
	"github.com/tjarratt/babble"
)

// Manager stores all instantiated agent servers and perform operations on them,
// (loading, starting, stopping, populating, generating certificates, etc...)
type Manager struct {
	// Servers
	Servers map[int]*Server
}

// NewManager instantiates a Server Manager, which will handle
// all requests from clients or from workspaces.
func NewManager() *Manager {
	manager := &Manager{
		Servers: make(map[int]*Server),
	}
	// Handle requests
	go manager.handleClientRequests()
	go manager.handleWorkspaceRequests()
	go manager.handleEndpointRequests()

	return manager
}

func (sm *Manager) handleWorkspaceRequests() {
	for {
		request := <-workspace.ServerRequests
		switch request.Action {
		case "create":
			sm.createServer(request.WorkspacePath)
			sm.loadServer(request)
		case "spawn":
			sm.loadServer(request)
		case "delete":
			// Remove all agents tied to server
			agentReq := messages.AgentRequest{
				Action:   "delete_all",
				ServerID: sm.Servers[request.WorkspaceID].ID,
			}
			messages.AgentRequests <- agentReq
			// Give time to server to send kill messages
			time.Sleep(time.Second * 30)
			delete(sm.Servers, request.WorkspaceID)
		case "status":
			sm.giveStatus(request)
		}
	}
}

func (sm *Manager) handleClientRequests() {
	for {
		request := <-messages.ForwardServerManager
		switch request.Command[1] {
		case "start":
			sm.startServer(request)
		case "stop":
			sm.stopServer(request)
		case "list":
			sm.listServers(request)
		case "reload":
			sm.reloadServer(request)
		case "generate_certificate":
			sm.generateCertificate(request)
		}
	}
}

func (sm *Manager) handleEndpointRequests() {
	for {
		req := <-messages.EndpointToServer
		switch req.Command[2] {
		case "default":
			for _, s := range sm.Servers {
				if s.Workspace == "default" {
					res := messages.ServerResponse{
						ServerID:      s.ID,
						ServerRunning: s.Running,
					}
					messages.ForwardServer <- res
				}
			}
		}
	}
}

func (sm *Manager) giveStatus(request workspace.ServerRequest) {
	serv := sm.Servers[request.WorkspaceID]
	var status string
	switch sm.Servers[request.WorkspaceID].Running {
	case true:
		status = fmt.Sprintf("%s[*]%s HTTP/2 Server listening on %s:%s",
			tui.GREEN, tui.RESET, serv.Interface, strconv.Itoa(serv.Port))
	case false:
		status = fmt.Sprintf("%s[*]%s HTTP/2 Server ready to listen on %s:%s",
			tui.GREEN, tui.RESET, serv.Interface, strconv.Itoa(serv.Port))
	}
	res := messages.ServerResponse{
		Status:        status,
		ServerID:      serv.ID,
		ServerRunning: serv.Running,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "server",
		Content:  res,
	}
	messages.Responses <- msg
}

func (sm *Manager) startServer(request messages.ClientRequest) {
	s := sm.Servers[request.CurrentWorkspaceID]
	go s.Run()
	m := fmt.Sprintf("%s[*]%s Starting %s listener on %s:%d %s(pre-shared key: %s%s)",
		tui.GREEN, tui.RESET, s.Protocol, s.Interface, s.Port, tui.DIM, s.Psk, tui.RESET)
	// Give server time to start and set its status
	time.Sleep(time.Millisecond * 50)

	res := messages.ServerResponse{
		User:          request.UserName,
		ServerRunning: s.Running,
		Status:        m,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "server",
		Content:  res,
	}
	messages.Responses <- msg
	// Send notification
	notif := messages.Notification{
		Type:          "server",
		WorkspaceID:   request.CurrentWorkspaceID,
		NotConcerned:  request.ClientID,
		ServerID:      s.ID,
		ServerRunning: s.Running,
	}
	messages.Notifications <- notif
}

func (sm *Manager) stopServer(request messages.ClientRequest) {
	// Stop previous server before reinstantiating it. Will need to change this if other than HTTP/2 is available
	sm.Servers[request.CurrentWorkspaceID].Server.(*http.Server).Close()
	// 1. Load configuration from file
	path, _ := fs.Expand("~/.wiregost/server/workspaces/")
	template := Server{}
	configBlob, _ := ioutil.ReadFile(path + "/" + request.CurrentWorkspace + "/" + "server.conf")
	err := json.Unmarshal(configBlob, &template)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Reuse the same workspace logger for the new server
	reuseLogger := sm.Servers[request.CurrentWorkspaceID].log

	// Instantiate server, and attach it to manager
	server, _ := New(template.Interface, template.Port, template.Protocol, template.Key, template.Certificate,
		template.Psk, template.Workspace, template.WorkspaceID, reuseLogger)
	sm.Servers[request.CurrentWorkspaceID] = &server

	// Send response
	status := fmt.Sprintf("%s[-]%s HTTP2 Server has been stopped.", tui.GREEN, tui.RESET)
	res := messages.ServerResponse{
		User:          request.UserName,
		ServerRunning: server.Running,
		Status:        status,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "server",
		Content:  res,
	}
	messages.Responses <- msg
	// Send notification
	notif := messages.Notification{
		Type:          "server",
		WorkspaceID:   request.CurrentWorkspaceID,
		NotConcerned:  request.ClientID,
		ServerID:      server.ID,
		ServerRunning: server.Running,
	}
	messages.Notifications <- notif
}

func (sm *Manager) listServers(request messages.ClientRequest) {
	servers := make([]map[string]string, 0)
	for _, v := range sm.Servers {
		list := make(map[string]string)
		list["workspace"] = v.Workspace
		list["address"] = v.Interface + ":" + strconv.Itoa(v.Port)
		list["protocol"] = v.Protocol
		list["state"] = strconv.FormatBool(v.Running)
		list["psk"] = v.Psk
		list["certificate"] = v.Certificate
		list["id"] = v.ID.String()
		servers = append(servers, list)
	}
	res := messages.ServerResponse{
		ServerList: servers,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "server",
		Content:  res,
	}
	messages.Responses <- msg
}

func (sm *Manager) generateCertificate(req messages.ClientRequest) {
	// Create server object, associated certificates, key and populate
	// Private Key
	path, _ := fs.Expand("~/.wiregost/server/workspaces/" + req.CurrentWorkspace)
	name := req.Command[2]
	privKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	file, _ := os.Create(path + "/" + name + ".key")
	file.Close()
	writeKey, _ := fs.Expand(path + "/" + name + ".key")

	privkeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		})

	ioutil.WriteFile(writeKey, privkeyPem, 0644)

	// Certificate
	certBytes, _ := util.GenerateTLSCert(nil, nil, nil, nil, nil, privKey, true)
	file, _ = os.Create(path + "/" + name + ".crt")
	file.Close()
	writeCert, _ := fs.Expand(path + "/" + name + ".crt")

	cert := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes.Certificate[0],
		})

	ioutil.WriteFile(writeCert, cert, 0644)
	status := fmt.Sprintf("%s[*]%s Certificate '%s.crt' and private key '%s.key' created in %s/ directory.",
		tui.GREEN, tui.RESET, name, name, req.CurrentWorkspace)
	res := messages.ServerResponse{
		User:   req.UserName,
		Status: status,
	}
	msg := messages.Message{
		ClientID: req.ClientID,
		Type:     "server",
		Content:  res,
	}
	messages.Responses <- msg

}

func (sm *Manager) createServer(path string) {
	// Create configuration file
	serverConf, _ := os.Create(path + "/" + "server.conf")
	defer serverConf.Close()
	w := strings.Split(path, "/")
	name := w[len(w)-1]

	// Create server object, associated certificates, key and populate Private Key
	privKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	file, _ := os.Create(path + "/" + name + ".key")
	file.Close()
	writeKey, _ := fs.Expand(path + "/" + name + ".key")

	privkeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		})

	ioutil.WriteFile(writeKey, privkeyPem, 0644)

	// Certificate
	certBytes, _ := util.GenerateTLSCert(nil, nil, nil, nil, nil, privKey, true)
	file, _ = os.Create(path + "/" + name + ".crt")
	file.Close()
	writeCert, _ := fs.Expand(path + "/" + name + ".crt")

	cert := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes.Certificate[0],
		})

	ioutil.WriteFile(writeCert, cert, 0644)

	// Psk
	babbler := babble.NewBabbler()
	babbler.Count = 1
	psk := babbler.Babble()

	// Server
	server := Server{
		Interface:   "127.0.0.1",
		Port:        sm.findFreePort(),
		Protocol:    "h2",
		Certificate: writeCert,
		Key:         writeKey,
		Psk:         psk,
	}

	var jsonData []byte
	jsonData, _ = json.MarshalIndent(server, "", "    ")
	confFile, _ := fs.Expand(path + "/" + "server.conf")
	ioutil.WriteFile(confFile, jsonData, 0755)
}

// This function instantiates a new Server object when starting Wiregost and all saved workspaces
func (sm *Manager) loadServer(request workspace.ServerRequest) {
	// 1. Load configuration from file
	template := Server{}
	configBlob, _ := ioutil.ReadFile(request.WorkspacePath + "/" + "server.conf")
	json.Unmarshal(configBlob, &template)

	// Instantiate server, and attach it to manager
	server, _ := New(template.Interface, template.Port, template.Protocol, template.Key, template.Certificate,
		template.Psk, request.Logger.WorkspaceName, request.Logger.WorkspaceID, request.Logger)
	sm.Servers[request.Logger.WorkspaceID] = &server
}

// This function instantiates a new Server object upon request of a client
func (sm *Manager) reloadServer(request messages.ClientRequest) {
	// Reuse the same workspace logger for the new server
	reuseLogger := sm.Servers[request.CurrentWorkspaceID].log

	// Load pushed params if there are some.
	currentServer := sm.Servers[request.CurrentWorkspaceID]
	params := request.ServerParams
	newParams := make(map[string]string, 0)
	if v, ok := params["server.address"]; ok {
		if params["server.address"] != currentServer.Interface {
			newParams["server.address"] = v
		}
	} else {
		newParams["server.address"] = currentServer.Interface
	}
	if v, ok := params["server.port"]; ok {
		if params["server.port"] != strconv.Itoa(currentServer.Port) {
			newParams["server.port"] = v
		}
	} else {
		newParams["server.port"] = strconv.Itoa(currentServer.Port)
	}
	if v, ok := params["server.protocol"]; ok {
		if params["server.protocol"] != currentServer.Protocol {
			newParams["server.protocol"] = v
		}
	} else {
		newParams["server.protocol"] = currentServer.Protocol
	}
	if v, ok := params["server.certificate"]; ok {
		if params["server.certificate"] != currentServer.Certificate {
			newParams["server.certificate"] = v
		}
	} else {
		newParams["server.certificate"] = currentServer.Certificate
	}
	if v, ok := params["server.key"]; ok {
		if params["server.key"] != currentServer.Key {
			newParams["server.key"] = v
		}
	} else {
		newParams["server.key"] = currentServer.Key
	}
	if v, ok := params["server.psk"]; ok {
		if params["server.psk"] != currentServer.Psk {
			newParams["server.psk"] = v
		}
	} else {
		newParams["server.psk"] = currentServer.Psk
	}

	// Reinstantiate server with new parameter set
	port, _ := strconv.Atoi(newParams["server.port"])
	server, err := New(newParams["server.address"], port,
		newParams["server.protocol"], newParams["server.key"],
		newParams["server.certificate"], newParams["server.psk"], request.CurrentWorkspace, request.CurrentWorkspaceID, reuseLogger)

	status := ""
	if err != nil {
		status = fmt.Sprintf("%s[!]%s There was an error creating a new server instance:\r\n%s", tui.RED, tui.RESET, err.Error())
	} else {
		status = fmt.Sprintf("%s[-]%s HTTP2 Server ready to listen on %s:%s, with provided parameters. (Configuration saved)",
			tui.GREEN, tui.RESET, server.Interface, strconv.Itoa(server.Port))
		// Here spawn a server based on parameters of the client request
		sm.Servers[request.CurrentWorkspaceID] = &server

		// Create a server struct and populate it, then save to configuration file.
		template := Server{
			Interface:   newParams["server.address"],
			Port:        port,
			Protocol:    newParams["server.protocol"],
			Certificate: newParams["server.certificate"],
			Key:         newParams["server.key"],
			Psk:         newParams["server.psk"],
			Workspace:   request.CurrentWorkspace,
			WorkspaceID: request.CurrentWorkspaceID,
		}
		path, _ := fs.Expand("~/.wiregost/server/workspaces/" + request.CurrentWorkspace)
		var jsonData []byte
		jsonData, err = json.MarshalIndent(template, "", "    ")
		if err != nil {
			fmt.Println(err.Error())
		}
		confFile, _ := fs.Expand(path + "/" + "server.conf")
		err = ioutil.WriteFile(confFile, jsonData, 0755)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	response := messages.ServerResponse{
		ServerRunning: server.Running,
		Status:        status,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "server",
		Content:  response,
	}
	messages.Responses <- msg
}

func (sm *Manager) findFreePort() (port int) {
	freePort := 0
	for _, s := range sm.Servers {
		if s.Port > freePort {
			freePort = s.Port
		}
	}
	return freePort + 1
}
