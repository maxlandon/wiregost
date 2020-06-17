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
	"github.com/jinzhu/gorm"

	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

// MigrateShema - Migrates all object definitions to Wiregost's PostgreSQL database.
func MigrateShema(db *gorm.DB) error {

	// Log level
	db.LogMode(true)

	// DB Options
	db.Set("gorm:auto_preload", true) // Always load relationships for an object.

	// Wiregost Users
	db.AutoMigrate(serverpb.User{})

	// User, Server & Implant Certificates
	db.AutoMigrate(serverpb.CertificateKeyPair{})

	// Workspaces

	// Hosts

	// IP Addresses

	// Ports & Services

	// Credentials

	// Routes

	// Listeners

	// Ghost Implant builds

	return nil
}
