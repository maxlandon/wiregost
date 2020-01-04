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

// var ServerReqs = make(chan messages.ClientRequest)

// Var port is used to serve as an incremental port number for spwaning
// different servers
var port = 443

type ServerResponse struct {
	User   string
	Status string
	Error  string
}

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
			sm.CreateConf(request.WorkspacePath)
		case "spawn":
			sm.SpawnServer(request)
		case "delete":
			delete(sm.Servers, request.WorkspaceId)
		case "status":
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
			res := ServerResponse{
				Status: status,
			}
			msg := messages.Message{
				ClientId: request.ClientId,
				Type:     "server",
				Content:  res,
			}
			dispatch.Responses <- msg
		}
		// Here spawn a server based on parameters of the file loaded from path
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
			status, _ := sm.Servers[request.CurrentWorkspaceId].Run()
			res := ServerResponse{
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
		case "stop":
			delete(sm.Servers, request.CurrentWorkspaceId)
			// path, _ := fs.Expand("~/.wiregost/workspaces/" + request.CurrentWorkspace)
			// sm.SpawnServer(path, request.CurrentWorkspaceId)
			status := fmt.Sprintf("%s[-]%s HTTP2 Server has been stopped.", tui.GREEN, tui.RESET)
			res := ServerResponse{
				User:   request.UserName,
				Status: status,
			}
			msg := messages.Message{
				ClientId: request.ClientId,
				Type:     "server",
				Content:  res,
			}
			dispatch.Responses <- msg

		case "reload":
			sm.ReloadServer(request)
		case "generate_certificate":
			sm.GenerateCertificate(request)
		}
	}
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
	res := ServerResponse{
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

func (sm *ServerManager) CreateConf(path string) {
	// Create configuration file
	serverConf, _ := os.Create(path + "/" + "server.conf")
	defer serverConf.Close()
	w := strings.Split(path, "/")
	name := w[len(w)-1]

	// Create server object, associated certificates, key and populate
	// Private Key
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
		Port:        port,
		Protocol:    "h2",
		Certificate: writeCert,
		Key:         writeKey,
		psk:         psk,
	}

	var jsonData []byte
	jsonData, _ = json.MarshalIndent(server, "", "    ")
	confFile, _ := fs.Expand(path + "/" + "server.conf")
	ioutil.WriteFile(confFile, jsonData, 0755)
	fmt.Println("Written server.conf file")
}

// This function instantiates a new Server object when starting Wiregost and all saved workspaces
func (sm *ServerManager) SpawnServer(request workspace.ServerRequest) {
	// 1. Load configuration from file
	template := Server{}
	configBlob, _ := ioutil.ReadFile(request.WorkspacePath + "/" + "server.conf")
	json.Unmarshal(configBlob, &template)
	fmt.Println(request.WorkspacePath)
	fmt.Println(template)

	// Instantiate server, and attach it to manager
	server, _ := New(template.Interface, template.Port, template.Protocol, template.Key,
		template.Certificate, template.psk, request.Logger.WorkspaceName, request.WorkspaceId)
	server.log = request.Logger
	server.Workspace = request.Logger.WorkspaceName
	server.WorkspaceId = request.WorkspaceId
	sm.Servers[request.WorkspaceId] = &server
	fmt.Println(sm.Servers)
}

// This function instantiates a new Server object upon request of a client
func (sm *ServerManager) ReloadServer(request messages.ClientRequest) {
	params := request.ServerParams
	port, _ := strconv.Atoi(params["server.port"])
	server, err := New(params["server.address"], port,
		params["server.protocol"], params["server.key"],
		params["server.certificate"], params["server.psk"], request.CurrentWorkspace, request.CurrentWorkspaceId)

	status := ""
	if err != nil {
		status = fmt.Sprintf("%s[!]%s There was an error creating a new server instance:\r\n%s", tui.RED, tui.RESET, err.Error())
	} else {
		status = fmt.Sprintf("%s[-]%s HTTP2 Server (%s) ready to listen on %s:%s, with provided parameters.",
			tui.GREEN, tui.RESET, server.ID, server.Interface, strconv.Itoa(server.Port))
		// Here spawn a server based on parameters of the client request
		sm.Servers[request.CurrentWorkspaceId] = &server

		path, _ := fs.Expand("~/.wiregost/workspaces/" + request.CurrentWorkspace)
		var jsonData []byte
		jsonData, _ = json.MarshalIndent(server, "", "    ")
		confFile, _ := fs.Expand(path + "/" + "server.conf")
		ioutil.WriteFile(confFile, jsonData, 0755)
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
