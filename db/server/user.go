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
	return res, nil
	// return nil, nil
}
func (*userServer) AddUsers(context.Context, *db.AddUser) (*db.Added, error) {
	return nil, nil
}
func (*userServer) UpdateUsers(context.Context, *db.UpdateUser) (*db.Updated, error) {
	return nil, nil
}
func (*userServer) DeleteUsers(context.Context, *db.DeleteUser) (*db.Deleted, error) {
	return nil, nil
}
