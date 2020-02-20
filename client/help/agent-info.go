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

package help

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
)

var (
	infoHelp = fmt.Sprintf(`%s%s Implant Info Commands%s 

%s About:%s Show various informations about implant/target 

%s Commands:%s
    ifconfig        %sPrint network/interfaces information of target%s
    info            %sPrint implant/target information%s
    whoami          %sShow username/hostname of target%s
    ping            %sPing the implant. (This does NOT send an ICMP packet, but an empty C2 message)%s
    getuid          %sGet implant process User ID%s
    getpid          %sGet implant Process ID%s
    getgid          %sGet implant process Group ID%s`,
		tui.BLUE, tui.BOLD, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)
)
