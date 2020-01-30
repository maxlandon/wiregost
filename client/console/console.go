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

package console

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/tui"
	"github.com/google/uuid"

	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/completers"
	"github.com/maxlandon/wiregost/data_service/models"
	"github.com/maxlandon/wiregost/data_service/remote"
)

var home, _ = os.UserHomeDir()

type Console struct {
	// Shell
	Shell   *readline.Instance
	prompt  Prompt
	vimMode string

	// Context
	context          context.Context
	menuContext      string
	currentModule    string
	currentWorkspace *models.Workspace

	currentAgentID uuid.UUID
	// Server state
	currentServerID uuid.UUID
	serverRunning   bool
	// Server connection parameters
	// SavedEndpoints   []Endpoint
	// CurrentEndpoint  Endpoint
	endpointPublicIP string
	connected        bool
}

func NewConsole() *Console {

	shell, _ := readline.NewEx(&readline.Config{
		HistoryFile:       home + "/.wiregost/client/.history",
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistoryLimit:      5000,
		HistorySearchFold: true,
		VimMode:           true,
	})

	console := &Console{
		Shell:       shell,
		menuContext: "main",
	}
	console.initContext()
	console.prompt = newPrompt(console)

	// Setup Autocompleter
	completer := &completers.AutoCompleter{
		MenuContext: &console.menuContext,
		Context:     &console.context,
	}
	console.Shell.Config.AutoComplete = completer

	// Set Vim mode
	console.vimMode = "insert"
	console.Shell.Config.FuncFilterInputRune = console.filterInput

	// Register all commands into their respective menus
	commands.RegisterCommands()

	// Launch
	console.Start()

	return console
}

func (c *Console) Start() {

	// Eventually close
	defer c.Shell.Close()

	// Command loop
	for {
		c.vimMode = "insert"
		c.refresh()

		line, err := c.Shell.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}

		unfiltered := strings.Split(line, " ")

		var args []string
		for _, arg := range unfiltered {
			if arg != "" {
				args = append(args, arg)
			}
		}

		if err = ExecCmd(args, c.menuContext, &c.context, c.currentWorkspace, c.currentModule); err != nil {
			fmt.Println(err)
		}
	}
}

func (c *Console) initContext() {
	// Set workspace
	workspaces, err := remote.Workspaces(nil)
	if err != nil {
		fmt.Println(tui.Red("Failed to fetch workspaces"))
	}
	for i, _ := range workspaces {
		if workspaces[i].IsDefault {
			c.currentWorkspace = &workspaces[i]
		}
	}

	// Set context object passed to commands
	rootCtx := context.Background()
	c.context = context.WithValue(rootCtx, "workspace_id", c.currentWorkspace.ID)
}

func (c *Console) refresh() {
	refreshPrompt(c.prompt, c.Shell)
	c.Shell.Refresh()
}

func (c *Console) filterInput(r rune) (rune, bool) {

	switch c.vimMode {
	case "insert":
		switch r {
		case readline.CharEsc:
			c.vimMode = "normal"
			_, m := c.prompt.render()
			c.Shell.SetPrompt(m)
			c.Shell.Refresh()
			return r, true

		case readline.CharCtrlL:
			readline.ClearScreen(c.Shell)
			c.Shell.Refresh()
			c.refresh()
			return r, false
		}
	case "normal":
		switch r {
		case 'i', 'I', 'a', 'A', 's', 'S', 'c':
			c.vimMode = "insert"
			_, p := c.prompt.render()
			c.Shell.SetPrompt(p)
			c.Shell.Refresh()
			return r, true

		case readline.CharCtrlL:
			readline.ClearScreen(c.Shell)
			c.Shell.Refresh()
			c.refresh()
			return r, true
		}
	}
	return r, true
}
