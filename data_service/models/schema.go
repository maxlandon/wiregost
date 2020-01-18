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
	"log"

	"github.com/evilsocket/islazy/tui"
	"github.com/go-pg/pg/orm"
)

func (db *DB) CreateSchema() error {
	// Defining/importing models
	tableModels := []interface{}{
		&Workspace{},
		// &User{},

		// Hosts
		&Host{},

		// Services
		&Service{},

		// Web Clients
		&Client{},

		// Realms
		&Realm{},

		// Credentials
		&PublicCredential{},
		&PrivateCredential{},
		&Credential{},

		// Logins
		&Login{},

		// Tasks
		&Task{},
		&TaskHost{},
		&TaskService{},
		&TaskCredential{},
		&TaskAgent{},

		// Origins
		&OriginCrackedPassword{},
		&OriginImport{},
		&OriginManual{},
		&OriginService{},
		&OriginAgent{},

		// Listeners
		&Listener{},
	}

	// Define table options
	tableOptions := &orm.CreateTableOptions{
		FKConstraints: true,
	}

	// Create tables for each model
	for _, model := range tableModels {
		err := db.CreateTable(model, tableOptions)
		if err != nil {
			log.Println(err.Error())
		}
	}

	// Insert default fields/items if needed
	db.AddWorkspaces([]string{"default"})
	workspaces, _ := db.Workspaces()
	workspaces[0].IsDefault = true
	db.UpdateWorkspace(*workspaces[0])

	// Check tables

	return nil
}

func (db *DB) SchemaIsUpdated() bool {
	// Use info as generic struct, used successively in this function
	var info []struct {
		ColumnName string
		DataType   string
	}

	// Check for workspaces table
	_, err := db.Query(&info, `
                        SELECT column_name, data_type
                        FROM information_schema.columns
                        WHERE table_name = 'workspaces'
                `)
	if err != nil {
		log.Printf("%s[!] Error: could not query DB for testing schema: %s\n", tui.RED, err.Error())
		return false
	}

	// Break if info is empty
	if len(info) == 0 {
		return false
	}

	// Check for hosts table
	_, err = db.Query(&info, `
                        SELECT column_name, data_type
                        FROM information_schema.columns
                        WHERE table_name = 'hosts'
                `)
	if len(info) == 0 {
		return false
	}

	// Check for credentials table
	_, err = db.Query(&info, `
                        SELECT column_name, data_type
                        FROM information_schema.columns
                        WHERE table_name = 'creds'
                `)
	if len(info) == 0 {
		return false
	}

	return true
}
