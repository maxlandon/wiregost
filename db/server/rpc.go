package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	db "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	"github.com/maxlandon/wiregost/server/assets"
)

// StartRPCServices - Register all gRPC server components
func StartRPCServices() (err error) {

	// Setup & bind server connection
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", assets.ServerConfiguration.DatabaseRPCPort))

	// gRPC
	server := grpc.NewServer()

	// Users
	db.RegisterUserDBServer(server, &userServer{})

	// Certificates

	// Serve (blocking)
	server.Serve(lis)

	return
}
