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

type Credential struct {
	ServiceID  int `sql:"service_id,notnull,on_delete:CASCADE"`
	Service    *Service
	User       string `sql:"credential_user"`
	Password   string `sql:"credential_password"`
	Active     bool
	Proof      string
	PType      string `sql:"credential_private_type"`
	SourceID   int
	SourceType string
	CreatedAt  string
	UpdatedAt  string
}

type CoreCredential struct {
	ID          int `sql:"credential_id,pk"`
	OriginID    int
	OriginType  string
	PrivateID   int `sql:"private_credential_id"`
	PublicID    int `sql:"public_credential_id"`
	RealmID     int `sql:"realm_id,notnull,on_delete:CASCADE"`
	Realm       *Realm
	WorkspaceID int `sql:"workspace_id,notnull,on_delete:CASCADE"`
	Workspace   *Workspace
	LoginsCount int
	Logins      []Login
	CreatedAt   string
	UpdatedAt   string
}

type PublicCredential struct {
	ID        int `sql:"public_credential_id,pk,notnull,on_delete:CASCADE"`
	Core      *CoreCredential
	UserName  string `sql:"credential_user,notnull"`
	Type      string
	CreatedAt string
	UpdatedAt string
}

type PrivateCredential struct {
	ID        int `sql:"private_credential_id,pk,notnull,on_delete:CASCADE"`
	Core      *CoreCredential
	Data      string `sql:"credential_password,notnull"`
	Type      string `sql:"credential_private_type,notnull"`
	CreatedAt string
	UpdatedAt string
	JTRFormat string
}
