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

// HostHandler handles all HTTP requests for querying Hosts to the database
type HostHandler struct {
	*Env
}

// ServeHTTP processes requests, dispatch them and returns results
func (hh *HostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Check if id is provided in URL, will influence dispatch
	path := strings.TrimPrefix(r.URL.Path, HostAPIPath)

	switch {
	// No path suffix, request applies to a Host range ---------------------//
	case path == "":
		switch {
		// Get some or all hosts from workspace
		case r.Method == "GET":
			hh.hosts(w, r)

		// Report a Host
		case r.Method == "POST":
			hh.reportHost(w, r)

		// Delete a host
		case r.Method == "DELETE":
			hh.deleteHost(w, r)
		}

	// Path is not nil, applies to a single host ---------------------------//
	case path != "":
		switch {
		// Path is a search: applies to a single, non-identified host
		case path == "search":
			switch {
			case r.Method == "POST":
				hh.getHost(w, r)
			}

		// Path is a Host ID. Applies to a single host.
		default:
			// Delete a single Host
			switch {
			case r.Method == "DELETE":
				hh.deleteHost(w, r)

			// Update a Host
			case r.Method == "PUT":
				hh.updateHost(w, r)
			}
		}
	}
}

func (hh *HostHandler) hosts(w http.ResponseWriter, r *http.Request) {

	// Get workspace_id context in Header
	ws, _ := strconv.ParseUint(r.Header.Get("Workspace_id"), 10, 32)
	wsID := uint(ws)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	switch {
	// No filter were provided
	case len(b) == 0:
		hosts, err := hh.DB.Hosts(wsID, nil)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		json.NewEncoder(w).Encode(hosts)
		// Filters were provided, decode them
	default:
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
}

func (hh *HostHandler) reportHost(w http.ResponseWriter, r *http.Request) {

	// Get workspace_id context in Header
	ws, _ := strconv.ParseUint(r.Header.Get("Workspace_id"), 10, 32)
	wsID := uint(ws)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

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

	host, err := hh.DB.ReportHost(wsID, opts)
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

func (hh *HostHandler) deleteHosts(w http.ResponseWriter, r *http.Request) {

	// Get workspace_id context in Header
	ws, _ := strconv.ParseUint(r.Header.Get("Workspace_id"), 10, 32)
	wsID := uint(ws)

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
}

func (hh *HostHandler) getHost(w http.ResponseWriter, r *http.Request) {

	// Get workspace_id context in Header
	ws, _ := strconv.ParseUint(r.Header.Get("Workspace_id"), 10, 32)
	wsID := uint(ws)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

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

	host, err := hh.DB.GetHost(wsID, opts)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(host)
}

func (hh *HostHandler) updateHost(w http.ResponseWriter, r *http.Request) {

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

func (hh *HostHandler) deleteHost(w http.ResponseWriter, r *http.Request) {

	// Get workspace_id context in Header
	ws, _ := strconv.ParseUint(r.Header.Get("Workspace_id"), 10, 32)
	wsID := uint(ws)

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
}
