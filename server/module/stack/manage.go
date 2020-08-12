package stack

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

// UseModule - The client requests to use a module. Depending on the current context, handle it.
func (m *stacks) UseModule(in context.Context, req *pb.UseRequest) (res *pb.Use, err error) {

	stack := m.GetUserStack(*req.Client) // User stack
	if stack == nil {
		err = errors.New("an error happened with your stack. It is nil.")
		return &pb.Use{Err: err.Error()}, err
	}

	// Find and instantiate the module. This will instantiate it.
	module, err := GetModule(req.Path)
	go module.SetupLog(false, req.Client, nil)

	// If err or module empty, send back an error
	if err != nil {
		return &pb.Use{Err: err.Error()}, nil
	} else if module == nil {
		err = fmt.Errorf("no module at path %s", req.Path)
		return &pb.Use{Err: err.Error()}, nil
	}

	// If current is empty, we simply add this module and return
	if stack.Current == nil {
		stack.Current = module
		res = &pb.Use{
			Loaded: module.ToProtobuf(),
		}
		return
	}

	// The current module is not empty, we let the current module handle
	// this request. It will determine if it accepts it or not.
	accepted, err := stack.Current.AddModule(module)

	// If there is an error, the module might have accepted it
	// but didn't for a reason (generally, submodule compatibility).
	// We return the error immediately, so the user can see for himself.
	if err != nil {
		return &pb.Use{Err: err.Error()}, nil
	}

	// Here, module does not support this module as a "submodule".
	// Therefore we must devise what to do with it. See later.
	if accepted == false {
		stack.Current = module
		return &pb.Use{Loaded: module.ToProtobuf()}, nil

	}

	// END: We come here, so that means we added a module somewhere.
	// We must check if this module has a peer with a Run() function
	// on the user's stack binary. If yes, we do the necessary.

	// Change this
	return &pb.Use{Err: "Should not have got here"}, nil
}

// PushModule - The last module pushed on the stack popped back as the current module.
func (m *stacks) PopModule(in context.Context, req *pb.PopRequest) (res *pb.Pop, err error) {

	stack := m.GetUserStack(*req.Client) // User stack

	// Empty stack
	if len(stack.Stack) == 0 {
		err = fmt.Errorf("the module stack is empty")
		return &pb.Pop{Err: err.Error()}, nil
	}

	stack.Current = stack.Stack[len(stack.Stack)-1] // Set the module as current.
	stack.Stack = stack.Stack[:len(stack.Stack)-1]  // Pop its reference on the stack.

	res = &pb.Pop{
		Next: stack.Current.ToProtobuf(),
	}

	return
}

// PushModule - The module currently loaded by the user is pushed to the user's Stack.
func (m *stacks) PushModule(in context.Context, req *pb.PushRequest) (res *pb.Push, err error) {

	stack := m.GetUserStack(*req.Client) // User stack

	// We just put it at end (top) of the stack
	stack.Stack = append(stack.Stack, stack.Current)

	return &pb.Push{}, nil
}
func (m *stacks) ClearStack(in context.Context, req *pb.ClearRequest) (res *pb.Clear, err error) {

	stack := m.GetUserStack(*req.Client) // User stack

	stack.Stack = []Module{} // Hope we took care of background jobs first.

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
