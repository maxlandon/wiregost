package commands

import "github.com/jessevdk/go-flags"

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

// This function is used to set all options for above Main, Module and Ghost parsers
func InitParsers() (err error) {

	// For each parser:

	// 1) Set long/short descriptions (what is this menu), Usage (all commands and/or categories).

	// 2) Set unknown option handler

	// Add groups

	return
}

// Command - A command (that may have subcommands) or a subcommand dedicated to a single field/area/function/task.
// type Command struct {
//         Name string       // Command name to input
//         Help string       // Help/Usage/Doc
//         Sub  []string     // Subcommand strings (can be found later)
//         Opts []*Option    // Command options
//         Args []*Argument  // Command arguments
//         Run  func() error // This function runs the command
// }

// Option - A dash ( --option / -o ) option for a command
// type Option struct {
//         Short       string
//         Long        string
//         Type        string
//         Description string
//         Required    bool
//         Length      int
// }

// Argument - A command/subcommand argument, AS OPPOSED TO A COMMAND DASH (--) OPTION.
// type Argument struct{}
