package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/evilsocket/islazy/tui"
	db "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/assets"
)

// StartRPCServices - Register all gRPC server components
func StartRPCServices() (err error) {

	// Setup & bind server connection
	fmt.Println(assets.ServerConfiguration.DatabaseRPCHost)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", assets.ServerConfiguration.DatabaseRPCPort))

	// gRPC
	server := grpc.NewServer()
	RegisterServices(server)

	// Notify correct start
	fmt.Println(tui.Green("DB:") + " Wiregost Database running")

	// Serve (blocking)
	server.Serve(lis)

	return
}

// RegisterServices - All RPC services bound to gRPC server
func RegisterServices(server *grpc.Server) {

	db.RegisterUserDBServer(server, &userServer{})               // Users
	serverpb.RegisterCertificateRPCServer(server, &certServer{}) // Certificates
}
