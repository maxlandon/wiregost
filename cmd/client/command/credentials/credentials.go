package credentials

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
	"github.com/spf13/cobra"
)

// Commands returns a command tree to manage and display credentials.
func Commands() *cobra.Command {
	credentialsCmd := &cobra.Command{
		Use:     "credentials",
		Short:   "Manage database credentials",
		GroupID: "database",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Display credentials (with filters or styles)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	credentialsCmd.AddCommand(listCmd)

	rmCmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove one or more credentials from the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	credentialsCmd.AddCommand(rmCmd)

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show one ore more credentials details",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	credentialsCmd.AddCommand(showCmd)

	return credentialsCmd
}
