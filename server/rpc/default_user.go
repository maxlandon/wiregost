package rpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	db "github.com/maxlandon/wiregost/db/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	"github.com/maxlandon/wiregost/server/certs"
)

// CreateDefaultUser - When starting, the server automatically checks if at least one user exists. Creates one if not.
func CreateDefaultUser() (err error) {

	// Check if, instead of users in DB, we already have user certificates: if yes, it is not normal, we should have both
	userCerts := certs.UserClientListCertificates()

	// GetUsers from db
	res, err := db.Users.GetUsers(context.Background(), &dbpb.User{}, grpc.EmptyCallOption{})

	// If error, we might have to check further instead of directly creating a new default user
	if err != nil {
		fmt.Println(err.Error())
	}

	// If nil, create user wiregost
	if res.Users == nil && len(userCerts) == 0 {
		// Add new user to DB
		user := &dbpb.User{
			Name:     "wiregost",
			Password: []byte("wiregost"),
			Admin:    true,
		}

		add := &dbpb.AddUser{
			User:              user,
			WithConsole:       true,
			WithConsoleFile:   true,
			WithServerDefault: true,
			BinaryName:        "wiregost_default_console",
		}
		created, err := db.Users.AddUsers(context.Background(), add, grpc.EmptyCallOption{})
		if err != nil || created.User == nil {
			return err
		}

		// Log creation of a new default user (will be logged twice because DB does it also)

		return nil
	}

	// If not nil, return
	return
}
