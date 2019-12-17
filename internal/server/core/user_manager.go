package core

// This file contains the code used by the AuthenticationManager.
// Its role is to :
//		- check newly connected users,
//		- register them during first-time use,
//		- verify their admin rights,
//		- If admin, handle user creation and deletion

import (
	"fmt"
	"strings"

	// Here for testing
	"os/exec"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/server/db"
	"google.golang.org/grpc/metadata"

	"golang.org/x/net/context"
)

type UserManager struct {
	ConnectedUsers []db.User

	// DB Access
	database *DBManager
}

func NewUserManager() *UserManager {
	userManager := &UserManager{}

	return userManager
}

// -------------------------------------------------------------
// DATABASE FUNCTIONS

func (um *UserManager) GetUsers() ([]db.User, error) {
	var users []db.User
	err := um.database.DB.Model(&users).Select()
	return users, err
}

// func AddUser() {
// }

// -------------------------------------------------------------
// RPC FUNCTIONS

func (um *UserManager) ConnectUser(ctx context.Context, in *ConnectRequest) (*ConnectResponse, error) {
	// Get metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")

		// Respond
		if ctx.Value(clearanceCtx) == "clear" && ctx.Value(adminCtx) == true {
			// Add to ConnectedUsers
			var connected = false
			for _, u := range um.ConnectedUsers {
				if u.Name == clientLogin && u.PasswordHashString == clientPassword {
					connected = true
				}
			}
			if connected == false {
				var user db.User
				err := um.database.DB.Model(&user).Where("name = ?", clientLogin).Where("password_hash_string = ?", clientPassword).Select()
				um.ConnectedUsers = append(um.ConnectedUsers, user)
				if err != nil {
					fmt.Println(err)
				}
				return &ConnectResponse{Clearance: "clear", Admin: true}, nil
			}
			return &ConnectResponse{Clearance: "clear", Admin: true}, nil
		}
		if ctx.Value(clearanceCtx) == "clear" && ctx.Value(adminCtx) == false {
			return &ConnectResponse{Clearance: "clear", Admin: false}, nil
		}
		if ctx.Value(clearanceCtx) == "reg" {
			users, _ := um.GetUsers()
			for _, v := range users {
				// If user is in db
				if v.Name == clientLogin {
					v.PasswordHashString = clientPassword
					_, err := um.database.DB.Model(&v).Set(
						"password_hash_string = ?password_hash_string").Where("name = ?name").Update()
					if err != nil {
						fmt.Println(tui.Red("Error: Failed to save PasswordHash to Database."))
						fmt.Println(err)
					}
					return &ConnectResponse{Clearance: "reg"}, nil
				}
			}
		}
	}
	if ctx.Value(clearanceCtx) == "none" {
		return &ConnectResponse{Clearance: "none"}, nil
	}

	return &ConnectResponse{}, nil
}

// Disconnecting a user
func (um *UserManager) DisconnectUser(ctx context.Context, in *DisconnectRequest) (*DisconnectResponse, error) {
	// Get metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")

		// Respond
		if ctx.Value(clearanceCtx) == "clear" {
			// Remove from connected users
			var newConnected []db.User
			for _, u := range um.ConnectedUsers {
				if clientLogin != u.Name {
					newConnected = append(newConnected, u)
				}
			}
			um.ConnectedUsers = newConnected
			return &DisconnectResponse{Disconnected: true}, nil
		}
	}

	return &DisconnectResponse{Disconnected: false}, nil
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
