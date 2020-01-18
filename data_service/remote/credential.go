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

package remote

import "github.com/maxlandon/wiregost/data_service/models"

const (
	credentialAPIPath = "/api/v1/credentials/"
)

// Creds queries all credentials to Data Service
func Creds() ([]models.Credential, error) {
	client := newClient()
	req, err := client.newRequest("GET", credentialAPIPath, nil)
	if err != nil {
		return nil, err
	}

	var workspaces []models.Credential
	err = client.do(req, &workspaces)

	return workspaces, err
}

// CreateCredential is used by client and server to add a Credential set.
func CreateCredential(c *models.Credential) error {
	client := newClient()
	req, err := client.newRequest("POST", credentialAPIPath, c)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}

// DeleteCredentials is used by client and server to delete one or more credential sets
func DeleteCredentials(ids []int) error {
	client := newClient()
	req, err := client.newRequest("DELETE", credentialAPIPath, ids)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}

// UpdateCredential is used by client and server to update a credential set
func UpdateCredential(ws models.Credential) error {
	client := newClient()
	req, err := client.newRequest("PUT", credentialAPIPath, ws)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}
