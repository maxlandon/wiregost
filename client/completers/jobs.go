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
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/lmorg/readline"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/commands/jobs"
)

func CompleteJobIDs(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get last path
	splitLine := strings.Split(string(line), " ")
	last := splitLine[len(splitLine)-1]

	jobs := jobs.GetJobs(Context.Server.RPC)

	for _, job := range jobs.Active {
		suggestions = append(suggestions, strconv.Itoa(int(job.ID))[(len(last)):])
		listSuggestions[strconv.Itoa(int(job.ID))[(len(last)):]] = tui.RESET + job.Description + tui.DIM + " (:" + strconv.Itoa(int(job.Port)) + ")"
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayList
}
