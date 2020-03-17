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

	"github.com/maxlandon/wiregost/data-service/models"
)

const (
	// ServiceAPIPath is the API path to services
	ServiceAPIPath = "/api/v1/services/"
)

// ServiceHandler handles all HTTP requests for querying services to the database
type ServiceHandler struct {
	*Env
}

// ServeHTTP processes requests, dispatch them and returns results
func (sh *ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Check if id is provided in URL, will influence dispatch
	path := strings.TrimPrefix(r.URL.Path, ServiceAPIPath)

	switch {
	// No path suffix, request applies to a Service range
	case path == "":
		switch {
		case r.Method == "GET":
			sh.services(w, r)

		case r.Method == "POST":
			sh.reportService(w, r)

		case r.Method == "DELETE":
			sh.deleteServices(w, r)
		}

	// Path is not nil, applies to a single Service
	case path != "":
		switch {
		case r.Method == "PUT":
			sh.updateService(w, r)
		}
	}
}

func (sh *ServiceHandler) services(w http.ResponseWriter, r *http.Request) {

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
		services, err := sh.DB.Services(wsID, nil)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		json.NewEncoder(w).Encode(services)
	// Filters were provided, decode them
	default:
		var opts map[string]interface{}
		err = json.Unmarshal(b, &opts)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		services, err := sh.DB.Services(wsID, opts)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		json.NewEncoder(w).Encode(services)
	}
}

func (sh *ServiceHandler) reportService(w http.ResponseWriter, r *http.Request) {

	h, _ := strconv.ParseUint(r.Header.Get("Host_id"), 10, 32)
	hostID := uint(h)

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

	service, err := sh.DB.ReportService(&hostID, opts)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(service)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (sh *ServiceHandler) deleteServices(w http.ResponseWriter, r *http.Request) {

	h, _ := strconv.ParseUint(r.Header.Get("Host_id"), 10, 32)
	hostID := uint(h)

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

	deleted, err := sh.DB.DeleteServices(hostID, opts)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if deleted != 1 {
		http.Error(w, "Some ids are not valid", 500)
	}
}

func (sh *ServiceHandler) updateService(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var service *models.Port
	err = json.Unmarshal(b, &service)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service, err = sh.DB.UpdateService(*service)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(service)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
