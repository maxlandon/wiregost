package rpc

import (
	"context"

	"google.golang.org/grpc"

	db "github.com/maxlandon/wiregost/db/client"
	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	"github.com/maxlandon/wiregost/server/assets"
)

type connectionServer struct {
	*clientpb.UnimplementedConnectionRPCServer
}

func (c *connectionServer) Authenticate(ctx context.Context, req *clientpb.AuthenticationRequest) (*clientpb.Authentication, error) {

	// If already 5 attempts, do not go further
	if ((*Clients.Unauthenticated)[req.MD.Token] != nil) && ((*Clients.ClientAttempts)[req.MD.Token] >= 5) {
		return &clientpb.Authentication{}, nil
	}

	// Add client to clients map (it is temporary)
	temp := &clientpb.Client{Token: req.MD.Token}
	Clients.AddClient(*temp)

	// Check DB for users matching
	dbRes, _ := db.Users.GetUsers(ctx, &dbpb.User{Name: req.Username}, grpc.EmptyCallOption{})

	// If no one found, remove client & increase counter (the counter will leave a trace of the token as key)
	if dbRes == nil {

		Clients.IncrementClientAttempts(temp.Token)
		Clients.RemoveClient(temp.Token)

		return &clientpb.Authentication{}, nil
	}

	// If password wrong, send back not ok, empty user and empty token
	if string(dbRes.Users[0].Password) != req.Password {

		Clients.IncrementClientAttempts(temp.Token)
		Clients.RemoveClient(temp.Token)

		return &clientpb.Authentication{}, nil
	}

	// If success, fill user info, and send back token
	res := &clientpb.Authentication{}
	res.Success = true
	res.Client = temp
	res.Client.Token = req.MD.Token

	res.Client.User = dbRes.Users[0]
	res.Client.User.Online = true

	return res, nil
}

func (c *connectionServer) GetConnectionInfo(context.Context, *clientpb.ConnectionInfoRequest) (*clientpb.ConnectionInfo, error) {

	info := &clientpb.ConnectionInfo{
		DBHost: assets.ServerConfiguration.DatabaseRPCHost,
		DBPort: int32(assets.ServerConfiguration.DatabaseRPCPort),
		Jobs:   3,
		Ghosts: 13,
	}

	return info, nil
}

func (c *connectionServer) GetVersion(context.Context, *clientpb.Empty) (*clientpb.Version, error) {

	ver := &clientpb.Version{
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
