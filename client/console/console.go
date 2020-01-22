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
	"os"

	"github.com/desertbit/grumble"
	"github.com/evilsocket/islazy/tui"
	"github.com/google/uuid"

	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/data_service/models"
	"github.com/maxlandon/wiregost/data_service/remote"
)

type Console struct {
	app    *grumble.App
	prompt Prompt

	// Environmment variables
	Env map[string]string

	// Context (temporary, check if we can get rid of this with grumble)
	context       context.Context
	menuContext   string
	currentModule string

	currentWorkspace   models.Workspace
	currentWorkspaceID *uint

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
	home, _ := os.UserHomeDir()

	console := &Console{}

	// Set console prompt
	prompt := newPrompt(console)

	// Set console app
	console.app = grumble.New(&grumble.Config{
		Name:            "Wiregost",
		Description:     tui.Blue(tui.Bold("Wiregost Client")),
		HistoryFile:     home + "/.wiregost/client/.history",
		HistoryLimit:    5000,
		Prompt:          prompt.render(),
		HelpSubCommands: true,
	})

	// Get default workspace
	workspaces, err := remote.Workspaces(nil)
	if err != nil {
		fmt.Println(tui.Red("Failed to fetch workspaces"))
	}
	for i, _ := range workspaces {
		if workspaces[i].IsDefault {
			console.currentWorkspace = workspaces[i]
			console.currentWorkspaceID = &workspaces[i].ID
			console.app.SetPrompt(prompt.render())
			commands.CurrentWorkspace.SetCurrentWorkspace(workspaces[i])
		}
	}
	// Set workspace refresher
	commands.CurrentWorkspace.AddObserver(func() {
		console.app.SetPrompt(prompt.render())
	})

	// Set console context
	rootCtx := context.Background()
	console.context = context.WithValue(rootCtx, "workspace_id", &commands.CurrentWorkspace.Workspace.ID)

	// Register Commands
	commands.RegisterCommands(&console.currentWorkspace, &console.context, console.app)

	// Start console
	console.Start()

	return console
}

func (c *Console) Start() error {

	// Run console
	c.app.Run()

	return nil
}
