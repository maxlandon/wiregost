package module

// PayloadDriver - A driver for payload construction, setup, compilation and staging.
// It plays a role similar to ExploitDriver, thereby offering multiple possibilities
// for configuring per OS, per transport, and for potentially adding security modules.
// These modules would enhance the set of granularity/functionality/automation we'd like.
//
// It is somehow a bit different from an ExploitDriver:
// - A Payload module has no Run(cmd string, args []string) method, because it is not exploiting
//   Instead, it sets up, through various functions, listeners and compilation details,
//   before executing a command such as generate, to_handler, etc...
//
// A PayloadDriver instance is created each time a user loads up a payload module.
type PayloadDriver struct {
	// Manages some handlers (they pull their configuration out of the implant's transports)
	// Manages a Payload module
}

// Instantiation methods.
