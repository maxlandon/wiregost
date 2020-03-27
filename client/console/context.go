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

	c.context = &commands.ShellContext{
		Shell:     c.Shell,
		Config:    c.config,
		DBContext: c.dbContext,
		Menu:      &c.menu,
		Module:    c.module,
		UserID:    c.userID,
		Workspace: c.workspace,
		Server:    c.server,
		Jobs:      &c.jobs,
		Ghosts:    &c.ghosts,
		Ghost:     c.Ghost,
		GhostPwd:  &c.GhostPwd,
	}
}
