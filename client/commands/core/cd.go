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
	"github.com/evilsocket/islazy/tui"

	. "github.com/maxlandon/wiregost/client/commands"
)

// Command -----------------------------------------------------------------------------------

// ChangeDirectory - Change the working directory of the client console
type ChangeDirectoryCmd struct {
	Positional struct {
		Path string `description:"Local path" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var ChangeDirectory ChangeDirectoryCmd

func RegisterCd() {
	// cd
	MainParser.AddCommand("cd", "", "", &ChangeDirectory)

	cd := MainParser.Find("cd")
	CommandMap[MAIN_CONTEXT] = append(CommandMap[MAIN_CONTEXT], cd)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], cd)
	cd.ShortDescription = "Change the client console working directory"
	cd.LongDescription = CdHelp

	cd.Args()[0].RequiredMaximum = 1
}

// Execute - Handler for ChangeDirectory
func (cd *ChangeDirectoryCmd) Execute(args []string) error {

	dir, err := fs.Expand(cd.Positional.Path)
	err = os.Chdir(dir)
	if err != nil {
		fmt.Printf(CommandError+"%s \n", err)
	} else {
		fmt.Printf(Info+"Changed directory to %s \n", dir)
	}

	return nil
}

// Help -----------------------------------------------------------------------------------

var CdHelp = fmt.Sprintf(`%s%s Command:%s cd <dir>%s

%s About:%s Change the client console working directory.`,
	tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
	tui.YELLOW, tui.RESET,
)
