package db

import (
	"errors"
)

type User struct {
	Id                 int
	Name               string
	PasswordHashString string
	PasswordHash       [32]byte
	AccessToken        string
	Administrator      bool

	// CurrentWorkspace *Workspace
}

// Function only available to admin users.
func (user *User) AddUser(name string) (err error) {
	if !user.Administrator {
		return errors.New("Cannot add user: You do not have admin rights")
	} else {

	}

	return nil
}
