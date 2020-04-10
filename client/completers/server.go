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
	"fmt"
	"strings"

	"github.com/lmorg/readline"
	"github.com/maxlandon/wiregost/client/assets"
)

func CompleteServer(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get last path
	splitLine := strings.Split(string(line), " ")
	last := splitLine[len(splitLine)-1]

	// Get configs
	configs := assets.GetConfigs()

	for _, c := range configs {
		conf := fmt.Sprintf("%s@%s:%d", c.User, c.LHost, c.LPort)
		if strings.HasPrefix(conf, string(last)) {
			suggestions = append(suggestions, conf[len(last):])
		}
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
}
