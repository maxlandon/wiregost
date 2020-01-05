package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/dispatch"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/util"
	"github.com/maxlandon/wiregost/internal/workspace"
	"github.com/tjarratt/babble"
)

type ServerManager struct {
	// Servers
	Servers map[int]*Server
}

// Each time the Manager has to spawn a Server, it should load its required parameters.
// Ideally, it should load them only before spawn, and not keep them as state for too long

func NewServerManager() *ServerManager {
	manager := &ServerManager{
		Servers: make(map[int]*Server),
	}
	// Handle requests
	go manager.handleClientRequests()
	go manager.handleWorkspaceRequests()

	return manager
}

func (sm *ServerManager) handleWorkspaceRequests() {
	for {
		request := <-workspace.ServerRequests
		switch request.Action {
		case "create":
			sm.CreateServer(request.WorkspacePath)
			sm.LoadServer(request)
		case "spawn":
			sm.LoadServer(request)
		case "delete":
			delete(sm.Servers, request.WorkspaceId)
		case "status":
			sm.GiveStatus(request)
		}
	}
}

func (sm *ServerManager) handleClientRequests() {
	for {
		request := <-dispatch.ForwardServerManager
		fmt.Println()
		fmt.Println()
		fmt.Println(sm.Servers)
		switch request.Command[1] {
		case "start":
			sm.StartServer(request)
		case "stop":
			sm.StopServer(request)
		case "list":
			sm.ListServers(request)
		case "reload":
			sm.ReloadServer(request)
		case "generate_certificate":
			sm.GenerateCertificate(request)
		}
	}
}

func (sm *ServerManager) GiveStatus(request workspace.ServerRequest) {
	serv := sm.Servers[request.WorkspaceId]
	var status string
	fmt.Println(sm.Servers[request.WorkspaceId].Running)
	switch sm.Servers[request.WorkspaceId].Running {
	case true:
		status = fmt.Sprintf("%s[*]%s HTTP/2 Server listening on %s:%s",
			tui.GREEN, tui.RESET, serv.Interface, strconv.Itoa(serv.Port))
	case false:
		status = fmt.Sprintf("%s[*]%s HTTP/2 Server ready to listen on %s:%s",
			tui.GREEN, tui.RESET, serv.Interface, strconv.Itoa(serv.Port))
	}
	res := messages.ServerResponse{
		Status: status,
	}
	msg := messages.Message{
		ClientId: request.ClientId,
		Type:     "server",
		Content:  res,
	}
	dispatch.Responses <- msg
}

func (sm *ServerManager) StartServer(request messages.ClientRequest) {
	status, _ := sm.Servers[request.CurrentWorkspaceId].Run()
	res := messages.ServerResponse{
		User:   request.UserName,
		Status: status,
	}
	msg := messages.Message{
		ClientId: request.ClientId,
		Type:     "server",
		Content:  res,
	}
	fmt.Println(sm.Servers[request.CurrentWorkspaceId].Running)
	dispatch.Responses <- msg
}

func (sm *ServerManager) StopServer(request messages.ClientRequest) {
	// 1. Load configuration from file
	path, _ := fs.Expand("~/.wiregost/workspaces/")
	template := Server{}
	configBlob, _ := ioutil.ReadFile(path + "/" + request.CurrentWorkspace + "/" + "server.conf")
	fmt.Println(configBlob)
	err := json.Unmarshal(configBlob, &template)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Reuse the same workspace logger for the new server
	reuseLogger := sm.Servers[request.CurrentWorkspaceId].log

	// Instantiate server, and attach it to manager
	server, _ := New(template.Interface, template.Port, template.Protocol, template.Key, template.Certificate,
		template.Psk, template.Workspace, template.WorkspaceId, reuseLogger)
	sm.Servers[request.CurrentWorkspaceId] = &server

	// Send response
	status := fmt.Sprintf("%s[-]%s HTTP2 Server has been stopped.", tui.GREEN, tui.RESET)
	res := messages.ServerResponse{
		User:   request.UserName,
		Status: status,
	}
	msg := messages.Message{
		ClientId: request.ClientId,
		Type:     "server",
		Content:  res,
	}
	dispatch.Responses <- msg
}

func (sm *ServerManager) ListServers(request messages.ClientRequest) {
	servers := make([]map[string]string, 0)
	for _, v := range sm.Servers {
		list := make(map[string]string)
		list["workspace"] = v.Workspace
		list["address"] = v.Interface + ":" + strconv.Itoa(v.Port)
		list["protocol"] = v.Protocol
		list["state"] = strconv.FormatBool(v.Running)
		list["psk"] = v.Psk
		list["certificate"] = v.Certificate
		servers = append(servers, list)
	}
	res := messages.ServerResponse{
		ServerList: servers,
	}
	msg := messages.Message{
		ClientId: request.ClientId,
		Type:     "server",
		Content:  res,
	}
	dispatch.Responses <- msg
}

func (sm *ServerManager) GenerateCertificate(req messages.ClientRequest) {
	// Create server object, associated certificates, key and populate
	// Private Key
	path, _ := fs.Expand("~/.wiregost/workspaces/" + req.CurrentWorkspace)
	name := req.Command[2]
	fmt.Println(tui.Dim("Generating RSA private key"))
	privKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	file, _ := os.Create(path + "/" + name + ".key")
	file.Close()
	writeKey, _ := fs.Expand(path + "/" + name + ".key")

	privkey_bytes := x509.MarshalPKCS1PrivateKey(privKey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		})

	ioutil.WriteFile(writeKey, privkey_pem, 0644)
	fmt.Println(tui.Dim("Ok"))

	// Certificate
	fmt.Println(tui.Dim("Generating TLS Certificate from private key"))
	cert_bytes, _ := util.GenerateTLSCert(nil, nil, nil, nil, nil, privKey, true)
	file, _ = os.Create(path + "/" + name + ".crt")
	file.Close()
	writeCert, _ := fs.Expand(path + "/" + name + ".crt")

	cert := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert_bytes.Certificate[0],
		})

	ioutil.WriteFile(writeCert, cert, 0644)
	fmt.Println(tui.Dim("Ok"))
	status := fmt.Sprintf("%s[*]%s Certificate '%s.crt' and private key '%s.key' created in %s/ directory.",
		tui.GREEN, tui.RESET, name, name, req.CurrentWorkspace)
	res := messages.ServerResponse{
		User:   req.UserName,
		Status: status,
	}
	msg := messages.Message{
		ClientId: req.ClientId,
		Type:     "server",
		Content:  res,
	}
	dispatch.Responses <- msg

}

func (sm *ServerManager) CreateServer(path string) {
	// Create configuration file
	serverConf, _ := os.Create(path + "/" + "server.conf")
	defer serverConf.Close()
	w := strings.Split(path, "/")
	name := w[len(w)-1]

	// Create server object, associated certificates, key and populate Private Key
	fmt.Println(tui.Dim("Generating RSA private key"))
	privKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	file, _ := os.Create(path + "/" + name + ".key")
	file.Close()
	writeKey, _ := fs.Expand(path + "/" + name + ".key")

	privkey_bytes := x509.MarshalPKCS1PrivateKey(privKey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		})

	ioutil.WriteFile(writeKey, privkey_pem, 0644)
	fmt.Println(tui.Dim("Ok"))

	// Certificate
	fmt.Println(tui.Dim("Generating TLS Certificate from private key"))
	cert_bytes, _ := util.GenerateTLSCert(nil, nil, nil, nil, nil, privKey, true)
	file, _ = os.Create(path + "/" + name + ".crt")
	file.Close()
	writeCert, _ := fs.Expand(path + "/" + name + ".crt")

	cert := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert_bytes.Certificate[0],
		})

	ioutil.WriteFile(writeCert, cert, 0644)
	fmt.Println(tui.Dim("Ok"))

	// Psk
	babbler := babble.NewBabbler()
	babbler.Count = 1
	psk := babbler.Babble()
	fmt.Println("Generated PSK: " + psk)

	// Server
	server := Server{
		Interface:   "127.0.0.1",
		Port:        sm.FindFreePort(),
		Protocol:    "h2",
		Certificate: writeCert,
		Key:         writeKey,
		Psk:         psk,
	}

	var jsonData []byte
	jsonData, _ = json.MarshalIndent(server, "", "    ")
	confFile, _ := fs.Expand(path + "/" + "server.conf")
	ioutil.WriteFile(confFile, jsonData, 0755)
	fmt.Println("Written server.conf file")
}

// This function instantiates a new Server object when starting Wiregost and all saved workspaces
func (sm *ServerManager) LoadServer(request workspace.ServerRequest) {
	// 1. Load configuration from file
	template := Server{}
	configBlob, _ := ioutil.ReadFile(request.WorkspacePath + "/" + "server.conf")
	json.Unmarshal(configBlob, &template)

	// Instantiate server, and attach it to manager
	server, _ := New(template.Interface, template.Port, template.Protocol, template.Key, template.Certificate,
		template.Psk, request.Logger.WorkspaceName, request.Logger.WorkspaceId, request.Logger)
	sm.Servers[request.Logger.WorkspaceId] = &server
}

// This function instantiates a new Server object upon request of a client
func (sm *ServerManager) ReloadServer(request messages.ClientRequest) {
	// Reuse the same workspace logger for the new server
	reuseLogger := sm.Servers[request.CurrentWorkspaceId].log

	// Load pushed params if there are some.
	currentServer := sm.Servers[request.CurrentWorkspaceId]
	params := request.ServerParams
	fmt.Println(params)
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
		newParams["server.certificate"], newParams["server.psk"], request.CurrentWorkspace, request.CurrentWorkspaceId, reuseLogger)

	status := ""
	if err != nil {
		status = fmt.Sprintf("%s[!]%s There was an error creating a new server instance:\r\n%s", tui.RED, tui.RESET, err.Error())
	} else {
		status = fmt.Sprintf("%s[-]%s HTTP2 Server ready to listen on %s:%s, with provided parameters. (Configuration saved)",
			tui.GREEN, tui.RESET, server.Interface, strconv.Itoa(server.Port))
		// Here spawn a server based on parameters of the client request
		sm.Servers[request.CurrentWorkspaceId] = &server

		// Create a server struct and populate it, then save to configuration file.
		template := Server{
			Interface:   newParams["server.address"],
			Port:        port,
			Protocol:    newParams["server.protocol"],
			Certificate: newParams["server.certificate"],
			Key:         newParams["server.key"],
			Psk:         newParams["server.psk"],
			Workspace:   request.CurrentWorkspace,
			WorkspaceId: request.CurrentWorkspaceId,
		}
		path, _ := fs.Expand("~/.wiregost/workspaces/" + request.CurrentWorkspace)
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
		fmt.Println("Written server.conf file")
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

func (sm *ServerManager) FindFreePort() (port int) {
	freePort := 0
	for _, s := range sm.Servers {
		if s.Port > freePort {
			freePort = s.Port
		}
	}
	return freePort + 1
}
