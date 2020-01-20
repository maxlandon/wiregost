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

package models

import (
	"fmt"
	"log"

	"github.com/evilsocket/islazy/tui"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB is the exported Database, against which are used entity methods
type DB struct {
	*gorm.DB
}

// New represents and gives access to a PostgreSQL database
func New(dbName string, user string, password string) *DB {
	creds := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", "localhost", 5432, user, dbName, password)
	gormDB, err := gorm.Open("postgres", creds)
	if err != nil {
		log.Fatal(tui.Red("Could not connect to PostgreSQL with provided parameters and credentials" + err.Error()))
	}
	db := &DB{gormDB}

	return db
}
