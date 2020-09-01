package stack

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/module"
)

// Module - The module interface is used by a stack to interact with a module type embedding
// a module.Module type that implements this interface. Some of the interface functions
// are reimplemented by these module subtypes, for adapting workflow to their task.
type Module interface {
	// Run - All modules have, at least from their module.Module base type, a Run() method.
	Run(cmd string, args []string) (result string, err error)
	// ToProtobuf - The module.Module can push info, opts & commands with it.
	ToProtobuf() (mod *pb.Module)
	// SetOption - The module has the ability to set its options.
	SetOption(opt *pb.Option) (err error)
	// AddModule - Some modules may be able to combine with other module types.
	// This method leaves them with how to handle their babies.
	AddModule(m module.Module) (ok bool, err error)
	// SetupLog - Called by module Stacks (server-side and stack-side), for
	SetupLog(remote bool, cli *clientpb.Client, rpc serverpb.EventsClient) (logger *logrus.Entry)
}

// RunModule - A user has requested to run one of the curent module's functions.
func (m *stacks) RunModule(context.Context, *pb.RunRequest) (*pb.Run, error) {

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

	// As well, obviously, if the module run command does not need a Transport/Payload
	// because it is not remote, do not use a Driver.Run() method

	return nil, status.Errorf(codes.Unimplemented, "method RunModule not implemented")
}

func (m *stacks) GetInfo(context.Context, *pb.InfoRequest) (*pb.Info, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}

func (m *stacks) GetOptions(context.Context, *pb.OptionsRequest) (*pb.Options, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOptions not implemented")
}

func (m *stacks) SetOption(context.Context, *pb.SetOptionRequest) (*pb.Option, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOption not implemented")
}

func (m *stacks) EditModule(context.Context, *pb.EditRequest) (*pb.Edit, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditModule not implemented")
}

// Module management methods (init new module, get/set module paths, etc).

// GetModule - This function finds a module by path (doing all the processing and checking
// if needed), instantiates it and returns it to the stack.
func GetModule(path string) (m Module, err error) {

	return
}
