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

import (
	"github.com/maxlandon/wiregost/data_service/models"
)

const (
	WorkspaceApiPath = "/api/v1/workspaces/"
)

// Workspaces queries all workspaces to Data Service
func Workspaces() ([]models.Workspace, error) {
	client := newClient()
	req, err := client.newRequest("GET", WorkspaceApiPath, nil)
	if err != nil {
		return nil, err
	}

	var workspaces []models.Workspace
	err = client.do(req, &workspaces)

	return workspaces, err
}

func AddWorkspaces(names []string) error {
	client := newClient()
	req, err := client.newRequest("POST", WorkspaceApiPath, names)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}

func DeleteWorkspaces(ids []int) error {
	client := newClient()
	req, err := client.newRequest("DELETE", WorkspaceApiPath, ids)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}

func UpdateWorkspace(ws models.Workspace) error {
	client := newClient()
	req, err := client.newRequest("PUT", WorkspaceApiPath, ws)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}
