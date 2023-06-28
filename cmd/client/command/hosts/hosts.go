package hosts

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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/maxlandon/aims/display"
	"github.com/maxlandon/aims/host"
)

// Commands returns a command tree to manage and display hosts.
func Commands() *cobra.Command {
	hostsCmd := &cobra.Command{
		Use:     "hosts",
		Short:   "Manage database hosts",
		GroupID: "database",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Display hosts (with filters or styles)",
		RunE: func(cmd *cobra.Command, args []string) error {
			h1 := host.Host{OSName: "Windows", OSFamily: "Vista", MAC: "18:32:ds:1d:12:c4", Arch: "x64", Purpose: "server"}
			h2 := host.Host{OSName: "Linux", OSFamily: "Fedora", MAC: "23:cc:eh:1d:23:i2", Arch: "x86", Purpose: "router"}

			table := display.Table([]display.Rower{&h1, &h2})

			fmt.Println(table.Render())

			return nil
		},
	}

	hostsCmd.AddCommand(listCmd)

	rmCmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove one or more hosts from the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	hostsCmd.AddCommand(rmCmd)

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show one ore more hosts details",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	hostsCmd.AddCommand(showCmd)

	return hostsCmd
}
