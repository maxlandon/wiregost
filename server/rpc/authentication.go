package rpc

import (
	"context"

	"google.golang.org/grpc"

	db "github.com/maxlandon/wiregost/db/client"
	client "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
)

type connectionServer struct {
	*client.UnimplementedConnectionRPCServer
}

func (c *connectionServer) Authenticate(ctx context.Context, req *client.AuthenticationRequest) (*client.Authentication, error) {

	// Check DB for users matching
	dbRes, _ := db.Users.GetUsers(ctx, &dbpb.User{Name: req.Username}, grpc.EmptyCallOption{})
	if dbRes == nil {
		return &client.Authentication{}, nil
	}

	// Check password
	if string(dbRes.Users[0].Password) != req.Password {
		return &client.Authentication{}, nil
	}

	// If success, fill user info
	res := &client.Authentication{}
	res.Success = true
	res.Client = &client.Client{}
	res.Client.User = dbRes.Users[0]

	// Generate token
	res.Client.Token = "test"

	// Send back user information

	// If error, send back not ok, empty user and empty token

	// Substract try out of 5

	return res, nil
	// return &client.Authentication{}, nil
	// return nil, nil
	// return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}

func (c *connectionServer) GetConnectionInfo(context.Context, *client.ConnectionInfoRequest) (*client.ConnectionInfo, error) {
	return &client.ConnectionInfo{}, nil
	// return nil, status.Errorf(codes.Unimplemented, "method GetConnectionInfo not implemented")
}

func (c *connectionServer) GetVersion(context.Context, *client.Empty) (*client.Version, error) {
	ver := &client.Version{
		ClientMajor:     "1",
		ClientMinor:     "0",
		ClientPatch:     "0",
		ServerMajor:     "1",
		ServerMinor:     "0",
		ServerPatch:     "0",
		ServerCommitTag: "fhlh83hkllfd8kjld7321hf908Ofwhw",
	}
	return ver, nil
	// return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
