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

package util

import (
	"fmt"
	"os"

	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"
)

func Table() *tablewriter.Table {

	table := tablewriter.NewWriter(os.Stdout)

	// Appearance
	table.SetCenterSeparator(fmt.Sprintf("%s|%s", tui.FOREBLACK, tui.RESET))
	table.SetColumnSeparator(fmt.Sprintf("%s|%s", tui.FOREBLACK, tui.RESET))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetAutoWrapText(true)
	table.SetColWidth(20)
	table.SetBorder(false)

	return table
}
