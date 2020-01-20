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
	"context"

	"github.com/maxlandon/wiregost/data_service/models"
)

const (
	workspaceAPIPath = "/api/v1/workspaces/"
)

// Workspaces queries all workspaces to Data Service
func Workspaces(ctx context.Context) ([]models.Workspace, error) {
	client := newClient()
	req, err := client.newRequest(ctx, "GET", workspaceAPIPath, nil)
	if err != nil {
		return nil, err
	}

	var workspaces []models.Workspace
	err = client.do(req, &workspaces)

	return workspaces, err
}

// AddWorkspaces is used by client and server to add one or more workspaces
func AddWorkspaces(ctx context.Context, names []string) error {
	client := newClient()
	req, err := client.newRequest(ctx, "POST", workspaceAPIPath, names)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}

// DeleteWorkspaces is used by client and server to delete one or more workspaces
func DeleteWorkspaces(ctx context.Context, ids []uint) error {
	client := newClient()
	req, err := client.newRequest(ctx, "DELETE", workspaceAPIPath, ids)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}

// UpdateWorkspace is used by client and server to update a workspace
func UpdateWorkspace(ctx context.Context, ws models.Workspace) error {
	client := newClient()
	req, err := client.newRequest(ctx, "PUT", workspaceAPIPath, ws)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}
