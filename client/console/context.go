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

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/data_service/remote"
	"github.com/maxlandon/wiregost/server/module/templates"
)

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

	// Set context object (needed for some commands), further included in ShellContext
	rootCtx := context.Background()
	c.context = context.WithValue(rootCtx, "workspace_id", c.currentWorkspace.ID)

	// Set current module to nil
	c.currentModule = ""
	c.module = &templates.Module{}

	// Set ShellContext struct, passed to all commands
	c.shellContext = &commands.ShellContext{
		// Context
		Context:          c.context,
		MenuContext:      &c.menuContext,
		CurrentModule:    &c.currentModule,
		Module:           c.module,
		CurrentWorkspace: c.currentWorkspace,

		// Server state
		Server: c.server,

		// Jobs
		Listeners: &c.listeners,

		// Agents
		Ghosts: &c.ghosts,
		// Keep for prompt, until not needed anymore
		CurrentAgent: c.currentAgent,
	}
}
