package main

// wiregost-stack - A compiled stack of modules interacting with the Wiregost server.
// This is the answer to static compilation problems arising with exploit module development in Go.
// This executable is prepared, setup and compiled by Wiregost's module.Manager, at users' request.
// It asks for and shares state with the server, having access to many of its components (handlers,
// routing system, libraries, etc) through a set of gRPC clients. This ensures module logic can be
// based upon real-time state coming from the server.
//
// RESPONSIBILITIES:
// This executable is responsible for registering to all RPC components it needs for function.
// This is not the role of the Module Manager, which should limit itself to coordinate the system.
//
// NON-RESPONSIBILITIES:
// This binary is not in charge of verifying user permissions before they run their modules:
// It is the responsibility of module drivers that are instantiated when users run them.

func main() {

	// Get configuration and information

	// Setup logging

	// Init communications with server's module.Manager

	// Get module paths (public ones and private/user ones)

	// Load modules and init stack state

	// Subscribe to Wiregost events

	// Subscribe to Wiregost components gRPC (handlers, sessions, payloads, routing, etc...)

	// Notify back to server everything is set up.

	// Block
}
