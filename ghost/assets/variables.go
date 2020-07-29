package assets

// This file is used to gather ALL compile-time values for a ghost implant.
// No compile-time variable should be located outside this file.
// The choice for a single file is to simplify as much as possible the -ldflags
// usage: we won't have to input dozens of different paths to all packages,
// changing accross commits, versions, etc...
//
// Comments should always be present next to the VariableName, like:
// => "true"
// => "value:format,pattern"
// This will help those having problems with this.
var (
	// LOGGING ------------------------------------------------------------------------------

	// DebugLocal - Local, command-line debugging
	DebugLocal string // => "true"
	// DebugRemote - All logs are sent back to the server. Many timings/strategies possible
	DebugRemote string // => "true"

	// OWNERSHIP & PERMISSIONS --------------------------------------------------------------

	// RestrictToOwner -
	RestrictToOwner string

	// TRANSPORTS ---------------------------------------------------------------------------

	// StartupTransport - The first transport to use at startup
	StartupTransport string // "protocol://host:port"
	// OtherTransports - All other transports to preconfigure
	OtherTransports string // "protocol://host:port,protocol://host:port,protocol://host:port,"

	// TRANSPORT SECURITY -------------------------------------------------------------------
)
