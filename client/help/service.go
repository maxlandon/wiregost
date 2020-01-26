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
	serviceHelp = fmt.Sprintf(`%s%s Command:%s services <verb> <options> <filters>%s

%s About:%s Manage database services (shows all workspace services if no verb and no options)
        (Type 'services add', 'services delete' or 'services update' for further command-specific examples)

%s Options:%s
    add             %sAdd services instead of searching%s
    delete          %sDelete services instead of searching%s
    update          %sUpdate a service instead of searching (need host_id)%s
    -S, --search    %sSearch string to filter by%s

%s Filters:%s
    service_id  %sID of a service. Available when listing them%s
    addresses   %sOne or several IPv4/IPv6 Addresses, comma-separated (192.168.1.15,230.16.13.15)%s
    port        %sPort on which the service listens%s
    proto       %sTransport protocol used (tcp/udp,...)%s
    state       %sState of the service's port (open/closed/filtered/unknown)%s
    name        %sService type name (http/https/dns/rpc...)%s
    info        %sExtra information about the service (usually its banner)%s

%s Examples:%s
    services add --addresses 192.168.1.24 --state open          %sManually add a service with host address and port state%s
    services --addresses 220.188.2.15 --name https              %sList all HTTPS services on host 220.188.2.15%s
    services --addresses 220.188.2.15,192.168.0.12              %sList services that match one of these addresses%s
    services update --service_id 23 --addresses 192.34.23.1     %sUpdate a service by changing its host(address) %s`,
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
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET)

	servicesAdd = fmt.Sprintf(`%s%s Usage %s%s

    // %sMandatory arguments%s: port, addresses 

    services add --addresses 192.168.1.24 --port 443 --name https       %sAdd a service on a host address, with port and name%s
    services add --addresses 220.188.2.15,220.188.1.1 --name dns        %sAdd a DNS service on two different hosts%s
    hosts add --addresses 220.188.2.15 --port 8080 --state filtered     %sAdd a service on address:port and with proto. 
                                                                        (If the address:port is already used, it will not create the service)%s`,
		tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
		tui.BOLD, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET)

	servicesDelete = fmt.Sprintf(`%s%s Usage %s%s

    services delete --addresses 192.168.1.24 --port 443 --name https       %sDelete service on this address, with port and name%s
    services delete --addresses 220.188.2.15,220.188.1.1 --name dns        %sDelete all DNS services on these two addresses%s
    hosts delete --service_id 3                                            %sDelete service with ID 3%s`,
		tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET)

	servicesUpdate = fmt.Sprintf(`%s%s Usage %s%s

    // %sMandatory arguments%s: service_id

    services update --service_id 3 --addresses 192.168.1.24 --port 443      %sChange a service  address and port%s
    hosts update --service_id 4 --addresses 220.188.2.15 --proto tcp        %sChange a service address and protocol%s`,
		tui.BLUE, tui.BOLD, tui.FOREWHITE, tui.RESET,
		tui.BOLD, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET)
)
