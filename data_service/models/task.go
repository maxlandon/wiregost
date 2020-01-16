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

type Task struct {
	WorkspaceID int `sql:"workspace_id,notnull,on_delete:CASCADE"`
	Workspace   *Workspace
	CreatedBy   string
	Module      string
	ModuleUUID  string
	CompletedAt string
	Path        string
	Info        string
	Description string
	Progress    int
	Options     string
	Error       string
	Result      string
	CreatedAt   string
	UpdatedAt   string
}

type TaskCredential struct {
	*Task
	ID           int `sql:"task_id,pk"`
	CredentialID int `sql:"credential_id,notnull,on_delete:CASCADE"`
	Credential   *CoreCredential
	CreatedAt    string
	UpdatedAt    string
}

type TaskHost struct {
	*Task
	ID        int `sql:"task_id,pk"`
	HostID    int `sql:"host_id,notnull,on_delete:CASCADE"`
	Host      *Host
	CreatedAt string
	UpdatedAt string
}

type TaskService struct {
	*Task
	ID        int `sql:"task_id,pk"`
	ServiceID int `sql:"service_id,notnull,on_delete:CASCADE"`
	Service   *Service
	CreatedAt string
	UpdatedAt string
}

type TaskAgent struct {
	ID        int `sql:"task_id,pk"`
	AgentID   string
	CreatedAt string
	UpdatedAt string
}
