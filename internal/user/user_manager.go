package user

import (
	"github.com/maxlandon/wiregost/internal/db"
	"github.com/maxlandon/wiregost/internal/messages"
)

var (
	AuthReqs = make(chan messages.ClientRequest)
	AuthResp = make(chan messages.ClientRequest)
)

type UserManager struct {
	ConnectedUsers []db.User
	// DB Access
	database *db.DBManager
}

func NewUserManager() *UserManager {
	userManager := &UserManager{
		database: db.NewDBManager(),
	}

	go userManager.Authenticate()

	return userManager
}

func (um *UserManager) GetUsers() ([]db.User, error) {
	var users []db.User
	err := um.database.DB.Model(&users).Select()
	return users, err
}

func (um *UserManager) Authenticate() {
	for {
		msg := <-AuthReqs
		users, _ := um.GetUsers()
		registered := false
		connected := false
		var user db.User
		for _, u := range users {
			if u.Name == msg.UserName && u.PasswordHashString == msg.UserPassword {
				registered = true
				user = u
				msg.UserId = u.Id
			}
		}
		for _, u := range um.ConnectedUsers {
			if u.Name == msg.UserName && u.PasswordHashString == msg.UserPassword {
				connected = true
				user = u
				msg.UserId = u.Id
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
			msg.UserId = 0
			AuthResp <- msg
		}
	}
}
