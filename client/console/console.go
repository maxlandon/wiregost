// +build linux
// Wiregost - Post-Exploitation & Implant Framework
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
	"fmt"
	"os"
	"strings"

	"github.com/maxlandon/readline"

	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/completers"
	"github.com/maxlandon/wiregost/client/connection"
	"github.com/maxlandon/wiregost/client/context"
	"github.com/maxlandon/wiregost/client/util"
	dbcli "github.com/maxlandon/wiregost/db/client"
	client "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
)

var (
	// Console - The client console object
	Console = newConsole()
)

// Console - Central object of the client UI
type console struct {
	Shell  *readline.Instance // Console readline input
	Config *client.ConsoleConfig
}

// newConsole - Instantiates a console with some default behavior
func newConsole() *console {

	console := &console{
		Shell:  readline.NewInstance(),
		Config: &client.ConsoleConfig{},
	}

	return console
}

// Connect - The console loads the server configuration, connects to it and attempts user authentication
func (c *console) Connect() (err error) {

	// Connect to server via TLS
	conn, err := connection.ConnectTLS()

	// Authenticate (5 tries)
	var cli client.ConnectionRPCClient
	cli, context.Context.Client = connection.Authenticate(conn)

	// Receive various infos sent by server when authenticated (ClientID, messages, users, version information, etc)
	info, config := context.GetConnectionInfo(cli)

	// Use console config received from the server
	c.Config = config

	// Connect to database on another connection
	dbcli.ConnectToDatabase("", int(info.DBPort), info.PublicKeyDB, info.PrivateKeyDB)

	// Print banner, user and client/server version information
	c.PrintBanner(context.GetVersion(cli), info)

	// Register all gRPC clients with the connection
	connection.RegisterRPCClients(conn)

	// Listen for incoming server/implant events
	go c.StartEventListener()

	return nil
}

// Setup - Setup various elements of the console.
func (c *console) Setup() (err error) {

	// Prompt
	c.InitPrompt()

	// Completion, Hints & Syntax
	c.Shell.TabCompleter = completers.TabCompleter
	c.Shell.HintText = completers.HintCompleter
	c.Shell.SyntaxHighlighter = completers.SyntaxHighlighter

	// Env
	err = util.LoadClientEnv()

	// Context

	// Commands
	err = commands.InitParsers()

	return
}

// Start - Start the console
func (c *console) Start() {

	// Connect to server and authenticate
	c.Connect()

	// Setup console
	c.Setup()

	// Input loop
	for {
		// Recompute prompt each time
		c.Refresh()

		// Readline
		line, _ := c.Readline()

		// Split & sanitize
		sanitized, empty := Sanitize(line)
		if empty {
			continue
		}

		// Process tokens
		parsed, _ := util.ParseEnvironmentVariables(sanitized)

		// Execute the command input
		c.ExecuteCommand(parsed)
	}
}

// Refresh - Computes prompt and current context
func (c *console) Refresh() {

	// Recompute prompt
	PromptBis.ComputePrompt()
}

// Readline - Add an empty line between input line and output
func (c *console) Readline() (line string, err error) {
	line, err = c.Shell.Readline()
	fmt.Println()
	return
}

// Sanitize - Trims spaces and other unwished elements from the input line
func Sanitize(line string) (sanitized []string, empty bool) {

	// Assume the input is not empty
	empty = false

	// Trim border spaces
	trimmed := strings.TrimSpace(line)
	if len(line) < 1 {
		empty = true
		return
	}

	unfiltered := strings.Split(trimmed, " ")

	// Catch any eventual empty items
	for _, arg := range unfiltered {
		if arg != "" {
			sanitized = append(sanitized, arg)
		}
	}

	return
}

// Exit - Kill the current client console
func (c *console) Exit() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Confirm exit (Y/y): ")
	text, _ := reader.ReadString('\n')
	answer := strings.TrimSpace(text)

	if (answer == "Y") || (answer == "y") {
		os.Exit(0)
	}

	fmt.Println()
}
