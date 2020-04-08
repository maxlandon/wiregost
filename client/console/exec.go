// Wiregost - Golang Exploitation Framework
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
	"github.com/maxlandon/wiregost/client/util"
)

func (c *Console) ExecCommand(args []string) error {
	if len(args) < 1 {
		return nil
	}

	// 1) Check context
	var cmd *flags.Command
	cmds := commands.CommandsByContext()
	for _, c := range cmds {
		if c.Name == args[0] {
			cmd = c
		}
	}

	if cmd == nil {
		return c.handleSpecialCommands(args)
	}

	_, err := commands.CommandParser.ParseArgs(args)
	if err != nil {
		return err // Not printed currently
	}

	return nil
}

// handleSpecialCommands - Handles all commands not registered to the parser
func (c *Console) handleSpecialCommands(args []string) error {

	switch c.menu {
	// Check context for availability
	case commands.MAIN_CONTEXT, commands.MODULE_CONTEXT:
		switch args[0] {
		case "exit":
			c.exit()
			return nil
			// Fallback: Use the system shell through the console
		default:
			return util.Shell(args)
		}
	}

	fmt.Printf(CommandError+"Invalid command: %s%s \n", tui.YELLOW, args[0])
	return nil
}
