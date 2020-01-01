package session

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/compiler"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/modules"
)

// Responses channels
var moduleReqs = make(chan modules.ModuleResponse)
var agentReqs = make(chan messages.AgentResponse)
var logReqs = make(chan messages.LogResponse)
var workspaceReqs = make(chan messages.WorkspaceResponse)
var endpointReqs = make(chan messages.EndpointResponse)
var serverReqs = make(chan messages.ServerResponse)
var stackReqs = make(chan messages.StackResponse)
var compilerReqs = make(chan compiler.CompilerResponse)
var logEventReqs = make(chan map[string]string, 1)

// ----------------------------------------------------------------------
// ENDPOINT LOADING

type Endpoint struct {
	IPAddress   string
	Port        int
	Certificate string
	Key         string
	FQDN        string
	IsDefault   bool
}

func (s *Session) LoadEndpointList() error {
	serverList := []Endpoint{}

	userDir, _ := fs.Expand("~/.wiregost/client/")
	if !fs.Exists(userDir) {
		os.MkdirAll(userDir, 0755)
		fmt.Println(tui.Dim("User directory was not found: creating ~/.wiregost/client/"))
	}
	path, _ := fs.Expand("~/.wiregost/client/server.conf")
	if !fs.Exists(path) {
		fmt.Println(tui.Red("Endpoint Configuration file not found: check for issues," +
			" or run the configuration script again"))
		os.Exit(1)
	} else {
		configBlob, _ := ioutil.ReadFile(path)
		json.Unmarshal(configBlob, &serverList)
	}

	// Format certificate path for each server, add server to EndpointManager
	for _, i := range serverList {
		i.Certificate, _ = fs.Expand(i.Certificate)
		s.SavedEndpoints = append(s.SavedEndpoints,
			Endpoint{IPAddress: i.IPAddress,
				Port:        i.Port,
				Certificate: i.Certificate,
				Key:         i.Key,
				FQDN:        i.FQDN,
				IsDefault:   i.IsDefault})
	}
	return nil
}

func (s *Session) GetDefaultEndpoint() error {
	for _, i := range s.SavedEndpoints {
		if i.IsDefault == true {
			s.CurrentEndpoint = i
			break
		}
	}
	return nil
}

// ----------------------------------------------------------------------
// ENDPOINT CONNECTION

func (s *Session) Send(cmd []string) error {
	msg := messages.ClientRequest{
		UserName:           s.user.Name,
		UserPassword:       s.user.PasswordHashString,
		CurrentWorkspace:   s.currentWorkspace,
		CurrentWorkspaceId: s.CurrentWorkspaceId,
		Context:            s.menuContext,
		CurrentModule:      s.currentModule,
		Command:            cmd,
	}
	enc := json.NewEncoder(s.writer)
	err := enc.Encode(msg)
	if err != nil {
		log.Fatal(err)
	}
	s.writer.Flush()
	return nil
}

func (s *Session) Connect() error {

	// Prepare TLS conf and connect.
	certFile, _ := fs.Expand(s.CurrentEndpoint.Certificate)
	keyFile, _ := fs.Expand(s.CurrentEndpoint.Key)
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	conf := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	address := s.CurrentEndpoint.IPAddress + ":" + strconv.Itoa(s.CurrentEndpoint.Port)
	s.connection, err = tls.Dial("tcp", address, &conf)
	if err != nil {
		fmt.Printf("%s Could not connect to server with given address: %s %s\n", tui.RED, err.Error(), tui.RESET)
		return err
	}

	s.reader = bufio.NewReader(s.connection)
	s.writer = bufio.NewWriter(s.connection)

	// Listen for incoming data
	go func() {
		for {
			var msg json.RawMessage
			env := messages.Message{
				Content: &msg,
			}
			dec := json.NewDecoder(s.reader)
			if err := dec.Decode(&env); err != nil {
				fmt.Println("Failed to decode raw message: " + err.Error())
				break
			}
			switch env.Type {
			case "module":
				var mod modules.ModuleResponse
				if err := json.Unmarshal(msg, &mod); err != nil {
					fmt.Println("Failed to decode Module message: " + err.Error())
				}
				moduleReqs <- mod
			case "agent":
				var agent messages.AgentResponse
				if err := json.Unmarshal(msg, &agent); err != nil {
					fmt.Println("Failed to decode agent reponse:")
				}
				agentReqs <- agent
			case "log":
				var log messages.LogResponse
				if err := json.Unmarshal(msg, &log); err != nil {
					fmt.Println("Failed to decode log response")
				}
				logReqs <- log
			case "stack":
				var stack messages.StackResponse
				if err := json.Unmarshal(msg, &stack); err != nil {
					fmt.Println("Failed to decode log response")
				}
				stackReqs <- stack
			case "workspace":
				var workspace messages.WorkspaceResponse
				if err := json.Unmarshal(msg, &workspace); err != nil {
					fmt.Println("Failed to decode log response")
				}
				workspaceReqs <- workspace
			case "server":
				var server messages.ServerResponse
				if err := json.Unmarshal(msg, &server); err != nil {
					fmt.Println("Failed to decode log response")
				}
				serverReqs <- server
			case "endpoint":
				var endpoint messages.EndpointResponse
				if err := json.Unmarshal(msg, &endpoint); err != nil {
					fmt.Println("Failed to decode log response")
				}
				endpointReqs <- endpoint
			case "compiler":
				var compiler compiler.CompilerResponse
				if err := json.Unmarshal(msg, &compiler); err != nil {
					fmt.Println("Failed to decode log response")
				}
				compilerReqs <- compiler
			case "logEvent":
				var logEvent map[string]string
				if err := json.Unmarshal(msg, &logEvent); err != nil {
					fmt.Println("Failed to decode log response")
				}
				// fmt.Println()
				// s.Logger.handleEvents(logEvent, s)
				fmt.Println(logEvent["message"])
				// fmt.Println()
				// s.Events.Add("server", LogMessage{testlog.INFO, "Test message"})
				// s.Shell.Refresh()
				// logEventReqs <- logEvent
			}
		}
	}()
	return nil
}

func (s *Session) Disconnect() error {
	err := s.connection.Close()
	if err != nil {
		return err
	}
	return nil
}
