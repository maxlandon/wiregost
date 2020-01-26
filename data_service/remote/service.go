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
	serviceAPIPath = "/api/v1/services/"
)

// Services queries all Services to Data Service, with optional search filters passed in a map
func Services(ctx context.Context, opts map[string]interface{}) (services []models.Port, err error) {
	client := newClient()
	req, err := client.newRequest(ctx, "GET", serviceAPIPath, opts)
	if err != nil {
		return nil, err
	}

	err = client.do(req, &services)

	return services, err
}

// ReportService adds a Service to the database
func ReportService(ctx context.Context, opts map[string]interface{}) (service *models.Port, err error) {
	client := newClient()
	req, err := client.newRequest(ctx, "POST", serviceAPIPath, opts)
	if err != nil {
		return nil, err
	}

	err = client.do(req, &service)

	return service, err
}

// UpdateService updates a Service properties
func UpdateService(s *models.Port) (service *models.Port, err error) {
	client := newClient()
	portID := string(s.ID)
	req, err := client.newRequest(nil, "PUT", serviceAPIPath+portID, s)
	if err != nil {
		return nil, err
	}

	err = client.do(req, &service)

	return service, err
}

// DeleteServices deletes one or more Services from the database
func DeleteServices(ctx context.Context, opts map[string]interface{}) error {
	client := newClient()
	req, err := client.newRequest(ctx, "DELETE", serviceAPIPath, opts)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}
