package stack

import (
	"sync"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

var (
	// Stacks - The module Stacks in Wiregost. There is one Stack instance per user
	// and it handles everything for all connected users and their consoles.
	Stacks = &stacks{
		Active: &map[*dbpb.User]*stack{},
		mutex:  &sync.Mutex{},
	}
)

type stacks struct {
	Active *map[*dbpb.User]*stack

	// ModuleUI (consoles) gRPC server
	// There is one gRPC server instance for Stackss. Each request to one of
	// the services contains a Client PB object, for dispatching to the good stack.
	*modulepb.UnimplementedStackServer // Embedding this makes it a gRPC server

	// Should pass on a logger from when the user connected.

	mutex *sync.Mutex
}

// Add - Adds a module stack for the user that just registered.
func (s *stacks) Add(client *clientpb.Client) (stack *stack, err error) {
	s.mutex.Lock()
	stack = newStack(client)
	(*s.Active)[client.User] = stack
	s.mutex.Unlock()
	return
}

func GetUserStack(client *clientpb.Client) (s *stack) {
	return (*Stacks.Active)[client.User]
}

// stack - A central element to Wiregost system: because modules can be recompiled and are fundamentally
// a different binary, they need to continuously share and have access to state of various components.
// On the other hand, console users need to interact with modules, which they do through this stack.
// There is one Manager instance running for each user, managing its own stack, its own drivers, etc...
//
// In sum, this stack is in charge of:
// - Starting, compiling, and reloading module stack binaries (must preserve cached stack list before shutdown)
// - Setting/registering all services (handlers, sessions, etc,) they need to interact with.
// - Interacting, to a lesser extent with drivers such as exploit_driver (furnishing him various components.)
// - Relaying user console commands and actions concerning modules (maybe not)
//
// The stack should provide easy access to modules for other server components such exploit_driver.
// This means the Module interface (behind which is hidden a specific type of module), can be passed
// around by this Manager.
type stack struct {

	// Current - A module loaded "on a cache": the next loaded module on the console
	// whill erase this one purely and simply. Users can push it on the stack.
	Current Module

	// User stack (a list of module that have been ignited server-side, on a driver).
	// Some modules loaded on this Stack might have stack binary counterparts: when
	// a stack binary is recompiled and restarted, we send back these modules' info
	// to the stack, which is then ready again.
	// Module is an interface for management of modules no matter their subtypes.
	Stack []Module

	// Client/User using this stack. This Client information is passed
	// down to modules and sessions, who need it for various things.
	User *dbpb.User

	// Module Stack gRPC Client
	Modules modulepb.StackClient

	// mutex
	mutex *sync.Mutex
}

func newStack(client *clientpb.Client) (m *stack) {
	return &stack{
		Current: nil,         // No module is currently loaded
		Stack:   []Module{},  // The list of modules on stack is empty
		User:    client.User, // The stack is being assigned to a User
		Modules: nil,         // There is no stack binary started yet.
	}
}

// AssignStack - Instantiates a module Stacks, dedicated to a single user.
// This function is called when a user is logged in (first console).
func AssignStack(client *clientpb.Client) (m *stack) {

	// Check if user already has a Stacks running.
	if stack, found := (*Stacks.Active)[client.User]; found {
		// If yes, return the pointer to this instance.
		return stack
	}

	// If not, return a newly instantiated one and start it
	m, _ = Stacks.Add(client)

	return
}

// Start - A function running concurrently, in which the Manager starts
// all its components (first and foremost, the stack binary).
func (m *stack) Start() (err error) {

	// Check stack binary is at specified path and compiled up-to-date

	// Start the stack binary

	// Connect to it with a ClientConn, and register the StackClient to it.

	// Additional checks if needed

	return
}
