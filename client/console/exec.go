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
	"github.com/maxlandon/wiregost/client/commands"
)

func (c *Console) ExecCommand(args []string) error {
	if len(args) < 1 {
		return nil
	}

	// 1) Ask parser for command
	cmd := commands.CommandParser.Find(args[0])
	if cmd == nil {
		return c.handleSpecialCommands(args)
	}

	// 2) If command is found, check for context
	_, err := commands.CommandParser.ParseArgs(args)
	if err != nil {
		return err // Not printed currently
	}

	return nil
}

func (c *Console) handleSpecialCommands(args []string) error {

	switch c.menu {
	// Check context for availability
	case commands.MAIN_CONTEXT:
		switch args[0] {
		// ! - Use the system shell through the console
		case "!":
			return commands.Shell(args)
		case "exit":
			c.exit()
		}
	}

	fmt.Printf(CommandError+"Invalid command: %s%s \n", tui.YELLOW, args[0])
	return nil
	// return fmt.Errorf(CommandError+"Invalid command: %s%s +\n", tui.YELLOW, args[0])
}
