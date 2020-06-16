package rpc

import (
	"net"

	"github.com/maxlandon/wiregost/server/assets"
	"google.golang.org/grpc"
)

// StartClientListener - Listen for incoming console connections
func StartClientListener(host string, port int) (server *grpc.Server, ln net.Listener, err error) {

	// Get server config (certificates, etc)
	config := LoadUserServerTLSConfig(assets.ServerConfiguration.ServerHost)

	// Logging

	// Setup gRPC server
	server = SetupGRPC(config)

	// Register RPC Services
	RegisterServices(server)

	// Listen and serve gRPC
	ServeGRPC(server)

	return
}
