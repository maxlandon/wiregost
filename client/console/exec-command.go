// Wiregost - Post-Exploitation & Implant Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package console

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
	"github.com/jessevdk/go-flags"

	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/context"
	"github.com/maxlandon/wiregost/client/util"
)

// ExecuteCommand - Dispatches an input line to its appropriate command.
func (c *console) ExecuteCommand(args []string) (err error) {

	ctx := context.Context // The Console Context

	// We redirect the input to the appropriate parser, depending on the console menu.
	// The error returned might be several things, so we handle some cases later,
	// like special commands
	var parserErr error
	switch ctx.Menu {

	case context.MainMenu:
		_, parserErr = commands.Main.ParseArgs(args)

	case context.ModuleMenu:
		_, parserErr = commands.Module.ParseArgs(args)

	case context.CompilerMenu:
		_, parserErr = commands.Compiler.ParseArgs(args)

	case context.GhostMenu:
		_, parserErr = commands.Ghost.ParseArgs(args)
	default:
	}

	// All errors that might go out of parsers are handled here
	if parserErr != nil {
		err = c.HandleParserErrors(parserErr, args)
	}

	// END: Reset variables for command options (go-flags)

	return nil
}

// ExecuteSpecialCommand - Handles all commands not registered to command parsers.
func (c *console) ExecuteSpecialCommand(args []string) error {

	// Check context for availability
	switch context.Context.Menu {
	case context.MainMenu, context.ModuleMenu:
		switch args[0] {
		case "exit":
			c.Exit()
			return nil
		default:
			// Fallback: Use the system shell through the console
			return util.Shell(args)
		}
	}

	// We should not get here, so we print an error-like message
	fmt.Printf(CommandError+"Invalid command: %s%s \n", tui.YELLOW, args[0])

	return nil
}

// HandleParserErrors - The parsers may return various types of Errors, handle them in this function.
func (c *console) HandleParserErrors(in error, args []string) (err error) {

	// If there is an error, cast it to a parser error, else return
	var parserErr *flags.Error
	if in == nil {
		return
	}
	parserErr = in.(*flags.Error) // We convert to a flag error

	// Handle errors on a case-by-case basis -------------------

	// If command is not found, handle special
	if parserErr.Type == flags.ErrUnknownCommand {
		return c.ExecuteSpecialCommand(args)
	}
	return
}
