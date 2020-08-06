package module

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

func (m *manager) PopModule(context.Context, *pb.PopRequest) (*pb.Pop, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PopModule not implemented")
}
func (m *manager) PushModule(context.Context, *pb.PushRequest) (*pb.Push, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushModule not implemented")
}
func (m *manager) ClearStack(context.Context, *pb.ClearRequest) (*pb.Clear, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearStack not implemented")
}
func (m *manager) ReloadModule(context.Context, *pb.ReloadRequest) (*pb.Reload, error) {

	// recompile

	// RESTART
	// If we have modules currently running on the stack binary, we can start another one
	// and wire everything needed to this new stack. This will prevent things from screwing
	// while allowing users to work on separate things at once.

	// Connect

	// Setup

	// Confirm

	return nil, status.Errorf(codes.Unimplemented, "method ReloadModule not implemented")
}

// Stack compilation/start/stop methods

// Stack connection/init/setup methods
