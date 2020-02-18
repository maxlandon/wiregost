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

	"github.com/evilsocket/islazy/tui"
)

var (
	Info    = fmt.Sprintf("%s[-]%s ", tui.BLUE, tui.RESET)
	Warn    = fmt.Sprintf("%s[!]%s ", tui.YELLOW, tui.RESET)
	Error   = fmt.Sprintf("%s[!]%s ", tui.RED, tui.RESET)
	Success = fmt.Sprintf("%s[*]%s ", tui.GREEN, tui.RESET)

	Infof   = fmt.Sprintf("%s[-] ", tui.BLUE)
	Warnf   = fmt.Sprintf("%s[!] ", tui.YELLOW)
	Errorf  = fmt.Sprintf("%s[!] ", tui.RED)
	Sucessf = fmt.Sprintf("%s[*] ", tui.GREEN)

	RPCError     = fmt.Sprintf("%s[RPC Error]%s ", tui.RED, tui.RESET)
	CommandError = fmt.Sprintf("%s[Command Error]%s ", tui.RED, tui.RESET)
	DBError      = fmt.Sprintf("%s[DB Error]%s ", tui.RED, tui.RESET)
)
