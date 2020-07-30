package core

import (
	"github.com/maxlandon/wiregost/client/context"
)

const compilerStr = "compiler"

// Compiler - This command switches to the compiler context.
type Compiler struct {
}

// Execute - Run
func (c *Compiler) Execute(args []string) (err error) {
	// Switch context. Prompt and context will automatically adapt.
	context.Context.Menu = context.CompilerMenu
	return
}
