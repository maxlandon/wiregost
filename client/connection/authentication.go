package connection

import (
	"google.golang.org/grpc"

	"github.com/google/uuid"
	client "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
)

// Authenticate - Perform full authentication process with server
func Authenticate(conn *grpc.ClientConn) (cli client.ConnectionRPCClient, user *dbpb.User, clientID uuid.UUID) {

	// Register ConnectionRPC client to connection
	cli = client.NewConnectionRPCClient(conn)

	// Prompt, store and send password (as a hash)

	// Send authentication request

	// Check answer, with success and token

	// If error, try again (five tries)

	// If success, store token, return

	// If failure, print error and exit program

	return
}
