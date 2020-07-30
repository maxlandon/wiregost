package compiler

import (
	"github.com/jessevdk/go-flags"
	"github.com/maxlandon/wiregost/client/context"
	help "github.com/maxlandon/wiregost/client/help/main/core"
)

const compilerExitStr = "exit"

// Exit - This command switches to the compiler context.
type Exit struct {
}

// Execute - Run
func (c *Exit) Execute(args []string) (err error) {

	// Switch context
	context.Context.Menu = context.MainMenu

	return
}

// InitCompilerExit - Register compiler command
func InitCompilerExit(parser *flags.Parser) (err error) {
	// Add
	// _, err = parser.AddCommand(compilerStr, help.CompilerShort, help.CompilerLong, &comp)
	exit, err := parser.AddCommand(compilerExitStr, help.CompilerExitShort, help.CompilerExitLong, &Exit{})

	exit.Aliases = []string{"main", "back"}

	return
}
