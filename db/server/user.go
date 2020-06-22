package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/evilsocket/islazy/tui"
	db "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/certs"
	"github.com/maxlandon/wiregost/server/generate"
)

type userServer struct {
	*db.UnimplementedUserDBServer
}

func (*userServer) GetUsers(context.Context, *db.User) (out *db.Users, err error) {

	// Get users from db
	DB.Find(&out.Users).Where("name = ?", "wiregost").Where("password = ?", "wiregost")

	return
}

func (*userServer) AddUsers(ctx context.Context, in *db.AddUser) (out *db.Added, err error) {

	// Check if no user has the same name
	var users []*db.User
	errs := DB.Find(&users).Where("name = ?", in.User.Name).GetErrors()
	if len(errs) != 0 {
		fmt.Println(tui.Red("Error asking for users"))
	}
	if len(users) != 0 {
		return nil, errors.New("User ")
	}

	// Add user to DB
	DB.Create(in.User)

	// Here we need to generate all necessary certificates
	pub, priv, err := certs.UserClientGenerateCertificate(in.User.Name)
	if err != nil {
		return nil, err
	}

	// Save certificates keypair to DB
	cert := &serverpb.CertificateKeyPair{
		Hostname:    in.User.Name,
		CAType:      certs.UserCA,
		KeyType:     certs.ECCKey, // Users always have an ECC key anyway
		Certificate: pub,
		PrivateKey:  priv,
	}
	// Add it
	DB.Create(cert)

	// Find a way, when generating a certificate for the user, to use it to further
	// generate certificates that will be used for compiled consoles belonging to this user.

	// Save the certificates to a configuration file for use with consoles
	if in.WithConsoleFile {
		generate.UserConsoleConfig(in.User, pub, priv, in.WithServerDefault)

	}

	// Compile a console for this user with the above certificates, and output binary in user directory
	if in.WithConsole {
		generate.CompileObfuscatedConsole(in.User.Name, in.BinaryName, in.Send)
	}

	return &db.Added{User: in.User}, nil
}

func (*userServer) UpdateUsers(context.Context, *db.UpdateUser) (*db.Updated, error) {

	// We should not have to touch the certificates
	return nil, nil
}

func (*userServer) DeleteUsers(context.Context, *db.DeleteUser) (*db.Deleted, error) {

	// Maybe we will have to revoke the certificates
	return nil, nil
}
