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
	"github.com/maxlandon/wiregost/data-service/remote"
)

// initContext - Shares the console state to commands
func (c *Console) initContext() {

	// Console Context --------------------------------------------------

	// Workspace
	workspaces, err := remote.Workspaces(nil)
	if err != nil {
		fmt.Println(tui.Red("Failed to fetch workspaces"))
	}
	for i := range workspaces {
		if workspaces[i].IsDefault {
			c.workspace = &workspaces[i]
		}
	}

	// Data Service
	rootCtx := context.Background()
	c.dbContext = context.WithValue(rootCtx, "workspace_id", c.workspace.ID)

	// Share Context --------------------------------------------------
	commands.Context.Shell = c.Shell
	commands.Context.Config = c.config
	commands.Context.DBContext = c.dbContext
	commands.Context.Menu = &c.menu
	commands.Context.Module = c.module
	commands.Context.UserID = c.userID
	commands.Context.Workspace = c.workspace
	commands.Context.Server = c.server
	commands.Context.Jobs = &c.jobs
	commands.Context.Ghosts = &c.ghosts
	commands.Context.Ghost = c.Ghost
	commands.Context.GhostPwd = &c.GhostPwd
}
