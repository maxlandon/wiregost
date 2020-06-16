package rpc

import (
	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	"google.golang.org/grpc"
)

// RegisterServices - Register all gRPC services available to console users
func RegisterServices(server *grpc.Server) {

	// Connection (authentication & information)
	clientpb.RegisterConnectionRPCServer(server, &connectionServer{})
}
