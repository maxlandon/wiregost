package connection

import (
	"google.golang.org/grpc"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	modpb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

var (
	// ConnectionRPC - All authentication and console utils
	ConnectionRPC clientpb.ConnectionRPCClient
	// UserRPC - User commands
	UserRPC dbpb.UserDBClient
	// ConfigRPC - Config commands
	ConfigRPC clientpb.ConfigRPCClient

	// CompilerRPC - Compiler commands
	CompilerRPC serverpb.CompilerClient

	// EventsRPC - Processes incoming events (modules, sessions, etc.)
	EventsRPC serverpb.EventsClient

	ModuleRPC modpb.StackClient
)

// RegisterRPCClients - Binds all gRPC clients to the newly established & authenticated connection.
func RegisterRPCClients(conn *grpc.ClientConn) (err error) {

	// User
	UserRPC = dbpb.NewUserDBClient(conn)

	// Config
	ConfigRPC = clientpb.NewConfigRPCClient(conn)

	// Compiler
	CompilerRPC = serverpb.NewCompilerClient(conn)

	// Events
	EventsRPC = serverpb.NewEventsClient(conn)

	// Module & Stack
	ModuleRPC = modpb.NewStackClient(conn)

	return
}
