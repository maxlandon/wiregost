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
	"strings"

	"github.com/maxlandon/wiregost/data_service/models"
)

const (
	// CredentialAPIPath is the API path to credentials
	CredentialAPIPath = "/api/v1/credentials/"
)

type CredentialHandler struct {
	// Env is needed to pass a DB connection pool
	*Env
}

// ServeHTTP dispatches and process credentials requests
func (ch *CredentialHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Check if id is provided in URL, will influence dispatch
	id := strings.TrimPrefix(r.URL.Path, CredentialAPIPath)

	switch {
	case id == "":
		switch {
		case r.Method == "GET":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			creds, err := ch.DB.Creds()
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}

			json.NewEncoder(w).Encode(creds)

		// Add a Host
		case r.Method == "POST":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var cred *models.Credential
			err = json.Unmarshal(b, &cred)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			cred, err = ch.DB.ReportAuthInfo(*cred)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			err = json.NewEncoder(w).Encode(cred)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

		// Delete several Hosts
		case r.Method == "DELETE":
			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var ids []int
			err = json.Unmarshal(b, &ids)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			deleted, err := ch.DB.DeleteCreds(ids)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			if deleted != len(ids) {
				http.Error(w, "Some ids are not valid", 500)
			}

		case r.Method == "PUT":
		}

	// ID is there, a credential is specified
	case id != "":
		switch {
		// Return a single credential, based on ID
		case r.Method == "GET":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var opts map[string]string
			err = json.Unmarshal(b, &opts)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			cred, err := ch.DB.GetCred(opts)
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}

			json.NewEncoder(w).Encode(cred)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

		case r.Method == "POST":
		case r.Method == "DELETE":

		// Update a Credential
		case r.Method == "PUT":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var cred *models.Credential
			err = json.Unmarshal(b, &cred)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			cred, err = ch.DB.UpdateCredential(*cred)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			err = json.NewEncoder(w).Encode(cred)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	}
}
