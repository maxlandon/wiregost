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

	"github.com/maxlandon/wiregost/client/context"
	"github.com/maxlandon/wiregost/client/util"
)

// ExecuteCommand - Dispatches an input line to its appropriate command.
func (c *console) ExecuteCommand(args []string) error {

	var cmd *flags.Command // Command detected and stored

	// 2) If command is not found, handle special
	if cmd == nil {
		return c.ExecuteSpecialCommand(args)
	}

	// END: Reset variables for command options (go-flags)

	return nil
}

// ExecuteSpecialCommand - Handles all commands not registered to command parsers.
func (c *console) ExecuteSpecialCommand(args []string) error {

	switch context.Context.Menu {
	// Check context for availability
	case context.MainMenu, context.ModuleMenu:
		switch args[0] {
		case "exit":
			c.Exit()
			return nil
			// Fallback: Use the system shell through the console
		default:
			return util.Shell(args)
		}
	}

	fmt.Printf(CommandError+"Invalid command: %s%s \n", tui.YELLOW, args[0])

	return nil
}
