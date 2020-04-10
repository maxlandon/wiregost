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

package module

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/util"
)

// ModuleInfoCmd - Show module information and options
type ModuleInfoCmd struct{}

var ModuleInfo ModuleInfoCmd

func RegisterModuleInfo() {
	MainParser.AddCommand(constants.ModuleInfo, "", "", &ModuleInfo)

	mi := MainParser.Find(constants.ModuleInfo)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], mi)
	mi.ShortDescription = "Show module information and options"
}

// Execute - Show module information and options
func (i *ModuleInfoCmd) Execute(args []string) error {
	m := Context.Module

	// Info
	fmt.Printf("%s%sModule:%s\r\t\t%s\r\n", tui.DIM, tui.BLUE, tui.RESET, m.Name)
	fmt.Printf("%s%sPlatform:%s \t%s (%s)\r\n", tui.DIM, tui.BLUE, tui.RESET, m.Platform, m.Targets)
	fmt.Printf("%s%sModule Authors:%s ", tui.DIM, tui.BLUE, tui.RESET)
	for a := range m.Author {
		fmt.Printf("%s ", m.Author[a])
	}
	fmt.Println()
	fmt.Printf("%s%sCredits:%s \t", tui.DIM, tui.BLUE, tui.RESET)
	for c := range m.Credits {
		fmt.Printf("%s ", m.Credits[c])
	}
	fmt.Println()
	fmt.Printf("%s%sLanguage:%s\r\t\t%s\n", tui.DIM, tui.BLUE, tui.RESET, m.Lang)
	fmt.Printf("%s%sPriviledged:%s \t%t\n", tui.DIM, tui.BLUE, tui.RESET, m.Priviledged)
	fmt.Println()
	fmt.Printf("%s%sDescription:%s\r\n", tui.DIM, tui.BLUE, tui.RESET)
	fmt.Println(tui.Dim(util.Wrap(m.Description, 100)))
	fmt.Println()

	// Options
	PrintOptions(m)

	// Notes
	if m.Notes != "" {
		fmt.Println()
		fmt.Printf("%sNotes:%s ", tui.YELLOW, tui.RESET)
		fmt.Println(tui.Dim(util.Wrap(m.Notes, 100)))
	}

	return nil
}
