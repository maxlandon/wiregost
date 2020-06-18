package client

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/evilsocket/islazy/tui"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/assets"
	// "github.com/maxlandon/wiregost/server/assets"
)

var (
	// Users - User commands
	Users dbpb.UserDBClient
	// Certs - Certificate management
	Certs serverpb.CertificateRPCClient
)

// RegisterDBClients - Binds all Database gRPC clients to another dedicated connection
func RegisterDBClients(conn *grpc.ClientConn) (err error) {

	// User
	Users = dbpb.NewUserDBClient(conn)

	return
}

// ConnectToDatabase - Client method used by consoles and server to query DB remotely
func ConnectToDatabase(host string, port int, pub string, priv string) (err error) {

	// cert := tls.LoadX509KeyPair(pub, priv)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())

	// Register all DB clients
	RegisterDBClients(conn)

	return
}

// ConnectServerToDB - Client method used by server to query DB remotely
func ConnectServerToDB() (err error) {

	// Certificates from server conf
	conf := assets.ServerConfiguration
	// cert := tls.LoadX509KeyPair(conf.DBCertificate, conf.DBPrivateKey)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", conf.DatabaseRPCHost, conf.DatabaseRPCPort), grpc.WithInsecure())

	// Register all DB clients
	RegisterDBClients(conn)

	// Notify connection
	fmt.Println(tui.Green("DB:") + " Wiregost server connected to DB")

	return
}
