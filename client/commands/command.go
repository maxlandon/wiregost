package commands

import (
	"github.com/jessevdk/go-flags"
)

// PARSERS -----------------------------------------------------------------------------

// Each of these parsers is in charge of all commands for its respective menu context.
// These will be used by completion functions, between others.

// Main - Main context commands
var Main = flags.NewNamedParser("main", flags.None)

// Module - Module context commands
var Module = flags.NewNamedParser("module", flags.None)

// Ghost - Ghost context commands
var Ghost = flags.NewNamedParser("ghost", flags.None)

// Compiler - Compiler context commands
var Compiler = flags.NewNamedParser("compiler", flags.None)

// -------------------------------------------------------------------------------------

// InitParsers - This function is used to set all options for above Main, Module and Ghost parsers
func InitParsers() (err error) {

	// Main parser setup and registration
	err = BindMain()
	if err != nil {
		return err
	}
	// Compiler parser setup and registration
	err = BindCompiler()
	if err != nil {
		return err
	}

	return
}
