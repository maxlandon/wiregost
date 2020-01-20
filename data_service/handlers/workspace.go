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

	"github.com/maxlandon/wiregost/data_service/models"
)

const (
	// WorkspaceAPIPath is the API path to workspaces
	WorkspaceAPIPath = "/api/v1/workspaces/"
)

// WorkspaceHandler handles all HTTP requests concerning workspace management.
type WorkspaceHandler struct {
	// Env is needed to pass a DB connection pool
	*Env
}

// ServeHTTP dispatches and process workspace requests
func (wh *WorkspaceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	// List workspaces
	case r.Method == "GET":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		workspaces, err := wh.DB.Workspaces()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		json.NewEncoder(w).Encode(workspaces)

	// Add workspaces
	case r.Method == "POST":
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var names []string
		err = json.Unmarshal(b, &names)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = wh.DB.AddWorkspaces(names)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	// Delete workspace
	case r.Method == "DELETE":
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var ids []uint
		err = json.Unmarshal(b, &ids)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		deleted, err := wh.DB.DeleteWorkspaces(ids)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if deleted != int64(len(ids)) {
			http.Error(w, "Some ids are not valid", 500)
		}

	// Update workspace
	case r.Method == "PUT":
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var workspace models.Workspace
		err = json.Unmarshal(b, &workspace)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = wh.DB.UpdateWorkspace(workspace)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
}
