package module

import (
	"sync"

	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

var (
	// Managers - The module manager in Wiregost.
	Manager = &manager{}
)

// manager - A central element to Wiregost system: because modules can be recompiled and are fundamentally
// a different binary, they need to continuously share and have access to state of various components.
// On the other hand, console users need to interact with modules, which they do through this manager.
// There is one Manager instance running for each user, managing its own stack, its own drivers, etc...
//
// In sum, this manager is in charge of:
// - Starting, compiling, and reloading module stack binaries (must preserve cached stack list before shutdown)
// - Setting/registering all services (handlers, sessions, etc,) they need to interact with.
// - Interacting, to a lesser extent with drivers such as exploit_driver (furnishing him various components.)
// - Relaying user console commands and actions concerning modules (maybe not)
//
// The manager should provide easy access to modules for other server components such exploit_driver.
// This means the Module interface (behind which is hidden a specific type of module), can be passed
// around by this Manager.
type manager struct {

	// User stack (a list of module that have been ignited server-side, on a driver)
	// Having a stack server-side is useful, because we can send the list back to a
	// stack binary after it has restarted: it can then parse options/values for these
	// modules.

	// Client/User using this manager
	// We might need to track clients by some way or another, if we want to
	// push back various things to the right console. Or we pass state down
	// to drivers and modules

	// ModuleUI (consoles) gRPC server

	// Module Stack gRPC server
	*modulepb.UnimplementedManagerServer // Embedding this makes it a gRPC server

	// mutex
	mutex *sync.Mutex
}

func newManager() (m *manager) {
	return
}

// StartManager - Instantiates a module manager, dedicated to a single user.
// This function is called when a user is logged in (first console).
func StartManager() (m *manager) {

	// Check if user already has a manager running.

	// If yes, return the pointer to this instance.

	// If not, return a newly instantiated one and start it

	return
}

// Start - A function running concurrently, in which the Manager starts
// all its components (first and foremost, the stack binary).
func (m *manager) Start() (err error) {
	return
}

// Stack compilation/start/stop methods

// Stack connection/init/setup methods

// Module management methods (init new module, get/set module paths, etc).
