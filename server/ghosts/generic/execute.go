package generic

import corepb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost/core"

type Execute interface{}

// Execute - Execute a program on the target
func (g *Ghost) Execute(path string, args []string) (exec *corepb.Execute) {
	return
}
