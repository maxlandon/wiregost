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
	"strconv"

	"github.com/maxlandon/wiregost/data_service/models"
)

const (
	hostAPIPath = "/api/v1/hosts/"
)

// Hosts queries all Hosts to Data Service, with optional search filters passed in a map
func Hosts(ctx context.Context, opts map[string]interface{}) (hosts []models.Host, err error) {
	client := newClient()
	req, err := client.newRequest(ctx, "GET", hostAPIPath, opts)
	if err != nil {
		return nil, err
	}

	err = client.do(req, &hosts)

	return hosts, err
}

// GetHost returns a single host, based on various options passed as search filters.
func GetHost(ctx context.Context, opts map[string]interface{}) (host *models.Host, err error) {
	client := newClient()

	req, err := client.newRequest(ctx, "POST", hostAPIPath+"search", opts)
	if err != nil {
		return nil, err
	}

	err = client.do(req, &host)

	return host, err
}

// ReportHost adds a Host to the database
func ReportHost(ctx context.Context, opts map[string]interface{}) (host *models.Host, err error) {
	client := newClient()
	req, err := client.newRequest(ctx, "POST", hostAPIPath, opts)
	if err != nil {
		return nil, err
	}

	err = client.do(req, &host)

	return host, err
}

// UpdateHost updates a Host properties
func UpdateHost(h *models.Host) (host *models.Host, err error) {
	client := newClient()
	hostID := string(h.ID)
	req, err := client.newRequest(nil, "PUT", hostAPIPath+hostID, h)
	if err != nil {
		return nil, err
	}

	err = client.do(req, &host)

	return host, err
}

// DeleteHost deletes a Host from the database
func DeleteHost(ctx context.Context, id int) error {
	client := newClient()
	hostID := strconv.Itoa(id)
	req, err := client.newRequest(ctx, "DELETE", hostAPIPath+hostID, id)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}
