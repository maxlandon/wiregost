package modules

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

// This file defines all the modules gRPC services (run, load, reload, edit, etc...)

// ModuleServer - Handles all module actions taken by a user's console.
// Assumingly, an instance is created each time a new console connects, regardless of if its user already has one opened.
type ModuleServer struct {
}

func (s *ModuleServer) RunModule(ctx context.Context, req *modulepb.ModuleActionRequest) (res *modulepb.ModuleAction, err error) {

	// We first get the identity of the console user + client

	// We add it to the ctx object

	// We run the module's run() function, passing the context object

	// We return the response of the module (blocking until received, but the module can send events async nonetheless)

	return nil, status.Errorf(codes.Unimplemented, "method RunModule not implemented")
}

func (s *ModuleServer) SetOption(context.Context, *modulepb.ModuleSetOptionRequest) (*modulepb.ModuleSetOption, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOption not implemented")
}

func (s *ModuleServer) UseModule(context.Context, *modulepb.StackUseRequest) (*modulepb.StackUse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseModule not implemented")
}

func (s *ModuleServer) PopModule(context.Context, *modulepb.StackPopRequest) (*modulepb.StackPop, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PopModule not implemented")
}
