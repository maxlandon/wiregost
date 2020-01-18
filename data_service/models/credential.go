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

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

type Credential struct {
	ID          int `sql:"credential_id,pk"`
	WorkspaceID int `sql:"workspace_id,notnull,on_delete:CASCADE"`
	Workspace   *Workspace
	ServiceID   int `sql:"service_id,on_delete:CASCADE"`
	Service     *Service
	RealmID     int `sql:"realm_id"`
	Realm       *Realm
	PrivateID   int `sql:"private_credential_id"`
	Private     *PrivateCredential
	PublicID    int `sql:"public_credential_id"`
	Public      *PublicCredential
	Active      bool
	Proof       string
	SourceID    int
	SourceType  string
	LoginsCount int
	Logins      []Login
	CreatedAt   string
	UpdatedAt   string
}

// NewCredential instantiates a new credential set, whether it is made
// of a public/private pair, one of both, or none.
func NewCredential() *Credential {
	// Get good random id
	rand.Seed(time.Now().Unix())
	id := rand.Int()

	cred := &Credential{
		ID:        id,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05"),
	}

	return cred
}

// PublicCredential is the public part of a Credential set
type PublicCredential struct {
	ID        int    `sql:"public_credential_id,pk,notnull"`
	UserName  string `sql:"credential_user,notnull"`
	Type      string `sql:"public_type,notnull"`
	CreatedAt string
	UpdatedAt string
}

// NewPublicCredential instantiates a new Public Credential
func NewPublicCredential() *PublicCredential {
	// Get good random id
	rand.Seed(time.Now().Unix())
	id := rand.Int()

	pub := &PublicCredential{
		ID:        id,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05"),
	}

	return pub
}

// PrivateCredential is the private part of a Credential set
type PrivateCredential struct {
	ID        int    `sql:"private_credential_id,pk,notnull"`
	Data      string `sql:"credential_password,notnull"`
	Type      string `sql:"private_type,notnull"`
	CreatedAt string
	UpdatedAt string
	JTRFormat string
}

// NewPrivateCredential instantiates a new Private Credential
func NewPrivateCredential() *PrivateCredential {
	// Get good random id
	rand.Seed(time.Now().Unix())
	id := rand.Int()

	priv := &PrivateCredential{
		ID:        id,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05"),
	}

	return priv
}

// Creds returns all credentials in the Database
func (db *DB) Creds() ([]*Credential, error) {
	var creds []*Credential
	err := db.Model(&creds).
		Relation("PrivateCredential").
		Relation("PublicCredential").
		Select()
	if err != nil {
		return nil, err
	}
	return creds, err
}

// EachCred returns all credentials from the worspace given as argument.
func (db *DB) EachCred(workspaceID int) ([]*Credential, error) {
	var creds []*Credential
	err := db.Model(&creds).
		Where("workspace_id = ?", workspaceID).
		Relation("PrivateCredential").
		Relation("PublicCredential").
		Select()
	if err != nil {
		return nil, err
	}

	return creds, err
}

// GetCred returns a credential based on options passed as argument
func (db *DB) GetCred(opts map[string]string) (*Credential, error) {
	c := new(Credential)

	// Find credential by ID, and return it if found, return error otherwise.
	id, found := opts["credential_id"]
	if found {
		id, _ := strconv.Atoi(id)
		err := db.Model(c).Where("credential_id= ?", id).
			Relation("PrivateCredential").
			Relation("PublicCredential").
			Select()
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	// Workspace ID is required if no CredentialID is given, and needs to be cast
	ws, found := opts["workspace_id"]
	if !found {
		return nil, errors.New("Workspace ID is required")
	}
	wsID, _ := strconv.Atoi(ws)

	// Find credential by username
	user, found := opts["username"]
	if found {
		err := db.Model(c).Where("workspace_id = ?", wsID).
			Where("address = ?", user).Select()
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	return nil, nil
}

// Find or create a credential matching this type/data
func (db *DB) FindOrCreateCred(opts map[string]string) (*Credential, error) {
	cred, err := db.GetCred(opts)
	// If not cred is found, create one and fill values given
	if cred == nil {
		c := NewCredential()
		ws, found := opts["workspace_id"]
		if found {
			c.WorkspaceID, _ = strconv.Atoi(ws)
		}

		return db.ReportAuthInfo(*c)
	}

	return cred, err

}

// Creates a set of credentials in the database
func (db *DB) ReportAuthInfo(c Credential) (*Credential, error) {
	// Add Credential (no need to set ID, already exists with NewHost())
	err := db.Insert(&c)
	if err != nil {
		return nil, err
	}

	return &c, err

}

// DeleteCreds deletes Credential entries based on the IDs passed as arguments
func (db *DB) DeleteCreds(ids []int) (rows int, err error) {
	c := new(Credential)
	var deleted int
	for _, id := range ids {
		res, err := db.Model(c).Where("credential_id = ?", id).Delete()
		deleted += res.RowsAffected()
		if err != nil {
			return deleted, err
		}
	}

	return deleted, nil
}

// UpdateCredential updates a Credential entry with the Host struct passed as argument.
func (db *DB) UpdateCredential(c Credential) (*Credential, error) {
	c.UpdatedAt = time.Now().Format("2006-01-02T15:04:05")
	err := db.Update(&c)
	if err != nil {
		return nil, err
	}

	return &c, err
}
