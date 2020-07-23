package generic

import corepb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost/core"

// Execute - This interface contains all methods for executing code on the target.
// Some of these methods might not be implemented by the `generic.Ghost` type,
// in which case, they will be implemented by a Windows/Linux/Darwin Ghost subtype.
type Execute interface {
	Execute(path string, args []string) (exec *corepb.Execute)
}

// Execute - Execute a program on the target
func (g *Ghost) Execute(path string, args []string) (exec *corepb.Execute) {
	return
}
