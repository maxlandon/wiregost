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
	"fmt"
	"os"

	"github.com/maxlandon/wiregost/cmd/client/command"
	"github.com/maxlandon/wiregost/internal/client/assets"
	"github.com/maxlandon/wiregost/internal/client/console"
	"github.com/maxlandon/wiregost/internal/client/transport"
	"github.com/maxlandon/wiregost/internal/proto/rpcpb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// ConsoleCmd returns a command that starts the console application.
func ConsoleCmd(con *console.Client, isServer bool) *cobra.Command {
	consoleCmd := &cobra.Command{
		Use:   "console",
		Short: "Start the closed-loop console application",
		RunE: func(cmd *cobra.Command, args []string) error {
			return con.Start()
		},
	}

	if !isServer {
		consoleCmd.GroupID = "core"
	}

	return consoleCmd
}

func setupConsole(con *console.Client, runLoop bool) (pre, post func(cmd *cobra.Command, args []string) error) {
	var ln *grpc.ClientConn
	var logFile *os.File

	pre = func(_ *cobra.Command, _ []string) error {
		appDir := assets.GetRootAppDir()
		logFile = initLogging(appDir)

		configs := assets.GetConfigs()
		if len(configs) == 0 {
			fmt.Printf("No config files found at %s (see --help)\n", assets.GetConfigDir())
			return nil
		}
		config := selectConfig()
		if config == nil {
			return nil
		}

		// Don't clobber output when simply running an implant command from system shell.
		if runLoop {
			fmt.Printf("Connecting to %s:%d ...\n", config.LHost, config.LPort)
		}

		var rpc rpcpb.CoreClient
		var err error

		rpc, ln, err = transport.ConnectClient(config)
		if err != nil {
			fmt.Printf("Connection to server failed %s", err)
			return nil
		}

		// Bind the RPC connection and commands.
		console.Setup(con, config, rpc, command.Bind(con, nil))

		return nil
	}

	// Close the RPC connection once exiting
	post = func(_ *cobra.Command, _ []string) error {
		if ln != nil {
			ln.Close()
		}

		if logFile != nil {
			logFile.Close()
		}

		return nil
	}

	return pre, post
}
