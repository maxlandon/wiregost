package cli

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/modules"
)

var (
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
)

// Responses channels
var moduleReqs = make(chan modules.ModuleResponse)
var agentReqs = make(chan messages.AgentResponse)
var logReqs = make(chan messages.LogResponse)
var workspaceReqs = make(chan messages.WorkspaceResponse)
var serverReqs = make(chan messages.ServerResponse)
var stackReqs = make(chan messages.StackResponse)

func (s *Session) Send(cmd []string) error {
	msg := messages.ClientRequest{
		UserName:           s.user.Name,
		UserPassword:       s.user.PasswordHashString,
		CurrentWorkspace:   s.currentWorkspace,
		CurrentWorkspaceId: s.CurrentWorkspaceId,
		Context:            s.shellMenuContext,
		Command:            cmd,
	}
	enc := json.NewEncoder(writer)
	err := enc.Encode(msg)
	if err != nil {
		log.Fatal(err)
	}
	writer.Flush()
	return nil
}

func Connect() error {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", ":5000", conf)
	if err != nil {
		log.Fatal(err)
	}
	connection = conn
	reader = bufio.NewReader(connection)
	writer = bufio.NewWriter(connection)
	go func() {
		for {
			var msg json.RawMessage
			env := messages.Message{
				Content: &msg,
			}
			dec := json.NewDecoder(reader)
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
			}
		}
	}()
	return nil

}

func Disconnect() error {
	err := connection.Close()
	if err != nil {
		return err
	}
	return nil
}
