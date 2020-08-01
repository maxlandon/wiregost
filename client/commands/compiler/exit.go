package compiler

import (
	"github.com/maxlandon/wiregost/client/context"
)

// Enter - This command switches to the compiler context.
type Enter struct {
}

// Execute - Run
func (c *Enter) Execute(args []string) (err error) {

	// Switch context
	context.Context.Menu = context.CompilerMenu

	return
}

// Exit - This command switches back to main/module/ghost context.
type Exit struct {
}

// Execute - Run
func (c *Exit) Execute(args []string) (err error) {

	// Switch context
	context.Context.Menu = context.MainMenu

	return
}
