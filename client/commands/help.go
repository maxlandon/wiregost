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

package commands

import (
	"fmt"

	"github.com/desertbit/grumble"
	"github.com/evilsocket/islazy/tui"

	consts "github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/help"
)

func RegisterHelpCommands(app *grumble.App) {

	// Command categories
	helpCategoriesCommand := &grumble.Command{
		Name: "help",
		Help: tui.Blue(tui.Bold("  Command Categories")),
	}

	// Workspace Commands
	helpCategoriesCommand.AddCommand(&grumble.Command{
		Name: "workspace",
		Help: tui.Dim("Manage Wiregost workspaces"),
		Run: func(ctx *grumble.Context) error {
			fmt.Println()
			fmt.Println(help.GetHelpFor(consts.WorkspaceStr))
			fmt.Println()
			return nil
		},
		HelpGroup: consts.DataServiceHelpGroup,
	})
	// Hosts Commands
	helpCategoriesCommand.AddCommand(&grumble.Command{
		Name: "hosts",
		Help: tui.Dim("Manage database hosts"),
		Run: func(ctx *grumble.Context) error {
			fmt.Println()
			fmt.Println(help.GetHelpFor(consts.HostsStr))
			fmt.Println()
			return nil
		},
		HelpGroup: consts.DataServiceHelpGroup,
	})

	// Finally register commands
	app.AddCommand(helpCategoriesCommand)
}
