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

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/maxlandon/wiregost/data_service/models"
)

const (
	// HostAPIPath is the API path to hosts
	HostAPIPath = "/api/v1/hosts/"
)

type HostHandler struct {
	// Env is needed to pass a DB connection pool
	*Env
}

// ServeHTTP dispatches and process host requests
func (hh *HostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Check if id is provided in URL, will influence dispatch
	id := strings.TrimPrefix(r.URL.Path, HostAPIPath)

	switch {
	// ID is not there, applies to a range
	case id == "":
		switch {
		// Get all Hosts in database
		case r.Method == "GET":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			hosts, err := hh.DB.Hosts()
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}

			json.NewEncoder(w).Encode(hosts)
			return

		// Add a Host
		case r.Method == "POST":
			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var host *models.Host
			err = json.Unmarshal(b, &host)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			host, err = hh.DB.ReportHost(*host)
			fmt.Println(err.Error())
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			err = json.NewEncoder(w).Encode(host)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

		case r.Method == "DELETE":
		case r.Method == "PUT":
		}

	// ID is there, a host is specified
	case id != "":
		switch {
		case r.Method == "GET":
		case r.Method == "POST":

		// Delete a single Host
		case r.Method == "DELETE":
			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var id int
			err = json.Unmarshal(b, &id)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			deleted, err := hh.DB.DeleteHost(id)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			if deleted != 1 {
				http.Error(w, "Some ids are not valid", 500)
			}

		// Update a Host
		case r.Method == "PUT":
			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var host *models.Host
			err = json.Unmarshal(b, &host)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			host, err = hh.DB.UpdateHost(*host)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			err = json.NewEncoder(w).Encode(host)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	}

}
