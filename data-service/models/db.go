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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evilsocket/islazy/tui"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/maxlandon/wiregost/data-service/assets"
)

// DB is the exported Database, against which are used entity methods
type DB struct {
	*gorm.DB
}

// New represents and gives access to a PostgreSQL database
func New(dbName string, user string, password string) (db *DB, err error) {
	creds := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", "localhost", 5432, user, dbName, password)
	gormDB, err := gorm.Open("postgres", creds)
	if err != nil {
		// switch err.Error() {
		// If role does not exist, that means the database doesn't either
		// case fmt.Sprintf("pq: role \"%s\" does not exist", user):
		err = InitDB(dbName, user, password)
		if err != nil {
			return nil, err
		} else {
			gormDB, err = gorm.Open("postgres", creds)
		}
		// }
	}
	db = &DB{gormDB}

	return db, nil
}

// InitDB - Creates Database user, password and Database
func InitDB(dbName, user, password string) error {

	// Create temporary SQL file (os/exec is just a mess when passing psql commands)
	saveTo := assets.GetDataServiceDir()

	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		err = os.MkdirAll(saveTo, os.ModePerm)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot write to wiregost root directory %s", err))
		}
	}

	fi, err := os.Stat(saveTo)
	if fi.IsDir() {
		filename := "default_db.sql"
		saveTo = filepath.Join(saveTo, filename)
	}

	f, err := os.OpenFile(saveTo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	if _, err = f.WriteString(sqlQueries); err != nil {
		panic(err)
	}

	// Create Database
	cmd := exec.Command("psql", "-U", "postgres", "-f", saveTo)
	_, err = cmd.Output()
	if err != nil {
		return err
	} else {
		fmt.Printf("%s[*]%s Successfully created Database %s with user %s and password %s",
			tui.GREEN, tui.RESET, dbName, user, password)
	}

	// Delete SQL queries file
	f.Close()
	err = os.Remove(saveTo)
	if err != nil {
		fmt.Println(tui.Red("Failed to delete temporary SQL file 'default_db.sql'"))
	}

	return nil
}

const sqlQueries = `CREATE DATABASE wiregost_db;
CREATE USER wiregost;
ALTER ROLE wiregost WITH PASSWORD 'wiregost';
GRANT ALL ON DATABASE wiregost_db TO wiregost;
`
