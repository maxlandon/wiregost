package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	client "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	"github.com/maxlandon/wiregost/server/assets"
)

// RegisterRPCServices - Register all gRPC server components
func RegisterRPCServices() (err error) {

	// Setup & bind server connection
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", assets.ServerConfiguration.DatabaseRPCPort))

	// gRPC
	server := grpc.NewServer()

	// Users
	client.RegisterConnectionRPCServer(server, &userServer{})

	// Certificates

	// Serve (blocking)
	server.Serve(lis)

	return
}
