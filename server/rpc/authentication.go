package rpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	client "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
)

type connectionServer struct {
	*client.UnimplementedConnectionRPCServer
}

func (c *connectionServer) Authenticate(context.Context, *client.AuthenticationRequest) (*client.Authentication, error) {

	// Check DB for users matching

	// Check password

	// If success {

	// Generate token

	// Send back user information

	// If error, send back not ok, empty user and empty token

	// Substract try out of 5

	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}

func (c *connectionServer) GetConnectionInfo(context.Context, *client.ConnectionInfoRequest) (*client.ConnectionInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConnectionInfo not implemented")
}

func (c *connectionServer) GetVersion(context.Context, *client.Empty) (*client.Version, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
