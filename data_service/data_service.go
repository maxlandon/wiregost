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
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/evilsocket/islazy/tui"

	"github.com/maxlandon/wiregost/data_service/handlers"
)

func main() {
	// Setup DB and environment -------------------------------

	// Load DB credentials and data_service parameters
	env := handlers.LoadEnv()

	// AutoMigrate Schema
	err := env.DB.MigrateSchema()
	if err != nil {
		log.Printf("%s*%s Error: Could not migrate database schema: %s", tui.RED, tui.RESET, err.Error())

	}

	mux := http.NewServeMux()

	// Register handlers ---------------------------------------
	wh := &handlers.WorkspaceHandler{env}
	mux.Handle(handlers.WorkspaceAPIPath, wh)

	hh := &handlers.HostHandler{env}
	mux.Handle(handlers.HostAPIPath, hh)

	sh := &handlers.ServiceHandler{env}
	mux.Handle(handlers.ServiceAPIPath, sh)
	//
	// ch := &handlers.CredentialHandler{env}
	// mux.Handle(handlers.CredentialAPIPath, ch)
	//
	// Start server --------------------------------------------
	log.Printf("%s*%s Wiregost Data Service listening for requests...", tui.GREEN, tui.RESET)
	err = http.ListenAndServeTLS(env.Service.Address+":"+strconv.Itoa(env.Service.Port), env.Service.Certificate, env.Service.Key, mux)
	if err != nil {
		fmt.Println(err.Error())
	}
}
