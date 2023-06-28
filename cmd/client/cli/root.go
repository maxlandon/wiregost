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
	"log"
	"os"
	"path"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"

	"github.com/maxlandon/wiregost/cmd/client/command"
	"github.com/maxlandon/wiregost/internal/client/console"
	"github.com/maxlandon/wiregost/internal/client/version"
)

const (
	logFileName = "wiregost-client.log"
)

var wiregostServerVersion = fmt.Sprintf("v%s", version.FullVersion())

func init() {
	// Create the console client, without any RPC or commands bound to it yet.
	// This created before anything so that multiple commands can make use of
	// the same underlying command/run infrastructure.
	console := console.NewClient(false)

	// By default, all commands require the console client to be set up, even if not ran.
	rootCmd.PersistentPreRunE, rootCmd.PersistentPostRunE = setupConsole(console, false)

	// Configurations
	rootCmd.AddCommand(importCmd())

	// Closed-loop console start command
	rootCmd.AddCommand(ConsoleCmd(console, false))

	// All commands
	command.BindCLI(console, rootCmd)

	// Completions
	carapace.Gen(rootCmd)
}

var rootCmd = &cobra.Command{
	Use:   "wiregost-client",
	Short: "Wiregost exploitation toolkit (client)",
}

// Initialize logging
func initLogging(appDir string) *os.File {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logFile, err := os.OpenFile(path.Join(appDir, logFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		panic(fmt.Sprintf("[!] Error opening file: %s", err))
	}
	log.SetOutput(logFile)
	return logFile
}

// Execute runs the wiregost client root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
