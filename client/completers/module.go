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
	"os"
	"path/filepath"
	"strings"

	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/server/assets"
)

// AutoCompleter is the autocompletion engine
type ModuleCompleter struct {
	Command *commands.Command
}

// Do is the completion function triggered at each line
func (mc *ModuleCompleter) Do(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {

	splitLine := strings.Split(string(line), " ")
	line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))

	// types := []string{"exploit", "post", "auxiliary", "payload"}

	dirs := []string{}

	// Get all dirs and subdirs in modules
	err := filepath.Walk(assets.GetModulesDir(),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if !strings.HasSuffix(path, "/docs") {
					path = strings.TrimPrefix(path, assets.GetModulesDir())
					path = strings.TrimPrefix(path, "/")
					dirs = append(dirs, path)
				}
			}
			return nil
		})
	if err != nil {
		return
	}

	// Filter out subdirs not containing modules
	filtered := difference(dirs, moduleSubDirs)

	for _, dir := range filtered {
		search := dir
		if !hasPrefix(line, []rune(search)) {
			sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
			options = append(options, sLine...)
			offset = sOffset
		}
	}
	// for _, dir := range types {
	//         search := dir
	//         if !hasPrefix(line, []rune(search)) {
	//                 sLine, sOffset := doInternal(line, pos, len(line), []rune(search+"/"))
	//                 options = append(options, sLine...)
	//                 offset = sOffset
	//         } else {
	//                 dirs := []string{}
	//                 err := filepath.Walk(assets.GetModulesDir(),
	//                         func(path string, info os.FileInfo, err error) error {
	//                                 if err != nil {
	//                                         return err
	//                                 }
	//                                 if info.IsDir() {
	//                                         if !strings.HasSuffix(path, "/docs") {
	//                                                 path = strings.TrimPrefix(path, assets.GetModulesDir())
	//                                                 if strings.HasPrefix(path, "/"+search) {
	//                                                         path = strings.TrimPrefix(path, "/"+search)
	//                                                         dirs = append(dirs, path)
	//                                                 }
	//                                         }
	//                                 }
	//                                 return nil
	//                         })
	//                 if err != nil {
	//                         return
	//                 }

	// splitLine := strings.Split(string(line), " ")
	// line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))

	// words := strings.Split(string(line), "/")
	// argInput := lastString(words)

	// For some arguments, the split results in a last empty item.
	// if words[len(words)-1] == "" {
	//         argInput = words[0]
	// }

	// dirs, _ := ioutil.ReadDir(assets.GetModulesDir() + "/" + string(line))
	// for _, dir := range dirs {
	//         search := dir
	//         if !hasPrefix(line, []rune(search)) {
	//                 sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
	//                 options = append(options, sLine...)
	//                 offset = sOffset
	//                 // options = append(options, []rune(search+"/"))
	//                 // offset = len(search)
	//         }
	// }
	// }
	// }
	return options, offset
}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

var moduleSubDirs = []string{
	// Post
	"post",
	"post/windows",
	"post/windows/x64",
	"post/windows/x64/go",
	"post/windows/x64/go/credentials",

	// Payload
	"payload",
	"payload/multi",
	"payload/multi/single",
	"payload/multi/stager",

	// Exploit
	"exploit",

	// Auxiliary
	"auxiliary",
}
