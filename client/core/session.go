// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package core

import (
	// Standard
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	// 3rd party
	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/tui"
	uuid "github.com/satori/go.uuid"

	// Wiregost
	"github.com/maxlandon/wiregost/data_service/models"
	"github.com/maxlandon/wiregost/data_service/remote"
)

// Session is the base object of a Client session.
type Session struct {
	// Shell
	Shell  *readline.Instance
	prompt Prompt

	// Auth
	// user *User

	// Context
	menuContext        string
	currentModule      string
	currentWorkspace   models.Workspace
	currentWorkspaceID int
	currentAgentID     uuid.UUID

	// Server state
	currentServerID uuid.UUID
	serverRunning   bool

	// Environmment variables
	Env map[string]string

	// Server connection parameters
	// SavedEndpoints   []Endpoint
	// CurrentEndpoint  Endpoint
	endpointPublicIP string
	connected        bool
}

func NewSession() *Session {
	session := &Session{
		menuContext: "main",
		Env:         make(map[string]string),
	}

	home, _ := os.UserHomeDir()

	// Set Shell and Completers
	// shellCompleter := session.getCompleter("main")
	session.Shell, _ = readline.NewEx(&readline.Config{
		HistoryFile: home + "/.wiregost/client/.history",
		// AutoComplete:      shellCompleter,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		// FilterInputRune: To be used later if needed
	})
	// Set Prompt
	session.prompt = newPrompt(session)

	// Get default workspace
	workspaces, err := remote.Workspaces()
	if err != nil {
		fmt.Println(tui.Red("Failed to fetch workspaces"))
	}
	for _, w := range workspaces {
		if w.IsDefault {
			session.currentWorkspace = w
		}
	}

	// Launch console but give time to connect
	time.Sleep(time.Millisecond * 50)
	session.start()

	return session
}

func (s *Session) start() {
	// Eventually close the session
	defer s.Shell.Close()

	// Read commands
	for {
		line, err := s.Shell.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		cmd := strings.Fields(line)
		//
		if len(cmd) > 0 {
			fmt.Println(cmd)
			//         switch s.menuContext {
			//         case "main":
			//                 s.mainMenuCommand(cmd)
			//         case "module":
			//                 s.moduleMenuCommand(cmd)
			//         case "compiler":
			//                 s.compilerMenuCommand(cmd)
			//         case "agent":
			//                 s.agentMenuCommand(cmd)
			//         }
		}

		// Refresh shell & prompt after each command, at least.
		s.refresh()
	}
}

func (s *Session) refresh() {
	refreshPrompt(s.prompt, s.Shell)
	s.Shell.Refresh()
}
