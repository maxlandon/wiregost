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

package commands

import (
	"context"
	"fmt"

	"github.com/desertbit/grumble"
	"github.com/evilsocket/islazy/tui"

	consts "github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/help"
)

func RegisterServiceCommands(cctx *context.Context, app *grumble.App) {

	servicesCommand := &grumble.Command{
		Name:     "services",
		Help:     tui.Dim("Manage database services"),
		LongHelp: help.GetHelpFor(consts.ServicesStr),

		Flags: func(f *grumble.Flags) {
			f.StringL("service_id", "", "ID of a service. Available when listing them")
			f.StringL("addresses", "", "One or several IPv4/IPv6 Addresses, comma-separated (192.168.1.15,230.16.13.15)")
			f.StringL("port", "", "Port on which the service listens")
			f.StringL("proto", "", "Transport protocol used (tcp/udp,...)")
			f.StringL("state", "", "State of the service's port (open/closed/filtered/unknown)")
			f.StringL("name", "", "Service type name (http/https/dns/rpc...)")
			f.StringL("info", "", "Extra information about the service (usually its banner)")
		},
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			// hosts(cctx, gctx)
			fmt.Println()
			return nil
		},
	}

	// Add services
	servicesCommand.AddCommand(&grumble.Command{
		Name:  "add",
		Help:  tui.Dim("Manually add services"),
		Usage: help.GetHelpFor(consts.ServicesAdd),

		Flags: func(f *grumble.Flags) {
			f.StringL("service_id", "", "ID of a service. Available when listing them")
			f.StringL("addresses", "", "One or several IPv4/IPv6 Addresses, comma-separated (192.168.1.15,230.16.13.15)")
			f.StringL("port", "", "Port on which the service listens")
			f.StringL("proto", "", "Transport protocol used (tcp/udp,...)")
			f.StringL("state", "", "State of the service's port (open/closed/filtered/unknown)")
			f.StringL("name", "", "Service type name (http/https/dns/rpc...)")
			f.StringL("info", "", "Extra information about the service (usually its banner)")
		},
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			// hosts(cctx, gctx)
			fmt.Println()
			return nil
		},
	})

	// Delete services
	servicesCommand.AddCommand(&grumble.Command{
		Name:  "delete",
		Help:  tui.Dim("Manually delete services"),
		Usage: help.GetHelpFor(consts.ServicesDelete),

		Flags: func(f *grumble.Flags) {
			f.StringL("service_id", "", "ID of a service. Available when listing them")
			f.StringL("addresses", "", "One or several IPv4/IPv6 Addresses, comma-separated (192.168.1.15,230.16.13.15)")
			f.StringL("port", "", "Port on which the service listens")
			f.StringL("proto", "", "Transport protocol used (tcp/udp,...)")
			f.StringL("state", "", "State of the service's port (open/closed/filtered/unknown)")
			f.StringL("name", "", "Service type name (http/https/dns/rpc...)")
			f.StringL("info", "", "Extra information about the service (usually its banner)")
		},
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			// hosts(cctx, gctx)
			fmt.Println()
			return nil
		},
	})

	// Update services
	servicesCommand.AddCommand(&grumble.Command{
		Name:  "update",
		Help:  tui.Dim("Manually update services"),
		Usage: help.GetHelpFor(consts.ServicesUpdate),

		Flags: func(f *grumble.Flags) {
			f.StringL("service_id", "", "ID of a service. Available when listing them")
			f.StringL("addresses", "", "One or several IPv4/IPv6 Addresses, comma-separated (192.168.1.15,230.16.13.15)")
			f.StringL("port", "", "Port on which the service listens")
			f.StringL("proto", "", "Transport protocol used (tcp/udp,...)")
			f.StringL("state", "", "State of the service's port (open/closed/filtered/unknown)")
			f.StringL("name", "", "Service type name (http/https/dns/rpc...)")
			f.StringL("info", "", "Extra information about the service (usually its banner)")
		},
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			// hosts(cctx, gctx)
			fmt.Println()
			return nil
		},
	})

	// Register root service command
	app.AddCommand(servicesCommand)
}
