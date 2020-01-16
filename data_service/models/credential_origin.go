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

package models

type OriginCrackedPassword struct {
	ID        int `sql:"origin_id,notnull,on_delete:CASCADE"`
	Core      *CoreCredential
	FileName  string
	TaskID    int
	CreatedAt string
	UpdatedAt string
}

type OriginImport struct {
	ID        int `sql:"origin_id,notnull,on_delete:CASCADE"`
	Core      *CoreCredential
	TaskID    int
	CreatedAt string
	UpdatedAt string
}

type OriginManual struct {
	ID        int `sql:"origin_id,notnull,on_delete:CASCADE"`
	Core      *CoreCredential
	CreatedAt string
	UpdatedAt string
}

type OriginService struct {
	ID             int `sql:"origin_id,notnull,on_delete:CASCADE"`
	Core           *CoreCredential
	ServiceID      int `sql:"service_id,notnull,on_delete:CASCADE"`
	Service        *Service
	ModuleFullName string
	CreatedAt      string
	UpdatedAt      string
}

type OriginAgent struct {
	ID                int `sql:"origin_id,notnull,on_delete:CASCADE"`
	Core              *CoreCredential
	PostReferenceName string
	AgentID           string
	CreatedAt         string
	UpdatedAt         string
}
