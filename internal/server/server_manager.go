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
			sm.SpawnServer(request.WorkspacePath, request.WorkspaceId)
		}
		// Here spawn a server based on parameters of the file loaded from path
	}
}

func (sm *ServerManager) handleClientRequests() {
	for {
		request := <-dispatch.ForwardServerManager
		sm.ReloadServer(request)
	}
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
func (sm *ServerManager) SpawnServer(path string, id int) {
	// 1. Load configuration from file
	template := Server{}
	configBlob, _ := ioutil.ReadFile(path + "/" + "server.conf")
	json.Unmarshal(configBlob, &template)
	fmt.Println(path)
	fmt.Println(template)

	// Instantiate server, and attach it to manager
	server, _ := New(template.Interface, template.Port, template.Protocol, template.Key, template.Certificate, template.psk)
	server.Running = false
	sm.Servers[id] = server
	fmt.Println(sm.Servers)
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
