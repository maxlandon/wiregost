package generic

import corepb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost/core"

type Proc interface {
	Dump(procID uint32) (dump corepb.ProcessDump)
}

// Dump - Dump memory from a target process
func (g *Ghost) Dump(procID uint32) (dump corepb.ProcessDump) {
	return
}
