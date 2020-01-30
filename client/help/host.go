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
	hostsHelp = fmt.Sprintf(`%s%s Command:%s hosts <verb> <options> <filters>%s

%s About:%s Manage database hosts (shows all workspace hosts if no verb and no options)
        (Type 'hosts add', 'hosts delete' or 'hosts update' for further command-specific examples)

%s Options:%s
    search          %sSearch hosts with filters%s
    add             %sAdd the hosts instead of searching%s
    delete          %sDelete the hosts instead of searching%s
    update          %sUpdate the hosts instead of searching (need host_id)%s

%s Filters:%s
    host_id     %sID of a host. Available when listing them%s
    hostnames   %sOne or several hostnames%s
    os_name     %sOS name of a host (Windows 10/Linux)%s
    os_family   %sOS family of a host%s
    os_flavor   %sOS flavor of a host%s
    os_sp       %sOS Service Pack (windows) or kernel version (Unix/Apple) of a host%s
    arch        %sCPU architecture of a host%s
    addresses   %sOne or several IPv4/IPv6 Addresses, comma-separated (192.168.1.15,230.16.13.15)%s
    info        %sInfo of a host %s
    name        %sName of a host %s
    comment     %sComment of a host %s
    up          %sOnly show hosts which are up %s

%s Examples:%s
    hosts add addresses=192.168.1.24 os-family=Windows          %sManually add a host with address and os_family%s
    hosts search addresses=220.188.2.15 arch=64                 %sList hosts that match address and arch%s
    hosts search addresses=220.188.2.15,192.168.0.12            %sList hosts that match one of these addresses%s
    hosts update host-id=23 addresses=192.34.23.1 arch=x86      %sUpdate a host by changing its addresses and its CPU arch%s`,
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

	hostsAdd = fmt.Sprintf(`%s%s Usage %s%s

    hosts add addresses=192.168.1.24 os-family=Windows                      %sadd a host with address and os_family%s
    hosts add addresses=220.188.2.15 arch=64                                %sAdd a host with address and arch%s
    hosts add addresses=220.188.2.15,192.168.0.12 comment="A Comment"       %sAdd a host with these 2 addresses and a comment
                                                                            (If any of these addresses is already used, it will not create the host)%s`,
		tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET)

	hostsDelete = fmt.Sprintf(`%s%s Usage %s%s

    hosts delete addresses=192.168.1.24 os-family=Windows       %sDelete hosts matching both address and os_family%s
    hosts delete addresses=220.188.2.15 arch=64                 %sDelete hosts matching both address and arch%s
    hosts delete host-id=2 "                                    %sDelete host with ID 2%s`,
		tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET)

	hostsUpdate = fmt.Sprintf(`%s%s Usage %s%s

    // %sMandatory arguments%s: host-id

    hosts update host-id=2 os-family=Windows                    %sUpdate host with ID 2, and change its OS family%s
    hosts update host-id=3 addresses=220.188.2.15 arch=64       %sUpdate a host's addresses (will overwrite existing ones) and CPU arch%s
    hosts update host-id=2 comment="Updated manually"           %sUpdate host with ID 2 and change its comment%s`,
		tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
		tui.BOLD, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET)
)
