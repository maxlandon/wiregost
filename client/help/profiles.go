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
	profileHelp = fmt.Sprintf(`%s%s Profile Commands%s 

%s About:%s Manage implant profiles 

%s Commands:%s
    profiles    %sList saved implant profiles and their options%s

%s Notes:%s
    - Profiles can be used with payload modules. 
    - The 'Implant' options of these modules, when using 'parse_profile' will be populated with the selected 
      profile options. Any profile can be used with any payload type: the module will only use the options 
      from the profile that are relevant to him.
    - Conversely, profiles can be generated using the payload module command 'to_profile'`,
		tui.BLUE, tui.BOLD, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.YELLOW, tui.RESET,
	)
)
