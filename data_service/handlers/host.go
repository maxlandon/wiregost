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
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/maxlandon/wiregost/data_service/models"
)

const (
	// HostAPIPath is the API path to hosts
	HostAPIPath = "/api/v1/hosts/"
)

type HostHandler struct {
	*Env
}

func (hh *HostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Check if id is provided in URL, will influence dispatch
	id := strings.TrimPrefix(r.URL.Path, HostAPIPath)

	// Get workspace_id context in Header
	ws, _ := strconv.ParseUint(r.Header.Get("Workspace_id"), 10, 32)
	wsID := uint(ws)

	switch {
	// No ID, request applies to a Host range ------------------------//
	case id == "":
		switch {
		// Get some or all hosts from workspace
		case r.Method == "GET":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			// fmt.Println(len(b))

			switch {
			// No filter were provided
			case len(b) == 0:
				// fmt.Println(b)
				hosts, err := hh.DB.Hosts(wsID, nil)
				if err != nil {
					http.Error(w, http.StatusText(500), 500)
					return
				}

				json.NewEncoder(w).Encode(hosts)
				// Filters were provided, decode them
			default:
				// fmt.Println(len(b))
				var opts map[string]interface{}
				err = json.Unmarshal(b, &opts)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}

				hosts, err := hh.DB.Hosts(wsID, opts)
				if err != nil {
					http.Error(w, http.StatusText(500), 500)
					return
				}

				json.NewEncoder(w).Encode(hosts)
			}

		// Report a Host
		case r.Method == "POST":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			host, err := hh.DB.ReportHost(wsID)
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

	// ID, applies to a specific Host ---------------------------------//
	case id != "":
		// Delete a single Host
		switch {
		case r.Method == "DELETE":
			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var opts map[string]interface{}
			err = json.Unmarshal(b, &opts)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			deleted, err := hh.DB.DeleteHost(wsID, opts)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			if deleted != 1 {
				http.Error(w, "Some ids are not valid", 500)
			}

		// Update a Host
		case r.Method == "PUT":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

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
