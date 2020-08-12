package clients

import (
	"context"
	"strconv"

	"google.golang.org/grpc"

	db "github.com/maxlandon/wiregost/db/client"
	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/version"
)

type connectionServer struct {
	*clientpb.UnimplementedConnectionRPCServer
}

func (c *connectionServer) Authenticate(ctx context.Context, req *clientpb.AuthenticationRequest) (*clientpb.Authentication, error) {

	// If already 5 attempts, do not go further
	if ((*Consoles.Unauthenticated)[req.MD.ID] != nil) && ((*Consoles.ClientAttempts)[req.MD.ID] >= 5) {
		return &clientpb.Authentication{}, nil
	}

	// Add client to clients map (it is temporary)
	temp := &clientpb.Client{ID: req.MD.ID}
	Consoles.AddClient(*temp)

	// Check DB for users matching
	dbRes, _ := db.Users.GetUsers(ctx, &dbpb.User{Name: req.Username}, grpc.EmptyCallOption{})

	// If no one found, remove client & increase counter (the counter will leave a trace of the token as key)
	if dbRes == nil {
		Consoles.IncrementClientAttempts(temp.ID)
		Consoles.RemoveClient(temp.ID)

		return &clientpb.Authentication{}, nil
	}

	// If password wrong, send back not ok, empty user and empty token
	if string(dbRes.Users[0].Password) != req.Password {
		Consoles.IncrementClientAttempts(temp.ID)
		Consoles.RemoveClient(temp.ID)

		return &clientpb.Authentication{}, nil
	}

	// If success, fill user info, and send back token
	res := &clientpb.Authentication{}
	res.Success = true
	res.Client = temp
	res.Client.ID = req.MD.ID

	res.Client.User = dbRes.Users[0]
	res.Client.User.Online = true

	// We confirm the client, which will register it to a module stack as well
	Consoles.ConfirmClient(*res.Client)

	return res, nil
}

func (c *connectionServer) GetConnectionInfo(context.Context, *clientpb.ConnectionInfoRequest) (*clientpb.ConnectionInfo, error) {

	info := &clientpb.ConnectionInfo{
		DBHost: assets.ServerConfiguration.DatabaseRPCHost,
		DBPort: int32(assets.ServerConfiguration.DatabaseRPCPort),
		Jobs:   3,
		Ghosts: 13,
		ConsoleConfig: &clientpb.ConsoleConfig{
			MainPrompt:    "",
			ImplantPrompt: "",
		},
	}

	return info, nil
}

func (c *connectionServer) GetVersion(context.Context, *clientpb.Empty) (*clientpb.Version, error) {

	ver := &clientpb.Version{
		ClientMajor:     "1",
		ClientMinor:     "0",
		ClientPatch:     "0",
		ServerMajor:     strconv.Itoa(version.SemanticVersion()[0]),
		ServerMinor:     strconv.Itoa(version.SemanticVersion()[1]),
		ServerPatch:     strconv.Itoa(version.SemanticVersion()[2]),
		ServerCommitTag: version.GitCommit,
	}
	return ver, nil
}
