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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/evilsocket/islazy/tui"

	"github.com/maxlandon/wiregost/client/version"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/certs"
	"github.com/maxlandon/wiregost/server/config"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/module/load"
	"github.com/maxlandon/wiregost/server/transport"
	"github.com/maxlandon/wiregost/server/users"
)

var (
	wiregostServerVersion = fmt.Sprintf("Client v%s\nServer v0.0.7", version.FullVersion())
)

const (
	logFileName = "console.log"
)

func main() {
	unpack := flag.Bool("unpack", false, "force unpack assets")
	flag.Parse()

	// Set Logging
	appDir := assets.GetRootAppDir()
	logFile := initLogging(appDir)
	defer logFile.Close()

	// Setup Certificate Infrastructure
	certs.SetupCAs()

	// Setup static assets
	assets.Setup(*unpack)
	if *unpack {
		os.Exit(0)
	}

	// Load server config
	servConf := config.LoadServerConfig()

	// Check at least one user exists
	err := users.CreateDefaultUser()
	if err != nil {
		log.Println(err.Error())
	}

	// Initialize Module Stacks
	load.LoadModules()
	core.InitStacks()

	// Start client listener
	listener, err := transport.StartClientListener(servConf.LHost, uint16(servConf.LPort))
	if err != nil {
		log.Printf("%s Failed to start MTLS client listener: %s", tui.RED, err.Error())
	} else {
		fmt.Printf("%s Started MTLS client listener at %s %s", tui.GREEN, listener.Addr(), tui.RESET)
		// log.Printf("%s Started MTLS client listener at %s", tui.GREEN, listener.Addr())
	}

	// Block and listen for connections
	transport.AcceptClientConnections(listener)
}

// Initialize logging
func initLogging(appDir string) *os.File {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logFile, err := os.OpenFile(path.Join(appDir, logFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	log.SetOutput(logFile)
	return logFile
}
