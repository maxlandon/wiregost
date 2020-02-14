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
	moduleHelp = fmt.Sprintf(`%s%s Module Commands%s 

%s About:%s Module usage 

%s General Commands:%s
    use             %sUse a module (either already on stack or not) (auto-completed modules)%s
    set             %sSet a value to an option of the current module. (auto-completed options)%s
    run             %sRun the main function of a module (compiling a payload, executing a post module, etc...)%s

%s Payload Commands:%s
    to_listener             %sBased on listener options provided, start a listener%s
    parse_profile <name>    %sUse an implant profile and parse it into the current module (auto-completed)%s
                            %s(Any payload type can parse any profile, it will just pick the options relevant to him)%s
    to_profile <name>       %sGenerate an implant profile based on the current module options%s

%s Examples:%s
    use payload/multi/single/reverse_dns    %sUse a a given payload (will load it on the stack)%s
    set DomainsDNS www.example.com          %sSet the DomainsDNS option with www.example.com (auto-completed options)%s
    parse_profile myExampleProfile          %sParse the profile named 'myExampleProfile' into the current payload module%s
    to_profile myNewProfile                 %sGenerate a profile named 'myNewProfile' from a payload module 'Implant' options%s`,
		tui.BLUE, tui.BOLD, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
	)
)
