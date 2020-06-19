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

package server

import (
	"crypto/x509"
	"fmt"

	"github.com/jinzhu/gorm"
)

// This file defines the GORM DB instance used to query Wiregost' PostgreSQL Database

// DB - The GORM DB instance
var DB *gorm.DB

// ConnectPostgreSQL - Connect to PostgreSQL
func ConnectPostgreSQL(name, user, password string) (db *gorm.DB, err error) {

	// Credentials
	creds := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", "localhost", 5432, user, name, password)
	DB, err = gorm.Open("postgres", creds)

	// Check for DB in ~/.wiregost filesystem, init it if not present
	if DatabaseNotSet() {
		InitDatabase()
	}

	// Connect with insecure, retrieve certificates for DB
	cert := ConnectInsecure()

	// Restart PostgreSQL with secure TLS credentials
	ConnectSecure(cert)

	return
}

// ConnectInsecure - Connect to DB with default or no encryption, retrieve certificates
func ConnectInsecure() (cert *x509.Certificate) {
	return
}

// ConnectSecure - restart PostgreSQL connection with good certs
func ConnectSecure(cert *x509.Certificate) {}

// DatabaseNotSet - Database does not exist yet
func DatabaseNotSet() (ok bool) {

	return
}

// InitDatabase - Initialize the PostgreSQL database, with default settings
func InitDatabase() error {

	return nil
}
