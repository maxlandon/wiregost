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
	"github.com/maxlandon/wiregost/db/models"
	"github.com/maxlandon/wiregost/db/server"
	"github.com/maxlandon/wiregost/server/assets"
)

// Start - Starts one or more components of the Data Service
func Start() error {

	// Connect to DB
	conf := assets.ServerConfiguration

	// db, _ := models.ConnectDatabase(conf.DBName, conf.DBUser, conf.DBPassword)
	_, _ = models.ConnectDatabase(conf.DBName, conf.DBUser, conf.DBPassword)

	// Load certificates/key pairs (stored in DB)

	// Migrate Schema
	// MigrateShema(db)

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
