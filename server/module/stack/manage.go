package stack

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

func (m *stacks) UseModule(context.Context, *pb.UseRequest) (*pb.Use, error) {

	// The manager checks his user stack: if not here continue, otherwise just return
	// this module's information

	// If current is empty { find module and put it as current }

	// If current not empty {}
	//      - Add module, catch ok, and error
	//      if err {
	//              return errors.New("Must be that module is incompatible with current")
	//      }
	//      if ok {
	//             // The problem is not an incompatibilty about platform/payload
	//             // Just that we are not allowed to load this module as a submodule,
	//             // so we use it as Current.
	//      }

	return nil, status.Errorf(codes.Unimplemented, "method UseModule not implemented")
}

func (m *stacks) PopModule(context.Context, *pb.PopRequest) (*pb.Pop, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PopModule not implemented")
}
func (m *stacks) PushModule(context.Context, *pb.PushRequest) (*pb.Push, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushModule not implemented")
}
func (m *stacks) ClearStack(context.Context, *pb.ClearRequest) (*pb.Clear, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearStack not implemented")
}
func (m *stacks) ReloadModule(in context.Context, req *pb.ReloadRequest) (*pb.Reload, error) {

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
