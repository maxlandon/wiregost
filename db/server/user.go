package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	db "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/certs"
	"github.com/maxlandon/wiregost/server/generate"
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
func (*userServer) AddUsers(ctx context.Context, in *db.AddUser) (out *db.Added, err error) {

	// Check if no user has the same name
	var users []*db.User
	DB.Find(users).Where("name = ?", in.User.Name)
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
		GenerateUserConsoleConfig(in.User, pub, priv, in.WithServerDefault)

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

// ClientConfig - Struct containing user information and certificate
type ClientConfig struct {
	User          string `json:"user"`
	LHost         string `json:"lhost"`
	LPort         int    `json:"lport"`
	CACertificate string `json:"ca_certificate"`
	PrivateKey    string `json:"private_key"`
	Certificate   string `json:"certificate"`
	IsDefault     bool   `json:"is_default"`
}

func GenerateUserConsoleConfig(user *db.User, pub []byte, priv []byte, isDefault bool) (err error) {

	// Make config
	caCertPEM, _, _ := certs.GetCertificateAuthorityPEM(certs.UserCA)
	config := ClientConfig{
		User:          user.Name,
		LHost:         assets.ServerConfiguration.ServerHost,
		LPort:         int(assets.ServerConfiguration.ServerPort),
		CACertificate: string(caCertPEM),
		PrivateKey:    string(priv),
		Certificate:   string(pub),
		IsDefault:     isDefault,
	}

	// Save to file
	configJSON, _ := json.Marshal(config)

	saveTo, _ := filepath.Abs(path.Join(assets.GetRootAppDir(), "users"))

	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		err = os.MkdirAll(saveTo, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Cannot write to wiregost root directory %s", err)
		}
	}

	fi, err := os.Stat(saveTo)
	if fi.IsDir() {
		def := ""
		if isDefault {
			def = "default"
		} else {
			def = "normal"
		}
		filename := fmt.Sprintf("%s_%s_%s.cfg", filepath.Base(user.Name), filepath.Base(assets.ServerConfiguration.ServerHost), def)
		saveTo = filepath.Join(saveTo, filename)
	}
	err = ioutil.WriteFile(saveTo, configJSON, 0600)
	if err != nil {
		return fmt.Errorf("Failed to write config to: %s (%v) \n", saveTo, err)
	}

	return
}
