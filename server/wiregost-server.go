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
	"github.com/maxlandon/wiregost/server/certs"
	"github.com/maxlandon/wiregost/server/clients"
	"github.com/maxlandon/wiregost/server/events"
)

// The Server executable is the main piece of the Wiregost system. It is also the one that should NEVER have to be restarted.
// This is a challenge to the extent that this executable has to start and coordinate many pieces of information, across
// many network and RPC layers, while sharing state with many components, without any recursive import.
// For instance: many users are writing many modules, using many sessions. If any user has to modify any bit of code into
// the full Wiregost system, which is fully COMPILED, we have to comply with Go's apparent mantra: micro-service architecture.

func main() {

	// CONFIGURATION
	// The server makes configuration available to all components of the framework.
	// This package never imports other packages from under the server/ directory.
	// This function will also be called by other executables, such as the module manager
	assets.LoadServerConfiguration()

	// Check assets presence and unpacking process. Includes various toolchains, modules code and file structure, dictionaries, etc.
	// This package is only called by the server, because the module_manager isn't in charge of its own compilation.
	// It should never call any other package under server/, as many packages will import it, such as log/
	assets.SetupAssets()

	// LOGGING
	// The logging infrastructure plays a central role: because the logrus library being used in Wiregost supports various things.
	// In particular, the logger is used in different ways by different components:
	// The server uses it mainly for logging code workflow for components and logging (lots of) user commands.
	// Implants use it for logging their workflow output output to files and to consoles if needed.
	// Logging infrastructure could be leveraged to provide a more fine-grained event provider system, or at leat embbed it.
	// But this raises the question of how many packages should be imported by the logging infrastructure, if it were to
	// be fully context-aware.

	// CERTIFICATE INFRASTRUCTURE
	// The certificate infrastructure is taken and copied from Sliver's work. It is a blessing they did this: Wiregost aims to
	// provide wide networking capabilities with various implants and host systems, as well as various server-side systems
	// talking to each other. This will add a bit of security.
	// Load certificates management. Only the server needs to verify things related to
	// transport security and authentication problems, as it will always be up and running.
	certs.SetupCertificateAuthorities()

	// DATABASE
	// The database system in Wiregost is working both similar and different from the module manager: the database is not a
	// standalone binary (it works within this server process) but it listens on a port for connections from user consoles,
	// who may want to query various informations, either for printing or for completion, between others.
	// Therefore, all database server-side code is setup and instantiated here. This includes auto-migration of all data models
	// Wiregost makes use of.
	// AutoMigrate & Setup Database
	db.Setup()

	// Start Database service (PostgreSQL backend + gRPC server)
	go db.Start()

	// Setup client connection to DB (the server is itself a client of the DB)
	dbcli.ConnectServerToDB()

	// EVENTS
	// Wiregost handles various kinds of events pushed by components: jobs/proxies termination, client connections, impant
	// registrations, module events, etc... We setup and register all event subscribers here, available for all packages.
	// The event manager should offer a gRPC server to consoles (for pushing them events) and module manager (for pushing
	// and receiving events). All events happening in Wiregost always go through this package for processing and dispatch.
	go events.Broker.Start()

	// MODULE SYSTEM
	// The module system is composed of a module manager standalone program, which holds all available modules in Wiregost.
	// It is somehow a "live stack", that communicates over gRPC with the server and the database, either for requiring
	// implant actions, for pushing content to user consoles, etc.
	// It is handled and controlled by this server, which can stop, restart and recompile a module manager. It starts and
	// communicates with one module manager for each connected user, so that each of them can use, write and modify modules
	// without bothering the others.
	// modules.StartManagers()

	// PERSISTENCE
	// We might have some persistence needs, such as automatic listeners with
	// various preset rules (routes to open, pivots to reach, etc...)

	// Start Listening for client consoles
	clients.StartClientListener(assets.ServerConfiguration.ServerHost, assets.ServerConfiguration.ServerPort)
}
