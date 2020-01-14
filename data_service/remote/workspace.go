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
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/maxlandon/wiregost/data_service/models"
)

const (
	// Temporary protocol string, should be pulled out of config
	Protocol         = "http://localhost:8000"
	WorkspaceApiPath = "/api/v1/workspaces/"
)

// Workspaces queries all workspaces to Data Service
func Workspaces() ([]models.Workspace, error) {
	resp, err := http.Get(Protocol + WorkspaceApiPath)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var workspaces []models.Workspace
	err = json.Unmarshal(body, &workspaces)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}

func FindWorkspace(name string) (models.Workspace, error) {

}

// func AddWorkspaces(names []string) error {
//
// }
//
// func DeleteWorkspaces(ids []int) error {
//
// }
//
// func UpdateWorkspace(ws models.Workspace) error {
//
// }
