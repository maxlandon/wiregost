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
//         "strings"
//
//         "github.com/lmorg/readline"
//
//         "github.com/maxlandon/wiregost/client/commands"
//         "github.com/maxlandon/wiregost/data-service/remote"
// )
//
// type workspaceCompleter struct {
//         Command *commands.Command
// }
//
// func completeWorkspaces(cmd *commands.Command, line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {
//         // Completions
//         var suggestions []string
//         listSuggestions := map[string]string{}
//
//         // Get last path
//         splitLine := strings.Split(string(line), " ")
//         last := trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))
//
//         // Get workspaces
//         workspaces, _ := remote.Workspaces(nil)
//
//         // Get completions
//         for _, ws := range workspaces {
//                 if strings.HasPrefix(ws.Name, string(last)) {
//                         suggestions = append(suggestions, ws.Name[len(last):])
//                 }
//         }
//
//         return string(line[:pos]), suggestions, listSuggestions, readline.TabDisplayGrid
// }
//
// // Do is the completion function triggered at each line
// func (wc *workspaceCompleter) Do(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {
//
//         // Complete command args
//         splitLine := strings.Split(string(line), " ")
//         line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))
//
//         switch splitLine[0] {
//         // Provide only workspace names
//         case "switch":
//                 workspaces, _ := remote.Workspaces(nil)
//                 for _, w := range workspaces {
//                         search := w.Name
//                         if !hasPrefix(line, []rune(search)) {
//                                 sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
//                                 options = append(options, sLine...)
//                                 offset = sOffset
//                         }
//                 }
//                 return
//         case "delete":
//                 workspaces, _ := remote.Workspaces(nil)
//                 for _, w := range workspaces {
//                         search := w.Name
//                         if !hasPrefix(line, []rune(search)) {
//                                 sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
//                                 options = append(options, sLine...)
//                                 offset = sOffset
//                         }
//                 }
//                 return
//         case "add":
//                 return
//                 // Provide all arguments
//         case "update":
//                 for _, arg := range wc.Command.Args {
//                         search := arg.Name
//                         if !hasPrefix(line, []rune(search)) {
//                                 sLine, sOffset := doInternal(line, pos, len(line), []rune(search+"="))
//                                 options = append(options, sLine...)
//                                 offset = sOffset
//                         } else {
//                                 words := strings.Split(string(line), "=")
//                                 argInput := lastString(words)
//                                 if arg.Type == "boolean" {
//                                         for _, search := range []string{"true ", "false "} {
//                                                 offset = 0
//                                                 if strings.HasPrefix(search, argInput) {
//                                                         options = append(options, []rune(search[len(argInput):]))
//                                                         offset = len(argInput)
//                                                 }
//                                         }
//                                         return
//                                 }
//                                 if arg.Type == "string" && arg.Name == "name" {
//
//                                         workspaces, _ := remote.Workspaces(nil)
//                                         names := []string{}
//                                         for _, w := range workspaces {
//                                                 offset = 0
//                                                 names = append(names, w.Name)
//                                                 if strings.HasPrefix(w.Name, argInput) {
//                                                         options = append(options, []rune(w.Name[len(argInput):]+" "))
//                                                         offset = len(argInput)
//                                                 }
//                                         }
//                                         return
//                                 }
//                         }
//                 }
//         }
//
//         return options, offset
// }
