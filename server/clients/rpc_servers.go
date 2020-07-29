package clients

import (
	"google.golang.org/grpc"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/generate"
)

// RegisterServices - Register all gRPC services available to console users
func RegisterServices(server *grpc.Server) {

	// Connection (authentication & information)
	clientpb.RegisterConnectionRPCServer(server, &connectionServer{})

	// Implant & Console Compilation
	serverpb.RegisterCompilerServer(server, &generate.Compiler{})
}
