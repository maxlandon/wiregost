package server

import (
	"context"

	db "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
)

type userServer struct {
	*db.UnimplementedUserDBServer
}

func (*userServer) GetUsers(context.Context, *db.User) (*db.Users, error) {

	res := &db.Users{}
	res.Users = append(res.Users, &db.User{Name: "wiregostlong", Password: []byte("wiregost")})

	// Get users from db
	// DB.Find(res.Users)

	return res, nil
}
func (*userServer) AddUsers(context.Context, *db.AddUser) (*db.Added, error) {

	// Check if no user has the same name

	// Add user to DB

	// Here we need to generate all necessary certificates

	// Save certificates keypair to DB

	// Find a way, when generating a certificate for the user, to use it to further
	// generate certificates that will be used for compiled consoles belonging to this user.

	// Also save the certificates to a configuration file for use with consoles

	// Compile a console for this user with the above certificates, and output binary in user directory

	return nil, nil
}
func (*userServer) UpdateUsers(context.Context, *db.UpdateUser) (*db.Updated, error) {

	// We should not have to touch the certificates
	return nil, nil
}
func (*userServer) DeleteUsers(context.Context, *db.DeleteUser) (*db.Deleted, error) {

	// Maybe we will have to revoke the certificates
	return nil, nil
}
