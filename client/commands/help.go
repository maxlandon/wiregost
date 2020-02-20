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

func RegisterHelpCommands() {

	// Command categories
	help := &Command{
		Name: "help",
		// Needed for completion
		SubCommands: []string{
			"core",
			"workspace",
			"hosts",
			"server",
			"user",
			"profiles",
			"jobs",
			"module",
			"stack",
			"sessions",
		},
		Args: []*CommandArg{
			&CommandArg{Name: "core"},
			&CommandArg{Name: "workspace"},
			&CommandArg{Name: "hosts"},
			&CommandArg{Name: "server"},
			&CommandArg{Name: "user"},
			&CommandArg{Name: "profiles"},
			&CommandArg{Name: "jobs"},
			&CommandArg{Name: "module"},
			&CommandArg{Name: "stack"},
			&CommandArg{Name: "sessions"},
		},
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				fmt.Printf(help.GetHelpFor("help"))
				fmt.Println()
			default:
				fmt.Println()
				fmt.Println(help.GetHelpFor(r.Args[0]))
			}
			return nil
		},
	}

	// Finally register commands
	AddCommand("main", help)
	AddCommand("module", help)
}
