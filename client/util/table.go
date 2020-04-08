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
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/lmorg/readline"
	"github.com/olekukonko/tablewriter"
)

// Table is a wrapper around tablewriter.Table, so that we can customize behavior
type Table struct {
	*tablewriter.Table
}

// NewTable - Constructor method, with much of the default behavior we want
func NewTable() *Table {
	tab := &Table{tablewriter.NewWriter(os.Stdout)}

	// Borders
	tab.SetBorder(false)
	tab.SetCenterSeparator(fmt.Sprintf("%s|%s", tui.FOREBLACK, tui.RESET))
	tab.SetColumnSeparator(fmt.Sprintf("%s|%s", tui.FOREBLACK, tui.RESET))
	tab.SetRowSeparator(tui.Dim("-"))

	// Headers
	tab.SetAutoFormatHeaders(false)
	tab.SetColWidth(30) // Default column width

	// Cells
	tab.SetAlignment(tablewriter.ALIGN_LEFT)
	tab.SetReflowDuringAutoWrap(false) // Multiple lines in a cell if needed
	tab.SetAutoWrapText(false)

	return tab
}

// Output - Renders the table
func (t *Table) Output() {
	// Render
	t.Render()
}

// SetColumns - Simpler way to specify columns titles
func (t *Table) SetColumns(names []string, widths []int) {

	titles := []string{}             // header titles
	colors := []tablewriter.Colors{} // header colors
	def := tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor}

	for _, name := range names {
		titles = append(titles, tui.Dim(name))
		colors = append(colors, def)
	}
	t.SetHeader(titles)

	// Minimum width
	for i, width := range widths {
		if width == 0 {
			t.SetColMinWidth(i, 30)
			continue
		}
		t.SetColMinWidth(i, width)
	}
}

// WrapColumns - Value for maximum shell width, used for printing tables
var WrapColumns int

func AutoWrap(text string) (wrapped string) {

	// We check the current terminal width and adapt to it for wrapping
	termWidth := readline.GetTermWidth()
	var lineWidth int
	switch width := termWidth; {
	case width > 0 && width <= 50:
		lineWidth = 10
	case width > 50 && width <= 80:
		lineWidth = 20
	case width > 80 && width <= 110:
		lineWidth = 40
	case width > 110 && width <= 130:
		lineWidth = 60
	case width > 130 && width <= 150:
		lineWidth = 80
	case width > 150 && width <= 170:
		lineWidth = 100
	case width > 170 && width <= 190:
		lineWidth = 120
	case width > 190 && width <= 210:
		lineWidth = 140
	case width > 210 && width <= 230:
		lineWidth = 160
	case width > 230 && width <= 250:
		lineWidth = 180
	case width > 250 && width <= 270:
		lineWidth = 200
	case width > 270 && width <= 290:
		lineWidth = 220
	case width > 290:
		lineWidth = 240
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return
	}
	wrapped = words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}
	return
}

// Wrap - Function used for description paragraphs and table columns
func Wrap(text string, lineWidth int) (wrapped string) {
	words := strings.Fields(text)
	if len(words) == 0 {
		return
	}
	wrapped = words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}
	return
}
