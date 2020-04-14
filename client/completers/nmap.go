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

// import (
//         "fmt"
//         "strings"
//
//         "github.com/evilsocket/islazy/tui"
//         "github.com/maxlandon/wiregost/client/commands"
// )
//
// type nmapCompleter struct {
//         Command *commands.Command
// }
//
// // Do is the completion function triggered at each line
// func (mc *nmapCompleter) Do(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {
//
//         splitLine := strings.Split(string(line), " ")
//         line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))
//
//         responses := []string{"response", "probe", "scan"}
//
//         for _, mod := range responses {
//                 search := mod
//                 if !hasPrefix(line, []rune(search)) {
//                         fmt.Printf("\n %sOptions%s\n", tui.YELLOW, tui.RESET)
//                         sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
//                         options = append(options, sLine...)
//                         offset = sOffset
//                 }
//         }
//         return options, offset
// }
