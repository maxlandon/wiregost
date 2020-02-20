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

package commands

import (
	"fmt"

	"github.com/maxlandon/wiregost/client/help"
)

func RegisterAgentHelpCommands() {

	// Command categories
	help := &Command{
		Name: "help",
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				// One page help with all commands
				fmt.Printf(help.GetHelpFor("complete-agent-help"))
				// Multi page command categories
				// fmt.Printf(help.GetHelpFor("agent-help"))
				fmt.Println()
			default:
				fmt.Println()
				fmt.Println(help.GetHelpFor(r.Args[0]))
			}
			return nil
		},
	}

	// Finally register commands
	AddCommand("agent", help)
}
