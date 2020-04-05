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

package util

import (
	"fmt"
	"strings"

	"github.com/evilsocket/islazy/tui"
)

var (
	// ParserError - Failed to parse some tokens in the input
	ParserError = fmt.Sprintf("%s[Parser Error]%s ", tui.RED, tui.RESET)
)

// SystemEnv - Contains all system environment variables
var SystemEnv = map[string]string{}

// ParseEnvVariables - Replaces all environment variable tokens in the input with their values
func ParseEnvVariables(args []string) (processed []string, err error) {

	for _, arg := range args {

		// Anywhere a $ is assigned means there is an env variable
		if strings.Contains(arg, "$") || strings.Contains(arg, "~") {

			//Split in case env is embedded in path
			envArgs := strings.Split(arg, "/")

			// If its not a path
			if len(envArgs) == 1 {
				processed = append(processed, handleSingleVar(arg))
			}

			// If len of the env var split is > 1, its a path
			if len(envArgs) > 1 {
				processed = append(processed, handleEmbeddedVar(arg))
			}
		} else if arg != "" && arg != " " {
			// Else, if arg is not an environment variable, return it as is
			processed = append(processed, arg)
		}

	}
	return
}

func handleSingleVar(arg string) string {
	if strings.HasPrefix(arg, "$") && arg != "" && arg != "$" { // It is an env var.
		envVar := strings.TrimPrefix(arg, "$")
		val, ok := SystemEnv[envVar]
		if !ok {
			return envVar
		}
		return val
	}
	if arg != "" && arg == "~" {
		return SystemEnv["HOME"]
	}

	return arg
}

func handleEmbeddedVar(arg string) string {

	envArgs := strings.Split(arg, "/")
	var path []string

	for _, arg := range envArgs {
		if strings.HasPrefix(arg, "$") && arg != "" && arg != "$" {
			envVar := strings.TrimPrefix(arg, "$")
			val, ok := SystemEnv[envVar]
			if !ok {
				// Err will be caught when command is ran anyway, or completion will stop...
				path = append(path, arg)
			}
			path = append(path, val)

		} else if arg != "" && arg == "~" {
			path = append(path, SystemEnv["HOME"])

		} else if arg != " " && arg != "" {
			path = append(path, arg)
		}
	}

	return strings.Join(path, "/")
}
