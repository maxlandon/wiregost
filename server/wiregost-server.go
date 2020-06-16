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
	"github.com/maxlandon/wiregost/db"
	dbcli "github.com/maxlandon/wiregost/db/client"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/modules"
	"github.com/maxlandon/wiregost/server/rpc"
)

func main() {

	// Load configuration
	assets.LoadServerConfiguration()

	// Check assets presence/unpacking
	assets.SetupAssets()

	// Setup logging

	// AutoMigrate & Setup Database, Start & Test Connection
	go db.Start()

	// Setup client connection to DB (the server is itself a client of the DB)
	dbcli.ConnectServerToDB()

	// Load certificates

	// Load modules
	modules.RegisterModules()

	// Init users module stacks
	modules.InitStacks()

	// Start Persistent implants

	// Start Listening for client consoles
	rpc.StartClientListener(assets.ServerConfiguration.ServerHost, assets.ServerConfiguration.ServerPort)
}
