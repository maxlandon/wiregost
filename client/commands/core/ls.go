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

package core

import (
	"fmt"
	"os"

	"github.com/evilsocket/islazy/fs"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/util"
)

// Command -----------------------------------------------------------------------------------

// ListDirectories - List directory contents
type ListDirectoriesCmd struct {
	Positional struct {
		Path      string   `description:"Local directory/file"`
		OtherPath []string `description:"Local directory/file" `
	} `positional-args:"yes"`
}

var ListDirectories ListDirectoriesCmd

func RegisterLs() {
	MainParser.AddCommand("ls", "", "", &ListDirectories)

	ls := MainParser.Find("ls")
	CommandMap[MAIN_CONTEXT] = append(CommandMap[MAIN_CONTEXT], ls)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], ls)
	ls.ShortDescription = "List directory contents"
}

// Execute - Command
func (ls *ListDirectoriesCmd) Execute(args []string) error {

	base := []string{"ls", "--color", "-l"}

	var fullPath string
	if ls.Positional.Path == "" {
		wd, _ := os.Getwd()
		fullPath, _ = fs.Expand(wd)
	} else {
		fullPath, _ = fs.Expand(ls.Positional.Path)
	}

	base = append(base, fullPath)

	fullOtherPaths := []string{}
	for _, path := range ls.Positional.OtherPath {
		full, _ := fs.Expand(path)
		fullOtherPaths = append(fullOtherPaths, full)
	}
	base = append(base, fullOtherPaths...)

	err := util.Shell(base)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
