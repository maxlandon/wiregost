package module

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

// RunModule - A user has requested to run one of the curent module's functions.
func (m *managers) RunModule(context.Context, *pb.RunRequest) (*pb.Run, error) {

	// Until now, everything that was happening with the current module was only
	// between the console and the stack binary. At the time of this call, both should
	// have an up-to-date (and identical) version of the Modulepb object: same info,
	// same commands, same options, same targets if specified, etc.

	// This function, therefore, must first ask the stack binary to send a copy of the concerned module.

	// Based on the module type, we should branch into several possibilities.

	// ANYWAY !!!!
	// This function should be in charge of asking the stack binary module to run its
	// main function: This will avoid us some tricky and ugly coding, like having a single
	// gRPC client function in one of the module's base types !
	// NOTE: Because all written modules will be subtypes of some module, there Run function
	// can only be called through an interface, remember that.

	return nil, status.Errorf(codes.Unimplemented, "method RunModule not implemented")
}

func (m *managers) GetInfo(context.Context, *pb.InfoRequest) (*pb.Info, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}

func (m *managers) GetOptions(context.Context, *pb.OptionsRequest) (*pb.Options, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOptions not implemented")
}

func (m *managers) SetOption(context.Context, *pb.SetOptionRequest) (*pb.Option, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOption not implemented")
}

func (m *managers) EditModule(context.Context, *pb.EditRequest) (*pb.Edit, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditModule not implemented")
}

// Module management methods (init new module, get/set module paths, etc).
