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
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/lmorg/readline"

	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/commands/jobs"
	"github.com/maxlandon/wiregost/client/commands/sessions"
	"github.com/maxlandon/wiregost/client/completers"
	"github.com/maxlandon/wiregost/client/config"
	"github.com/maxlandon/wiregost/client/core"
	"github.com/maxlandon/wiregost/client/util"
	"github.com/maxlandon/wiregost/data-service/models"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
)

// Console is the client console object, and stores all client-side state.
type Console struct {
	// Shell
	Shell     *readline.Instance     // Console readline
	prompt    Prompt                 // Prompt string
	menu      string                 // Menu in which the shell is
	config    *config.Config         // Shell configuration
	server    *core.WiregostServer   // Wiregost Server
	workspace *models.Workspace      // Current workspace
	dbContext context.Context        // DB context
	module    *clientpb.Module       // Current module
	userID    int32                  // Unique user ID for module requests
	jobs      int                    // Number of jobs
	ghosts    int                    // Number of agents
	Ghost     *clientpb.Ghost        // Current implant
	GhostPwd  string                 // Current implant working directory
	context   *commands.ShellContext // Passes the shell state to commands
}

func newConsole() *Console {

	console := &Console{
		Shell:   readline.NewInstance(),
		menu:    commands.MAIN_CONTEXT,
		config:  config.LoadConsoleConfig(),
		userID:  rand.Int31(),
		module:  &clientpb.Module{}, // Needed even if empty
		Ghost:   &clientpb.Ghost{},  // Same
		context: &commands.ShellContext{},
	}

	return console
}

// Setup - Set all state for the shell
func (c *Console) Setup() {
	// Shell & Context
	c.initContext()

	// Completion, hints and syntax
	c.Shell.TabCompleter = completers.TabCompleter
	c.Shell.HintText = completers.HintText
	c.Shell.SyntaxHighlighter = completers.SyntaxHighlighter

	// Prompt
	c.prompt = newPrompt(c, c.config.Prompt, c.config.ImplantPrompt)

	// Commands
	RegisterCommands()

	// Env
	util.LoadSystemEnv()
}

// Start - Start the Shell
func Start() {

	// Instantiate and setup
	console := newConsole()
	console.Setup()

	// Connect to server
	err := console.connect(getDefaultServerConfig())
	if err != nil {
		log.Fatal(tui.Red(err.Error()))
	} else {
		fmt.Println()
		go console.eventLoop(console.server)
		commands.Context.Server = console.server
	}

	// Input loop
	for {
		// Refresh prompt
		console.hardRefresh()

		line, err := console.Shell.Readline()

		// To be deleted or modified if we use the flags library
		line = strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}

		sanitized, empty := splitAndSanitize(line)
		if empty {
			continue
		}

		// Leave a space between command and output
		fmt.Println()

		// Process tokens
		parsed, err := util.ParseEnvVariables(sanitized)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err = console.ExecCommand(parsed); err != nil {
			// For now we don't print errors here, because they appear twice because of parser
			// fmt.Printf(err.Error())
		}

	}
}

// // hardRefresh prints a new prompt
func (c *Console) hardRefresh() {
	// Menu context
	if len(c.module.Path) != 0 {
		c.menu = commands.MODULE_CONTEXT
	} else {
		c.menu = commands.MAIN_CONTEXT
	}
	if c.Ghost.Name != "" {
		c.menu = commands.GHOST_CONTEXT
	}

	// Jobs
	jobs := jobs.GetJobs(c.context.Server.RPC)
	c.jobs = len(jobs.Active)

	// Sessions
	sessions := sessions.GetGhosts(c.context.Server.RPC)
	c.ghosts = len(sessions.Ghosts)

	// Prompt
	refreshPrompt(c.prompt, c.Shell)
}

func (c *Console) exit() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Confirm exit (Y/y): ")
	text, _ := reader.ReadString('\n')
	answer := strings.TrimSpace(text)

	if (answer == "Y") || (answer == "y") {
		os.Exit(0)
	}
}

// splitAndSanitize - Various minor input sanitization steps
func splitAndSanitize(input string) (sanitized []string, empty bool) {

	// Assume the input is not empty
	empty = false

	// Trim last space
	line := strings.TrimSpace(input)
	if len(line) < 1 {
		empty = true
		return
	}

	unfiltered := strings.Split(line, " ")

	// Catch any eventual empty items
	for _, arg := range unfiltered {
		if arg != "" {
			sanitized = append(sanitized, arg)
		}
	}

	return
}
