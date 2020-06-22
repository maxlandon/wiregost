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

package db

import (
	"github.com/maxlandon/wiregost/db/server"
	"github.com/maxlandon/wiregost/server/assets"
)

// Setup - Establish connection with PostgreSQL, check DB exists and handles all along
func Setup() {

	// Configuration
	conf := assets.ServerConfiguration

	// Check required software is installed (PostgreSQL) and required access level.
	err := CheckPostgreSQLAccess()
	if err != nil {
		switch err.Error() {
		case ErrDatabaseDoesNotExist.Error(), ErrWiregostRoleDoesNotExist.Error():
			// Create DB
			err = InitDatabase()
			if err != nil {
				// We will most likely exit with a precise message on the cause
			}
		}
		// We will most likely exit with a precise message on the cause
	}

	// Connect to DB (no error checking, InitDatabase() already acts as a checker/automatic installer)
	server.ConnectPostgreSQL(conf.DBName, conf.DBUser, conf.DBPassword)

	// Load certificates/key pairs (stored in DB)

}

// Start - Starts one or more components of the Data Service
func Start() error {

	// Migrate Schema
	MigrateShema(server.DB)

	// Register & Start gRPC services (blocking)
	server.StartRPCServices()

	return nil
}

// StartRESTGateway - Start listening for REST requests
func StartRESTGateway() error {
	return nil
}

// StopRESTGateway - Stop the REST server
func StopRESTGateway() error {
	return nil
}
