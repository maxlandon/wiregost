package settings

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
	"github.com/reeflective/console/commands/readline"
	"github.com/spf13/cobra"

	"github.com/maxlandon/wiregost/internal/client/console"
)

// Commands returns the “ command and its subcommands.
func Commands(con *console.Client) []*cobra.Command {
	if con.IsCLI {
		return nil
	}

	settingsCmd := &cobra.Command{
		Use:   "settings",
		Short: "Manage client settings",
		Run: func(cmd *cobra.Command, args []string) {
			SettingsCmd(cmd, con, args)
		},
		GroupID: "core",
	}
	settingsCmd.AddCommand(&cobra.Command{
		Use:   "save",
		Short: "Save the current settings to disk",
		Run: func(cmd *cobra.Command, args []string) {
			SettingsSaveCmd(cmd, con, args)
		},
	})
	settingsCmd.AddCommand(&cobra.Command{
		Use:   "tables",
		Short: "Modify tables setting (style)",
		Run: func(cmd *cobra.Command, args []string) {
			SettingsTablesCmd(cmd, con, args)
		},
	})
	settingsCmd.AddCommand(&cobra.Command{
		Use:   "always-overflow",
		Short: "Disable table pagination",
		Run: func(cmd *cobra.Command, args []string) {
			SettingsAlwaysOverflow(cmd, con, args)
		},
	})
	settingsCmd.AddCommand(&cobra.Command{
		Use:   "small-terminal",
		Short: "Set the small terminal width",
		Run: func(cmd *cobra.Command, args []string) {
			SettingsSmallTerm(cmd, con, args)
		},
	})
	settingsCmd.AddCommand(&cobra.Command{
		Use:   "console-logs",
		Short: "Log console output (toggle)",
		Run: func(ctx *cobra.Command, args []string) {
			SettingsConsoleLogs(ctx, con)
		},
	})

	// Bind a readline subcommand to the `settings` one, for allowing users to
	// manipulate the shell instance keymaps, bindings, macros and global options.
	settingsCmd.AddCommand(readline.Commands(con.App.Shell()))

	return []*cobra.Command{settingsCmd}
}
