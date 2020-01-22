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
	consts "github.com/maxlandon/wiregost/client/constants"
)

var (
	cmdHelp = map[string]string{
		// Database
		consts.WorkspaceStr: workspaceHelp,
		consts.HostsStr:     hostsHelp,
	}

	// Database ------------------------------------------------------------------------------//

	hostsHelp = fmt.Sprintf(`%s%s Command:%s hosts <verb> <options> <filters>%s

%s About:%s Manage database hosts (shows all workspace hosts if no verb and no options)

%s Options:%s
    add             %sAdd the hosts instead of searching%s
    delete          %sDelete the hosts instead of searching%s
    update          %sUpdate the hosts instead of searching (need host_id)%s
    -u, --up        %sOnly show hosts which are up %s
    -S, --search    %sSearch string to filter by%s
    -i, --info      %sChange the info of a host %s
    -n, --name      %sChange the name of a host %s
    -m, --comment   %sChange the comment of a host %s

%s Filters:%s
    host_id     %sID of a host. Available when listing them, cannot use default%s
    hostnames   %sOne or several hostnames%s
    os_name     %sOS name of a host (Windows 10/Linux)%s
    os_family   %sOS family of a host%s
    os_flavor   %sOS flavor of a host%s
    os_sp       %sOS Service Pack (windows) or kernel version (Unix/Apple) of a host%s
    arch        %sCPU architecture of a host%s
    addresses   %sOne or several IPv4/IPv6 Addresses, comma-separated (192.168.1.15,230.16.13.15)%s

%s Examples:%s
    hosts add --addresses 192.168.1.24 --os_family Windows          %sManually add a host with address and os_family%s
    hosts --addresses 220.188.2.15 --arch 64                        %sList hosts that match address and arch%s
    hosts --addresses 220.188.2.15,192.168.0.12                     %sList hosts that match one of these addresses%s
    hosts update --host_id 23 --addresses 192.34.23.1 --arch x86    %sUpdate a host by appending a new address and changing arch%s`,
		tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET)

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

// GetHelpFor - Get help string for a command
func GetHelpFor(cmdName string) string {
	if 0 < len(cmdName) {
		if helpTempl, ok := cmdHelp[cmdName]; ok {
			return helpTempl
		}
	}
	return ""
}
