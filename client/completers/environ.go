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

package completers

import (
	"strings"

	"github.com/lmorg/readline"
	"github.com/maxlandon/wiregost/client/util"
)

// CompleteEnvironmentVariables - Returns all environment variables as suggestions
func CompleteEnvironmentVariables(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	var suggestions []string
	listSuggestions := map[string]string{}
	args := strings.Split(string(line), " ")
	last := trimSpaceLeft([]rune(args[len(args)-1]))

	// Check if last input is made of several different variables
	allVars := strings.Split(string(last), "/")
	lastVar := allVars[len(allVars)-1]

	var evaluated = map[string]string{}

	for k, v := range util.SystemEnv {
		if strings.HasPrefix("$"+k, lastVar) {
			suggestions = append(suggestions, k[(len(lastVar)-1):]+"/")
			evaluated[k] = v
		}
	}

	return lastVar, suggestions, listSuggestions, readline.TabDisplayGrid
}
