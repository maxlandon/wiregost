package cli

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

	"github.com/maxlandon/wiregost/cmd/client/command"
	"github.com/maxlandon/wiregost/cmd/server/command/operator"
	"github.com/maxlandon/wiregost/internal/client/console"
	"github.com/maxlandon/wiregost/internal/server/transport"
)

// initConsoleCmd sets up the console runner with an in-memory listener for the server.
func initConsoleCmd(con *console.Client) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		rpc, conn, err := transport.ConnectClient(nil)
		if err != nil {
			return err
		}

		defer conn.Close()

		// Bind the RPC connection, and all commands including admin ones.
		console.Setup(con, nil, rpc, command.Bind(con, operator.ServerCommands))

		return con.Start()
	}
}
