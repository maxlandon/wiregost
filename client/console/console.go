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
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/google/uuid"

	"github.com/maxlandon/wiregost/client/assets"
	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/completers"
	"github.com/maxlandon/wiregost/data_service/models"
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

	// Server state
	currentServer  *assets.ClientConfig
	serverPublicIP string

	// Jobs
	listeners int

	// Agents
	ghosts int
	// Keep for prompt, until not needed anymore
	currentAgentID uuid.UUID

	// CommandShellContext
	shellContext *commands.ShellContext
}

func NewConsole() *Console {

	shell, _ := readline.NewEx(&readline.Config{
		HistoryFile:       assets.GetRootAppDir() + "/.history",
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

	return console
}

func Start() {

	// Instantiate console
	c := NewConsole()

	// Connect to server
	c.connect()

	// Eventually close
	defer c.Shell.Close()

	// Command loop
	for {
		// Refresh Vim mode each time is needed here
		c.vimMode = "insert"
		c.refresh()

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

		var args []string
		for _, arg := range unfiltered {
			if arg != "" {
				args = append(args, arg)
			}
		}

		if err = ExecCmd(args, c.menuContext, c.shellContext); err != nil {
			fmt.Println(err)
		}
	}
}

// [ Generic console functions ] ---------------------------------------------------------------//

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
