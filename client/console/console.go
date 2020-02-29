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
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/fs"

	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/completers"
	"github.com/maxlandon/wiregost/client/core"
	"github.com/maxlandon/wiregost/data_service/models"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
)

var home, _ = os.UserHomeDir()

type Console struct {
	// Shell
	Shell   *readline.Instance
	prompt  Prompt
	mode    string
	vimMode string

	// Context
	context     context.Context
	menuContext string

	currentModule    string
	module           *clientpb.Module
	moduleUserID     int32
	currentWorkspace *models.Workspace

	// Server
	server *core.WiregostServer

	// Jobs
	listeners int

	// Agents
	ghosts              int
	CurrentAgent        *clientpb.Ghost
	AgentPwd            string
	SessionPathComplete bool

	// CommandShellContext
	shellContext *commands.ShellContext
}

func NewConsole() *Console {

	// [ Config ]
	conf := LoadConsoleConfig()
	history, _ := fs.Expand(conf.HistoryFile)

	// [ New console ]
	console := &Console{
		menuContext:         "main",
		SessionPathComplete: conf.SessionPathCompletion,
	}

	// Set ModuleRequestID
	console.moduleUserID = rand.Int31()

	console.initContext()
	console.prompt = newPrompt(console, conf.Prompt, conf.ImplantPrompt)

	// [ Console input ]
	shell, _ := readline.NewEx(&readline.Config{
		HistoryFile:       history,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistoryLimit:      5000,
		HistorySearchFold: true,
	})

	console.Shell = shell

	// Set keyboard mode
	if conf.Mode == "vim" {
		shell.Config.VimMode = true
	} else if conf.Mode == "emacs" {
		shell.Config.VimMode = false
	} else {
		shell.Config.VimMode = false
	}

	// Set Vim mode
	console.vimMode = "insert"
	console.Shell.Config.FuncFilterInputRune = console.filterInput

	// [ Autocompleters ]
	completer := &completers.AutoCompleter{
		MenuContext: &console.menuContext,
		Context:     console.shellContext,
	}
	console.Shell.Config.AutoComplete = completer

	// [ Commands ]
	commands.RegisterCommands()

	return console
}

func Start() {

	// Instantiate console
	c := NewConsole()

	// Connect to server
	config := getDefaultServerConfig()
	err := c.connect(config)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println()
		os.Exit(1)
	} else {
		go c.eventLoop(c.server)
	}

	// Eventually close
	defer c.Shell.Close()

	// Command loop
	for {
		// Refresh Vim mode each time is needed here
		c.hardRefresh()

		line, err := c.Shell.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			ex := c.exit()
			if ex {
				break
			}
			continue
		}

		line = strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}

		unfiltered := strings.Split(line, " ")

		// Handle exits
		if unfiltered[0] == "exit" {
			ex := c.exit()
			if ex {
				break
			}
			continue
		}

		// Sanitize input
		var args []string
		for _, arg := range unfiltered {
			if arg != "" {
				args = append(args, arg)
			}
		}

		// Exec command
		if err = ExecCmd(args, c.menuContext, c.shellContext); err != nil {
			fmt.Println(err)
		}

	}
}

// [ Generic console functions ] ---------------------------------------------------------------//

// hardRefresh prints a new prompt
func (c *Console) hardRefresh() {
	// Input
	if c.mode == "vim" {
		c.Shell.Config.VimMode = true
		c.vimMode = "insert"
	} else if c.mode == "emacs" {
		c.Shell.Config.VimMode = false
	}

	// Menu context
	if c.currentModule != "" {
		c.menuContext = "module"
	} else {
		c.menuContext = "main"
	}
	if c.CurrentAgent.Name != "" {
		c.menuContext = "agent"
	}

	// Jobs
	jobs := commands.GetJobs(c.shellContext.Server.RPC)
	c.listeners = len(jobs.Active)

	// Sessions
	sessions := commands.GetGhosts(c.shellContext.Server.RPC)
	c.ghosts = len(sessions.Ghosts)

	// Prompt
	refreshPrompt(c.prompt, c.Shell)
	c.Shell.Refresh()
}

// softRefresh does not print a new prompt, it simply updates the current one
func (c *Console) softRefresh() {
}

func (c *Console) filterInput(r rune) (rune, bool) {

	switch c.Shell.IsVimMode() {
	// If in Vim mode, apply filters
	case true:
		switch c.vimMode {
		case "insert":
			switch r {
			case readline.CharEsc:
				c.vimMode = "normal"
				_, m := c.prompt.render(true)
				c.Shell.SetPrompt(m)
				c.Shell.Refresh()
				return r, true

			case readline.CharCtrlL:
				readline.ClearScreen(c.Shell)
				c.Shell.Refresh()
				c.hardRefresh()
				return r, false
			}
		case "normal":
			switch r {
			case 'i', 'I', 'a', 'A', 's', 'S', 'c':
				c.vimMode = "insert"
				_, p := c.prompt.render(true)
				c.Shell.SetPrompt(p)
				c.Shell.Refresh()
				return r, true

			case readline.CharCtrlL:
				readline.ClearScreen(c.Shell)
				c.Shell.Refresh()
				c.hardRefresh()
				return r, true
			}
		}
		return r, true

	// If in Emacs, no filters needed
	case false:
		return r, true
	}

	return r, true
}

func (c *Console) exit() bool {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Confirm exit (Y/y): ")
	text, _ := reader.ReadString('\n')
	answer := strings.TrimSpace(text)

	if (answer == "Y") || (answer == "y") {
		c.Shell.Close()
		return true
	} else {
		return false
	}
}
