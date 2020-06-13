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
	"fmt"

	"github.com/google/uuid"
	"github.com/lmorg/readline"

	"github.com/maxlandon/wiregost/client/assets"
	// "github.com/maxlandon/wiregost/client/core"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

var (
	// Console - The client console object
	Console = newConsole()
)

// Console - Central object of the client UI
type console struct {
	ClientID uuid.UUID             // Unique identifier for this console
	User     *serverpb.User        // User information sent back after auth
	Shell    *readline.Instance    // Console readline input
	Config   *assets.ConsoleConfig // Console configuration
	// Server   *core.WiregostServer  // Server connection infrastructure
	Module *modulepb.Module // Module currently on stack
	Ghost  *ghostpb.Ghost   // Current ghost implant
	Ghosts int
	Jobs   int
}

// newConsole - Instantiates a console with some default behavior
func newConsole() *console {

	console := &console{
		Shell:  readline.NewInstance(),
		Module: &modulepb.Module{},
		Ghost:  &ghostpb.Ghost{},
	}

	return console
}

// Connect - The console loads the server configuration, connects to it and atempts user authentication
func (c *console) Connect() {

	// Load server connection configuration (check files in ~/.wiregost first, then binary)

	// Connect to server

	// Authenticate (5 tries)

	// Receive various infos sent by server when authenticated (ClientID, messages, users, etc)

	// Setup
	c.Setup()
}

// Setup - Setup various elements of the console.
func (c *console) Setup() {

	// Console configuration (from server first, ~/.wiregost second)

	// Prompt

	// Completion, Hints & Syntax

	// Commands

	// Env

	// Share context
}

// ShareContext - The console exposes its context to other packages
func (c *console) ShareContext() {

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
		// line, _ := c.Readline()

		// Split & sanitize
		// sanitized, empty := Sanitize(line)

		// Process tokens

		// Execute the command input
	}
}

// Refresh - Computes prompt and current context
func (c *console) Refresh() {

}

// Readline - Add an empty line between input line and output
func (c *console) Readline() (line string, err error) {
	line, err = c.Shell.Readline()
	fmt.Println()
	return
}

// Exit - Kill the current client console
func (c *console) Exit() {

}

// Sanitize - Trims spaces and other unwished elements from the input line
func Sanitize(line string) (sanitized []string, empty bool) {

	// Trim border spaces

	// Catch eventual empty items

	return
}
