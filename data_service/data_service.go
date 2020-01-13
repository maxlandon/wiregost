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
	"net/http"

	"github.com/maxlandon/wiregost/data_service/handlers"
	"github.com/maxlandon/wiregost/data_service/models"
)

func main() {
	// Setup DB and environment -------------------------------
	// Setup DB
	db := models.New("wiregost_db", "wiregost", "wiregost")

	// Setup env for passing DB connection around
	env := &handlers.Env{db}

	// Instantiate ServerMultiplexer
	mux := http.NewServeMux()

	// Register handlers ---------------------------------------
	wh := &handlers.WorkspaceHandler{env}
	mux.Handle(handlers.WorkspaceApiPath, wh)

	// Start server --------------------------------------------
	fmt.Println("Listening for requests...")
	http.ListenAndServe(":8000", mux)
}
