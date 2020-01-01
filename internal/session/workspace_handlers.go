package session

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/messages"
)

func (s *Session) WorkspaceList(cmd []string) {
	s.Send(cmd)
	workspace := <-workspaceReqs
	fmt.Println(workspace)
}

func (s *Session) WorkspaceSwitch(cmd []string) {
	s.currentWorkspace = cmd[2]
	s.Send(cmd)
	workspace := <-workspaceReqs
	s.CurrentWorkspaceId = workspace.WorkspaceId
}

func (s *Session) WorkspaceDelete(cmd []string) {
	if cmd[2] == s.currentWorkspace {
		fmt.Println()
		fmt.Printf("%s[!]%s Cannot delete current workspace", tui.RED, tui.RESET)
		fmt.Println()
	} else {
		s.Send(cmd)
		workspace := <-workspaceReqs
		fmt.Println()
		fmt.Println(workspace.Result)
	}
}

func (s *Session) WorkspaceNew(cmd []string) {
	// Send params if they are set
	params := make(map[string]string)
	for k, v := range s.Env {
		if strings.HasPrefix(k, "workspace") {
			params[k] = v
		}
	}
	msg := messages.ClientRequest{
		UserName:           s.user.Name,
		UserPassword:       s.user.PasswordHashString,
		CurrentWorkspace:   s.currentWorkspace,
		CurrentWorkspaceId: s.CurrentWorkspaceId,
		Context:            s.menuContext,
		CurrentModule:      s.currentModule,
		Command:            cmd,
		WorkspaceParams:    params,
	}
	enc := json.NewEncoder(s.writer)
	err := enc.Encode(msg)
	if err != nil {
		log.Fatal(err)
	}
	s.writer.Flush()
	workspace := <-workspaceReqs
	fmt.Println()
	fmt.Println(workspace.Result)
}
