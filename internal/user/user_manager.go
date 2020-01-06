package user

import (
	"github.com/maxlandon/wiregost/internal/db"
	"github.com/maxlandon/wiregost/internal/messages"
)

var (
	// AuthReqs is used to send authentication requests
	AuthReqs = make(chan messages.ClientRequest)
	// AuthResp is used to send authentication responses
	AuthResp = make(chan messages.ClientRequest)
)

// Manager has access to the list of saved users and keeps track of connected ones.
type Manager struct {
	ConnectedUsers []db.User
	// DB Access
	database *db.Manager
}

// NewManager instantiates a new User Manager, and handles authentication requests.
func NewManager() *Manager {
	userManager := &Manager{
		database: db.NewDBManager(),
	}

	go userManager.authenticate()

	return userManager
}

func (um *Manager) getUsers() ([]db.User, error) {
	var users []db.User
	err := um.database.DB.Model(&users).Select()
	return users, err
}

func (um *Manager) authenticate() {
	for {
		msg := <-AuthReqs
		users, _ := um.getUsers()
		registered := false
		connected := false
		var user db.User
		for _, u := range users {
			if u.Name == msg.UserName && u.PasswordHashString == msg.UserPassword {
				registered = true
				user = u
				msg.UserID = u.ID
			}
		}
		for _, u := range um.ConnectedUsers {
			if u.Name == msg.UserName && u.PasswordHashString == msg.UserPassword {
				connected = true
				user = u
				msg.UserID = u.ID
			}
		}
		if registered == true && connected == true {
			AuthResp <- msg
		}
		if registered == true && connected == false {
			um.ConnectedUsers = append(um.ConnectedUsers, user)
			AuthResp <- msg

		}
		if registered == false {
			msg.UserID = 0
			AuthResp <- msg
		}
	}
}
