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
	workspaceHelp = fmt.Sprintf(`%s%s Command:%s workspace <verb> <filters>%s

%s About:%s Manage workspaces (shows workspaces if no verb and no options)

%s Verbs:%s
    add             %sAdd a workspace %s
    switch          %sSwitch to a workspace %s
    delete          %sDelete a workspace (need workspace_id)%s
    update          %sUpdate a workspace with fields provided as filters (need workspace_id)%s

%s Filters:%s
    --name              %sName of workspace to add/update %s 
    --boundary          %sOne or several IPv4/IPv6 Addresses/Ranges, comma-separated (192.168.1.15,230.16.13.15)%s
    --description       %sA description for the workspace%s 
    --limit-to-network  %s(true/false) Will limit the interaction of tools like scanners and exploits to the workspace network boundary%s

%s Examples:%s
    workspace switch Wiregost                                               %sSwitch to workspace WireGost%s
    workspace add --name WireGost --boundary 192.168.1.0/24                 %sAdd a workspace named Wiregost, with a boundary but not restricted to it%s
    workspace update --name default --description "A test description"      %sUpdate a workspace by appending a new address and changing arch%s`,
		tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
		tui.YELLOW, tui.RESET,
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
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET)
)
