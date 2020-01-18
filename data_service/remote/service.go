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
	"errors"
	"net/http"
	"strconv"

	"github.com/maxlandon/wiregost/data_service/models"
)

const (
	serviceAPIPath = "/api/v1/services/"
)

// Services queries all Services to Data Service
func Services() ([]models.Service, error) {
	client := newClient()
	req, err := client.newRequest("GET", serviceAPIPath, nil)
	if err != nil {
		return nil, err
	}

	var services []models.Service
	err = client.do(req, &services)

	return services, err
}

// GetService returns a single Service, based on various options passed as search filters.
func GetService(opts map[string]string) (*models.Service, error) {
	client := newClient()
	var req *http.Request
	var err error

	// Check for ID (Currently only way to get a single host. No search
	// based on other options is possible here, because of how the data_service
	// dispatches requests)
	id, found := opts["service_id"]
	if found {
		req, err = client.newRequest("GET", serviceAPIPath+id, opts)
		if err != nil {
			return nil, err
		}
	} else {
		err := errors.New("No ServiceID is specified")
		return nil, err
	}

	var service *models.Service
	err = client.do(req, &service)

	return service, err
}

// ReportService adds a Service to the database
func ReportService(s *models.Service) (*models.Service, error) {
	client := newClient()
	req, err := client.newRequest("POST", serviceAPIPath, s)
	if err != nil {
		return nil, err
	}

	var service *models.Service
	err = client.do(req, &service)

	return service, err
}

// UpdateService updates a Service properties
func UpdateService(s *models.Service) (*models.Service, error) {
	client := newClient()
	serviceID := strconv.Itoa(s.ID)
	req, err := client.newRequest("PUT", serviceAPIPath+serviceID, s)
	if err != nil {
		return nil, err
	}

	var service *models.Service
	err = client.do(req, &service)

	return service, err
}

// DeleteServices deletes one or more Services from the database
func DeleteServices(ids []int) error {
	client := newClient()
	req, err := client.newRequest("DELETE", serviceAPIPath, ids)
	if err != nil {
		return err
	}
	err = client.do(req, nil)

	return err
}
