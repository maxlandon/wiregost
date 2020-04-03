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
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/util"
)

// ParseEnvVariables - Replaces all environment variable tokens in the input with their values
func ParseEnvVariables(args []string) (processed []string, err error) {

	for _, arg := range args {

		// Anywhere a $ is assigned means there is an env variable
		if strings.Contains(arg, "$") {

			//Split in case env is embedded in path
			envArgs := strings.Split(arg, "/")

			// If its not a path
			if len(envArgs) == 1 {
				arg := envArgs[0]
				if strings.HasPrefix(arg, "$") && arg != "" && arg != "$" { // It is an env var.
					envVar := strings.TrimPrefix(arg, "$")
					val, ok := util.SystemEnv[envVar]
					if !ok {
						return processed, fmt.Errorf(ParserError+"Variable not in OS env: %s",
							tui.YELLOW, envVar, tui.RESET)
					}
					processed = append(processed, val)
					continue
				}
			}

			// If len of the env var split is > 1, its a path
			if len(envArgs) > 1 {
				var path []string

				for _, arg := range envArgs {
					// If item is an env var
					if strings.HasPrefix(arg, "$") && arg != "" && arg != "$" {
						envVar := strings.TrimPrefix(arg, "$")
						val, ok := util.SystemEnv[envVar]
						if !ok {
							return processed, fmt.Errorf(ParserError+"Variable not in OS env: %s%s%s",
								tui.YELLOW, envVar, tui.RESET)
						}
						path = append(path, val)
					} else {
						path = append(path, arg)
					}
				}
				// Make full processed path and return
				processed = append(processed, strings.Join(path, "/"))
			}
		} else {
			// Else, if arg is not an environment variable, return it as is
			processed = append(processed, arg)
		}

	}
	return
}
