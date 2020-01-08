package session

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/compiler"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/modules"
)

// ----------------------------------------------------------------------
// ENDPOINT CONNECTION

func (s *Session) connectionStatus(conn messages.EndpointResponse) {
	switch conn.Status {
	case "authenticated":
		fmt.Printf("Connected as " + tui.Bold(tui.Yellow(s.user.Name)+" (Administrator rights)."))
		fmt.Println(" Server at " + s.CurrentEndpoint.IPAddress + ":" + strconv.Itoa(s.CurrentEndpoint.Port) +
			" (FQDN: " + s.CurrentEndpoint.FQDN + ", default: " + strconv.FormatBool(s.CurrentEndpoint.IsDefault) + ")")
		fmt.Println()
	case "rejected":
		fmt.Printf("%s[!]%s Connection closed: user authentication failed.\n", tui.RED, tui.RESET)
		fmt.Println()
	}
}

func (s *Session) send(cmd []string) error {
	msg := messages.ClientRequest{
		UserName:           s.user.Name,
		UserPassword:       s.user.PasswordHashString,
		CurrentWorkspace:   s.currentWorkspace,
		CurrentWorkspaceID: s.CurrentWorkspaceID,
		Context:            s.menuContext,
		CurrentModule:      s.currentModule,
		CurrentServerID:    s.currentServerID,
		CurrentAgentID:     s.currentAgentID,
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

func (s *Session) connect() error {

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

	// Send authenticated connection request
	s.send([]string{"connect"})

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
			case "connection":
				var conn messages.EndpointResponse
				if err := json.Unmarshal(msg, &conn); err != nil {
					fmt.Println("Failed to decode Module message: " + err.Error())
				}
				s.connectionStatus(conn)
			case "module":
				var mod modules.ModuleResponse
				if err := json.Unmarshal(msg, &mod); err != nil {
					fmt.Println("Failed to decode Module message: " + err.Error())
				}
				s.moduleReqs <- mod
			case "agent":
				var agent messages.AgentResponse
				if err := json.Unmarshal(msg, &agent); err != nil {
					fmt.Println("Failed to decode agent reponse:")
				}
				s.agentReqs <- agent
			case "log":
				var log messages.LogResponse
				if err := json.Unmarshal(msg, &log); err != nil {
					fmt.Println("Failed to decode log response")
				}
				s.logReqs <- log
			case "workspace":
				var workspace messages.WorkspaceResponse
				if err := json.Unmarshal(msg, &workspace); err != nil {
					fmt.Println("Failed to decode log response")
				}
				s.workspaceReqs <- workspace
			case "server":
				var server messages.ServerResponse
				if err := json.Unmarshal(msg, &server); err != nil {
					fmt.Println("Failed to decode log response")
				}
				s.serverReqs <- server
			case "endpoint":
				var endpoint messages.EndpointResponse
				if err := json.Unmarshal(msg, &endpoint); err != nil {
					fmt.Println("Failed to decode log response")
				}
				s.endpointReqs <- endpoint
			case "compiler":
				var compiler compiler.Response
				if err := json.Unmarshal(msg, &compiler); err != nil {
					fmt.Println("Failed to decode log response")
				}
				s.compilerReqs <- compiler
			case "logEvent":
				var event messages.LogEvent
				if err := json.Unmarshal(msg, &event); err != nil {
					fmt.Println("Failed to decode log response")
				}
				// Print event to console
				Log(event)
			case "notification":
				var notif messages.Notification
				if err := json.Unmarshal(msg, &notif); err != nil {
					fmt.Println("Failed to decode log response")
				}
				s.handleNotifications(notif)
			}
		}
	}()
	return nil
}

func (s *Session) disconnect() error {
	err := s.connection.Close()
	if err != nil {
		return err
	}
	return nil
}
