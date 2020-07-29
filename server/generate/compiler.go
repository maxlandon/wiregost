package generate

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/log"
)

var (
	compilerLog = log.ServerLogger("compilation", "generate", "rpc")
)

// Compiler - Manages all compilations in Wiregost
type Compiler struct {
	*serverpb.UnimplementedCompilerServer
}

// GetGhostProfiles - Returns all saved ghost implants build profiles.
func (c *Compiler) GetGhostProfiles(ctx context.Context, req *serverpb.GhostProfilesRequest) (res *serverpb.GhostProfiles, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGhostProfiles not implemented")
}

// CompileGhost - Request to compile a ghost implant.
func (c *Compiler) CompileGhost(ctx context.Context, req *serverpb.BuildGhostRequest) (res *serverpb.BuildGhost, err error) {

	// Concurrently compile implant
	go GhostImplant(*req.Profile)

	return
}

// CompileConsole - Request to compile a user console.
func (c *Compiler) CompileConsole(ctx context.Context, req *serverpb.BuildConsoleRequest) (res *serverpb.BuildConsole, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompileConsole not implemented")
}
