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
	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/util"
)

// AutoCompleter is the autocompletion engine
type OptionCompleter struct {
	Command *commands.Command
}

// Do is the completion function triggered at each line
func (oc *OptionCompleter) Do(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {

	switch *ctx.MenuContext {
	case "module":
		for _, v := range util.SortOptionKeys(ctx.Module.Options) {
			search := v
			if !hasPrefix(line, []rune(search)) {
				sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}
		return
	}

	return options, offset
}
