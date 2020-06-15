package connection

import (
	"google.golang.org/grpc"

	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
)

var (
	// UserRPC - User commands
	UserRPC dbpb.UserDBClient
)

// RegisterRPCClients - Binds all gRPC clients to the newly established & authenticated connection.
func RegisterRPCClients(conn *grpc.ClientConn) (err error) {

	// User
	UserRPC = dbpb.NewUserDBClient(conn)

	return
}
