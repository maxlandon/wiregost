package clients

import (
	"google.golang.org/grpc"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/generate"
	"github.com/maxlandon/wiregost/server/module/stack"
)

// RegisterServices - Register all gRPC services available to console users
func RegisterServices(server *grpc.Server) {

	// Connection (authentication & information)
	clientpb.RegisterConnectionRPCServer(server, &connectionServer{})

	// Implant & Console Compilation
	serverpb.RegisterCompilerServer(server, &generate.Compiler{})

	// Stack & Modules
	modulepb.RegisterStackServer(server, stack.Stacks)
}
