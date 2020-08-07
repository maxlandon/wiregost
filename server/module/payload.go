package module

import (
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

// Payload - A pseudo-module dedicated to payload generation and staging.
// It is not meant to be used and loaded on the Module Stack binary, because it would mean
// exposing the whole handler system with an RPC. Although, these modules will lie in the Module tree.
//
// A Payload module has no Run(cmd string, args []string) method, because it is not exploiting
// Instead, it sets up, through various functions, listeners and compilation details,
// before executing a command such as generate, to_handler, etc...
type Payload struct {
	// Base module. Allows us to us the module's logging system. The exploit driver
	// will decide, based on user choices, if the local debugging (happening in
	// most of these functions below) should be pushed back to the user console.
	*module

	// Ghost implant core (transport-agnostic) profile (format, arch, etc.)
	Info *ghostpb.Ghost

	Stage       []byte         // A ghost implant as bytes, for stagers
	StageConfig *ghostpb.Ghost // The configuration of the stage, for compabitibility checks
	IsCompiled  bool           // If this is true, that means a ghost implant is already ready as Bytes

	// Some security (limits, permissions, certificates, etc)

	// Some transport strings, passed at some point by a PayloadDriver
	Transports []string

	// Some working hours (Same, probably passed by a PayloadDriver, used by transports)
	// We might have an option to override payload working hours with transport ones, or the inverse

	// Some OS specifics
}

// Generate (with optional saving of the implant as file, or as bytes).

func (m *Payload) Init(meta *modulepb.Info) (err error) {

	// Parses the protobuf metadata (base module function)
	err = m.information(meta)

	// Checks various fields and adds some if needed. (Type-specific)
	meta.Type = modulepb.Type_PAYLOAD // Set module type

	// Setup logger here. It should be somehow hooked to the stack binary
	// gRPC connection with the Module Manager, with the ModuleLog stream.

	// The logger cannot be initiated from the module.Module base type,
	// because the method is unexported, and therefore cannot be called
	// from the stack binary.
	// Depending on which binary is calling this function we might pass
	// it a different logger (on the server a local one, on stack a remote)
	var isRemote bool

	if isRemote {
		// This instance is on the stack binary: start logger wired to gRPC.
	} else {
		// This instance is on the server: start a local logger.
	}

	return nil
}
