package rpc

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/server/assets"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// SetupGRPC - Setup gRPC security, and register all RPC services
func SetupGRPC(config *tls.Config) (server *grpc.Server) {

	// Get & set credentials
	creds := credentials.NewTLS(config)

	// Options
	opts := []grpc.ServerOption{
		grpc.Creds(creds),
	}

	// Set authentication interceptors

	// Instantiate gRPC server
	server = grpc.NewServer(opts...)

	return
}

// ServeGRPC - Start the Wiregost client gRPC server
func ServeGRPC(server *grpc.Server) {

	// Start listener
	ln, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", assets.ServerConfiguration.ServerHost, assets.ServerConfiguration.ServerPort))

	// Start server
	fmt.Println(tui.Green("gRPC:") + " Wiregost server running")
	server.Serve(ln)
}
