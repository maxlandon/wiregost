package db

import (
	"errors"
)

// User type stores credential information for a Wiregost user.
type User struct {
	ID                 int
	Name               string
	PasswordHashString string
	PasswordHash       [32]byte
	AccessToken        string
	Administrator      bool
}

// AddUser is a function only available to admin users.
func (user *User) AddUser(name string) (err error) {
	if !user.Administrator {
		return errors.New("Cannot add user: You do not have admin rights")
	}

	return nil
}
