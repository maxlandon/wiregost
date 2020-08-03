/*-------------------- Module package ---------------------

Checklist:
- Rename package name (package line)
- Rename `type Post struct` with module name
*/
package main // Rename this with module (directory) name

import (
	"github.com/maxlandon/wiregost/server/modules/post"
)

// Module - A Module Module (Change "Module")
type Module struct {
	// The module embeds (inherits from) a post.Module type, itself embedding
	// (inheriting from) a base.Module type. This makes it a valid module to Wiregost.
	// This gives your module full access to the server's sessions and their
	// functionalities, to the server's logging and console event system.
	// These underlying types also automatically handle various setup details so that
	// you don't have to deal with it yourself.
	*post.Post
}

// Run - Execute the module, with optional parameters for specifying actions to take or other details.
func (m *Module) Run(args []string) (result string, err error) {

	return
}
