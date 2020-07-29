package generate

import (
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

// GhostImplant - Root function for compiling a ghost implant.
// It takes care of parsing fields for each "category", setup the according
// details and values to be compiled, and then builds the implant.
func GhostImplant(prof serverpb.GhostBuild) (err error) {

	compilerLog.Infof("Starting compilation for an implant")

	// Base compilation strings

	// Compiler & Server version details

	// Architecture details

	// OS details

	// Binary Format details

	// Transports

	// Transport Security

	// Owner & Core Permissions

	// Routes & Routing Permissions

	// Add job to server

	// Compilation

	// If no errors, create implant directory, and add files if necessary.

	// Save file and/or send implant back to console (return it)

	return
}
