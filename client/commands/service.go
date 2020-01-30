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
	"strconv"
	"strings"

	"github.com/desertbit/grumble"
	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"

	consts "github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/help"
	"github.com/maxlandon/wiregost/client/util"
	"github.com/maxlandon/wiregost/data_service/models"
	"github.com/maxlandon/wiregost/data_service/remote"
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
			services(cctx, gctx)
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
			newApp(app, cctx, gctx)
			// addService(cctx, gctx)
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

func services(cctx *context.Context, gctx *grumble.Context) {

	// Get Services
	var services []models.Port
	var err error
	opts := serviceFilters(gctx)
	if len(opts) == 0 {
		services, err = remote.Services(*cctx, nil)
		if err != nil {
			fmt.Printf("%s[!]%s Error: %s\n",
				tui.RED, tui.RESET, err.Error())
			return
		}
	} else {
		services, err = remote.Services(*cctx, opts)
		if err != nil {
			fmt.Printf("%s[!]%s Error: %s\n",
				tui.RED, tui.RESET, err.Error())
			return
		}
	}

	// Print table
	servicesTable(cctx, &services)
}

func addService(cctx *context.Context, gctx *grumble.Context) {
}

func deleteServices(cctx *context.Context, gctx *grumble.Context) {
}

func updateService(cctx *context.Context, gctx *grumble.Context) {
}

func serviceFilters(ctx *grumble.Context) (opts map[string]interface{}) {
	opts = make(map[string]interface{}, 0)

	return opts
}

func servicesTable(cctx *context.Context, services *[]models.Port) {
	// Get host addresses for services
	hopts := make(map[string]interface{}, 0)
	keys := make(map[uint]bool)
	list := []uint{}
	for _, s := range *services {
		if _, value := keys[s.HostID]; !value {
			keys[s.HostID] = true
			list = append(list, s.HostID)
		}
	}
	hopts["host_id"] = list

	hosts, err := remote.Hosts(*cctx, hopts)
	if err != nil {
		fmt.Printf("%s[!]%s Error: %s\n",
			tui.RED, tui.RESET, err.Error())
	}
	ips := make(map[uint][]string)
	for _, h := range hosts {
		addrs := []string{}
		for _, a := range h.Addresses {
			addrs = append(addrs, a.Addr)
		}
		ips[h.ID] = addrs
	}

	// Table
	table := util.Table()
	table.SetHeader([]string{"ID", "Address", "Port", "Proto", "Name", "State", "Info"})
	table.SetColWidth(80)
	table.SetColMinWidth(6, 80)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)
	table.SetAutoFormatHeaders(true)

	data := [][]string{}
	for _, s := range *services {
		servID := uint64(s.ID)
		port := uint64(s.Number)
		row := []string{strconv.FormatUint(servID, 10), strings.Join(ips[s.HostID], " "),
			strconv.FormatUint(port, 10), s.Protocol, s.Service.Name, s.State.State, s.Service.ExtraInfo}
		data = append(data, row)
	}

	table.AppendBulk(data)
	table.Render()

}

func newApp(app *grumble.App, cctx *context.Context, gctx *grumble.Context) {
	app = grumble.New(&grumble.Config{
		Name:        "Wiregost",
		Description: tui.Blue(tui.Bold("Wiregost Client")),
		// HistoryFile:     home + "/.wiregost/client/.history",
		HistoryLimit:    5000,
		Prompt:          "newContext >",
		HelpSubCommands: true,
	})
	// RegisterCommands(nil, cctx, app)
	app.Run()
}
