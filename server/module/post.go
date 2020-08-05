package module

import (
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	"github.com/maxlandon/wiregost/server/ghosts"
)

// Post - A module dedicated to post-exploitation activities on a remote target, via a session instance.
type Post struct {
	*module               // Base module
	Session *ghosts.Ghost // Session is an interface accepting different implants
}

// NewPost - Instantiates a new post, and handles base module instantiation
func NewPost() (post *Post) {
	return
}

// Init - Module initialization process: parses metadata/information, register options.
// This function is called by the Stack binary, which then sends back module information
// as Protobuf to the server.
func (m *Post) Init(meta *modulepb.Info) (err error) {

	// Checks various fields and adds some if needed. (Type-specific)
	meta.Type = modulepb.Type_POST // Set module type

	// Parses the protobuf metadata (base module function)
	err = m.information(meta)

	return nil
}

// GetSession - Returns the Session corresponding to the Post "Session" option.
func (m *Post) GetSession(id uint32) (err error) {

	// requested := ghosts.Ghosts.Get(id)
	// ghost := requested.Core.Info()

	// We check permissions here and now, as we cannot pass
	// the module's context to each implant method call in module
	// Any calls to implant RPC stubs will trigger permission checks anyway.
	// _, err = security.CheckCorePermissions(ghost, m.Client.User)
	// if err != nil {
	//         return
	// }
	// m.Session = requested

	return
}
