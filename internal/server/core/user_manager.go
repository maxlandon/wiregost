package core

// This file contains the code used by the AuthenticationManager.
// Its role is to :
//		- check newly connected users,
//		- register them during first-time use,
//		- verify their admin rights,
//		- If admin, handle user creation and deletion

import (
	"log"
	// Here for testing
	"os/exec"

	"github.com/maxlandon/wiregost/internal/server/db"

	"golang.org/x/net/context"
)

type UserManager struct {
	ConnectedUsers []db.User
}

func NewUserManager() *UserManager {
	userManager := &UserManager{}

	return userManager
}

// -------------------------------------------------------------
// DATABASE FUNCTIONS

func GetUsers() ([]db.User, error) {
	database := db.NewDBManager()
	var users []db.User
	err := database.Database.Model(&users).Select()
	return users, err
}

// func AddUser() {
// }

// -------------------------------------------------------------
// RPC FUNCTIONS

// Registering a user, during first connection
func (um *UserManager) RegisterUser(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {

	log.Printf("Received registration request from user %s", in.Name)

	users, _ := GetUsers()
	for _, v := range users {
		if v.Name == in.Name {
			return &RegisterResponse{Registered: true}, nil
		}
	}
	return &RegisterResponse{Registered: false}, nil
}

func (um *UserManager) ConnectUser(ctx context.Context, in *ConnectRequest) (*ConnectResponse, error) {

	return &ConnectResponse{}, nil
}

// Disconnecting a user
func (um *UserManager) DisconnectUser(ctx context.Context, in *DisconnectRequest) (*DisconnectResponse, error) {

	return &DisconnectResponse{}, nil
}

// Listing all users
func (um *UserManager) ListUsers(*ListUsersRequest, UserManager_ListUsersServer) error {

	// Here for shutting up the compiler
	cmd := exec.Command("sh")
	err := cmd.Run()
	return err
}

// Creating a user (admin)
func (um *UserManager) CreateUser(ctx context.Context, in *CreateUserRequest) (*CreateUserResponse, error) {

	return &CreateUserResponse{}, nil
}

// Deleting a user (admin)
func (um *UserManager) DeleteUser(ctx context.Context, in *DeleteUserRequest) (*DeleteUserResponse, error) {

	return &DeleteUserResponse{}, nil
}

// Giving admin rights to an already existing user
func (um *UserManager) GiveAdminRights(ctx context.Context, in *AdminRightsRequest) (*AdminRightsResponse, error) {

	return &AdminRightsResponse{}, nil
}
