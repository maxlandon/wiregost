package command

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

import (
	"github.com/reeflective/console"
	"github.com/spf13/cobra"

	"github.com/maxlandon/wiregost/cmd/client/command/credentials"
	"github.com/maxlandon/wiregost/cmd/client/command/exit"
	"github.com/maxlandon/wiregost/cmd/client/command/hosts"
	"github.com/maxlandon/wiregost/cmd/client/command/services"
	"github.com/maxlandon/wiregost/cmd/client/command/settings"
	client "github.com/maxlandon/wiregost/internal/client/console"
)

// Bind returns the entire wiregost client command
// tree for the closed loop console application.
func Bind(console *client.Client, adminCmds func(con *client.Client) []*cobra.Command) console.Commands {
	return func() *cobra.Command {
		rootCmd := &cobra.Command{}

		// Console specific
		rootCmd.AddCommand(exit.Command(console)...)

		// Core command tree.
		BindCLI(console, rootCmd)

		// Admin (server-only) commands.
		if adminCmds != nil {
			rootCmd.AddGroup(&cobra.Group{ID: "multiplayer", Title: "multiplayer"})
			rootCmd.AddCommand(adminCmds(console)...)
		}

		return rootCmd
	}
}

// BindCLI returns the entire wiregost client command
// tree, slightly modified for single-exec CLI usage.
func BindCLI(console *client.Client, rootCmd *cobra.Command) {
	// Core
	rootCmd.AddGroup(&cobra.Group{ID: "core", Title: "core"})
	rootCmd.SetCompletionCommandGroupID("core")
	rootCmd.SetHelpCommandGroupID("core")
	rootCmd.AddCommand(settings.Commands(console)...)

	// Database Objects
	rootCmd.AddGroup(&cobra.Group{ID: "database", Title: "database"})
	rootCmd.AddCommand(hosts.Commands())
	rootCmd.AddCommand(credentials.Commands())
	rootCmd.AddCommand(services.Commands())

	return
}
