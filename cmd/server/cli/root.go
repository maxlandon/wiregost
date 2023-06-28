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
	"path/filepath"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"

	"github.com/maxlandon/wiregost/cmd/client/cli"
	"github.com/maxlandon/wiregost/internal/client/console"
	"github.com/maxlandon/wiregost/internal/server/certs"
	"github.com/maxlandon/wiregost/internal/server/configs"
	"github.com/maxlandon/wiregost/internal/server/daemon"
	wglog "github.com/maxlandon/wiregost/internal/server/log"
	"github.com/maxlandon/wiregost/internal/server/multiplayer"
)

const (
	// Unpack flags
	forceFlagStr = "force"

	// Operator flags
	nameFlagStr  = "name"
	lhostFlagStr = "lhost"
	lportFlagStr = "lport"
	saveFlagStr  = "save"

	// Cert flags
	caTypeFlagStr = "type"
	loadFlagStr   = "load"

	// console log file name
	logFileName = "console.log"
)

func init() {
	// Create the console client, without any RPC or commands bound to it yet.
	// This created before anything so that multiple commands can make use of
	// the same underlying command/run infrastructure.
	console := console.NewClient(true)
	console.IsCLI = true

	// By default, all commands require the console to be set up.
	rootCmd.PersistentPreRunE, rootCmd.PersistentPostRunE = setupServerPreRunners(console)

	// Console
	consoleCmd := cli.ConsoleCmd(console, true)
	consoleCmd.RunE = initConsoleCmd(console)
	rootCmd.AddCommand(consoleCmd)

	// Operator
	operatorCmd.Flags().StringP(nameFlagStr, "n", "", "operator name")
	operatorCmd.Flags().StringP(lhostFlagStr, "l", "", "multiplayer listener host")
	operatorCmd.Flags().Uint16P(lportFlagStr, "p", uint16(31337), "multiplayer listener port")
	operatorCmd.Flags().StringP(saveFlagStr, "s", "", "save file to ...")
	rootCmd.AddCommand(operatorCmd)

	// Certs
	cmdExportCA.Flags().StringP(saveFlagStr, "s", "", "save CA to file ...")
	rootCmd.AddCommand(cmdExportCA)

	cmdImportCA.Flags().StringP(loadFlagStr, "l", "", "load CA from file ...")
	rootCmd.AddCommand(cmdImportCA)

	// Daemon
	daemonCmd.Flags().StringP(lhostFlagStr, "l", daemon.BlankHost, "multiplayer listener host")
	daemonCmd.Flags().Uint16P(lportFlagStr, "p", daemon.BlankPort, "multiplayer listener port")
	rootCmd.AddCommand(daemonCmd)

	// Completions
	carapace.Gen(rootCmd)
}

// Initialize logging
func initConsoleLogging(appDir string) *os.File {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logFile, err := os.OpenFile(filepath.Join(appDir, logFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	log.SetOutput(logFile)
	return logFile
}

var rootCmd = &cobra.Command{
	Use:   "wiregost-server",
	Short: "Wiregost exploitation toolkit (server)",
}

func setupServerPreRunners(con *console.Client) (pre, post func(cmd *cobra.Command, args []string) error) {
	var logFile *os.File

	pre = func(cmd *cobra.Command, args []string) error {
		// Setup logging
		appDir := wglog.GetLogDir()
		logFile = initConsoleLogging(appDir)

		// Setup CAs, database schema
		certs.SetupCAs()

		// Start persistent server jobs
		serverConfig := configs.GetServerConfig()
		multiplayer.StartPersistentJobs(serverConfig)

		return nil
	}

	post = func(cmd *cobra.Command, args []string) error {
		return logFile.Close()
	}

	return pre, post
}

// Execute runs the wiregost server root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
