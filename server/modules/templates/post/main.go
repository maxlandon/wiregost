package main // Rename this with module (directory) name

import (
	"github.com/maxlandon/wiregost/server/modules/post"
)

// Module - Wiregost post-exploitation module.
type Module struct {
	// The module embeds (inherits from) a post.Module type, itself embedding (inheriting from) a
	// base.Module type. This makes it a valid module to Wiregost, and gives you access to :
	//
	// The base.Module type (that you don't see here) gives you access to the server's logging and
	// event system, provides various utility methods for querying options, and performs automatic and
	// serialized safety checks (see the function PreRunChecks() call below in the Run() method).
	//
	// This post.Post type gives you access to a Session object, with itself provides all methods
	// and functionality of a ghost implant, through various interfaces: many concrete types, generally
	// depending on the target's OS, will hide behind and be used through the Session object.
	*post.Post
}

// Run - Execute the module, with optional command for specifying actions to take, and parameters for
// handling or other details. Processing these args is left at the author's discretion.
func (m *Module) Run(cmd string, args []string) (result string, err error) {

	// Switch on console provided command: it has already been checked valid by the server.
	switch cmd {
	case "check":
		// We might have a function/method Check() handling this case and its subcases.
	case "exploit":
		// We might have a function/method Exploit() performing the exploit/action.
	}

	// We might as well have to specific subcommand for this module, and directly hold the logic in this function.

	return "Status message for end of module execution", nil
}
