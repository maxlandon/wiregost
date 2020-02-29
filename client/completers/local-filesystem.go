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
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/evilsocket/islazy/fs"
	"github.com/maxlandon/wiregost/client/commands"
)

// AutoCompleter is the autocompletion engine
type PathCompleter struct {
	Command *commands.Command
}

// Do is the completion function triggered at each line
func (pc *PathCompleter) Do(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {

	splitLine := strings.Split(string(line), " ")
	line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))

	// 1) Get the absolute path. There are two cases:
	//      - The path is "rounded" with a slash: no filter to keep.
	//      - The path is not a slash: a filter to keep for later.
	// We keep a boolean for remembering which case we found
	linePath := ""
	path := ""
	lastPath := ""
	if strings.HasSuffix(string(line), "/") {
		// Trim the non needed slash
		linePath = strings.TrimSuffix(string(line), "/")
		linePath = filepath.Dir(string(line))
		// Get absolute path
		path, _ = fs.Expand(string(linePath))

	} else if string(line) == "" {
		linePath = "."
		path, _ = fs.Expand(string(linePath))
	} else {
		linePath = string(line)
		linePath = filepath.Dir(string(line))
		// Get absolute path
		path, _ = fs.Expand(string(linePath))
		// Save filter
		lastPath = filepath.Base(string(line))
	}

	// 2) We take the absolute path we found, and get all dirs in it.
	var dirs []string
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}

	switch lastPath {
	case "":
		for _, dir := range dirs {
			search := dir + "/"
			if !hasPrefix([]rune(lastPath), []rune(search)) {
				sLine, sOffset := doInternal([]rune(lastPath), pos, len([]rune(lastPath)), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}
	default:
		filtered := []string{}
		for _, dir := range dirs {
			if strings.HasPrefix(dir, lastPath) {
				filtered = append(filtered, dir)
			}
		}

		for _, dir := range filtered {
			search := dir + "/"
			if !hasPrefix([]rune(lastPath), []rune(search)) {
				sLine, sOffset := doInternal([]rune(lastPath), pos, len([]rune(lastPath)), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}

	}

	return options, offset
}
