package clients

import (
	"fmt"
	"net"
	"strconv"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/server/assets"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
)

// SetupGRPC - Setup gRPC security, and register all RPC services
func SetupGRPC() (server *grpc.Server) {

	// Get server config (certificates, etc)
	// config := LoadUserServerTLSConfig(assets.ServerConfiguration.ServerHost)

	// Get & set credentials
	// creds := credentials.NewTLS(config)

	// Options
	// opts := []grpc.ServerOption{
	//         grpc.Creds(creds),
	// }

	// Set authentication interceptors

	// Instantiate gRPC server
	server = grpc.NewServer()
	// server = grpc.NewServer(opts...)

	return
}

// Serve - Listen for incoming console connections
func Serve() (server *grpc.Server, ln net.Listener, err error) {

	// Logging

	// Checking for at least one user
	err = CreateDefaultUser()

	// Setup gRPC server
	server = SetupGRPC()

	// Register RPC Services
	RegisterServices(server)

	// Listen and serve gRPC
	ServeGRPC(server)

	return
}

// ServeGRPC - Start the Wiregost client gRPC server
func ServeGRPC(server *grpc.Server) {

	// Start listener
	ln, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", assets.ServerConfiguration.ServerHost, assets.ServerConfiguration.ServerPort))

	// Start server
	fmt.Println(tui.Green("gRPC:") + " Wiregost server running on " +
		assets.ServerConfiguration.ServerHost + ":" +
		strconv.Itoa(assets.ServerConfiguration.ServerPort))

	server.Serve(ln)
}
