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

var ()

func RegisterHostCommands(cctx *context.Context, app *grumble.App) {

	hostsCommand := &grumble.Command{
		Name:     consts.HostsStr,
		Help:     tui.Dim("Manage database hosts"),
		LongHelp: help.GetHelpFor(consts.HostsStr),

		Flags: func(f *grumble.Flags) {
			// Filters
			f.UintL("host_id", 0, "ID of a host. Available when listing them, cannot use default")
			f.StringL("addresses", "", "One or several IPv4/IPv6 Addresses, comma-separated (192.168.1.15,230.16.13.15)")
			f.StringL("hostnames", "", "One or several hostnames")
			f.StringL("os_name", "", "OS name of a host (Windows 10/Linux)")
			f.StringL("os_family", "", "OS family of a host")
			f.StringL("os_flavor", "", "OS flavor of a host")
			f.StringL("os_sp", "", "OS Service Pack (windows) or kernel (Unix/Apple) of a host")
			f.StringL("arch", "", "CPU architecture of a host")
			f.StringL("name", "", "Name given to a host")
			f.StringL("info", "", "Informations on a host")
			f.StringL("comment", "", "Comment about a host")
			f.BoolL("up", false, "Hosts that are up")
		},
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			hosts(cctx, gctx)
			fmt.Println()
			return nil
		},
		HelpGroup: consts.DataServiceHelpGroup,
	}

	// Add hosts
	hostsCommand.AddCommand(&grumble.Command{
		Name: "add",
		Help: tui.Dim("Add a host manually"),
		Flags: func(f *grumble.Flags) {
			// Filters
			f.StringL("addresses", "", "One or several IPv4/IPv6 Addresses, comma-separated (192.168.1.15,230.16.13.15)")
			f.StringL("hostname", "", "One or several hostnames")
			f.StringL("os_name", "", "OS name of a host (Windows 10/Linux)")
			f.StringL("os_family", "", "OS family of a host")
			f.StringL("os_flavor", "", "OS flavor of a host")
			f.StringL("os_sp", "", "OS Service Pack (windows) or kernel (Unix/Apple) of a host")
			f.StringL("arch", "", "CPU architecture of a host")
			f.StringL("name", "", "Name given to a host")
			f.StringL("info", "", "Informations on a host")
			f.StringL("comment", "", "Comment about a host")
			f.BoolL("up", false, "Hosts that are up")
		},
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			addHost(cctx, gctx)
			fmt.Println()
			return nil
		},
		HelpGroup: consts.DataServiceHelpGroup,
	})

	// Delete hosts
	hostsCommand.AddCommand(&grumble.Command{
		Name: "delete",
		Help: tui.Dim("Delete one or more hosts manually"),
		Flags: func(f *grumble.Flags) {
			// Filters
			f.UintL("host_id", 0, "ID of a host. Available when listing them, cannot use default")
			f.StringL("addresses", "", "One or several IPv4/IPv6 Addresses, comma-separated (192.168.1.15,230.16.13.15)")
			f.StringL("hostname", "", "One or several hostnames")
			f.StringL("os_name", "", "OS name of a host (Windows 10/Linux)")
			f.StringL("os_family", "", "OS family of a host")
			f.StringL("os_flavor", "", "OS flavor of a host")
			f.StringL("os_sp", "", "OS Service Pack (windows) or kernel (Unix/Apple) of a host")
			f.StringL("arch", "", "CPU architecture of a host")
			f.StringL("name", "", "Name given to a host")
			f.StringL("info", "", "Informations on a host")
			f.StringL("comment", "", "Comment about a host")
			f.BoolL("up", false, "Hosts that are up")
		},
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			deleteHosts(cctx, gctx)
			fmt.Println()
			return nil
		},
		HelpGroup: consts.DataServiceHelpGroup,
	})

	// Update host
	hostsCommand.AddCommand(&grumble.Command{
		Name: "update",
		Help: tui.Dim("Update a host manually"),
		Flags: func(f *grumble.Flags) {
			// Filters
			f.UintL("host_id", 0, "ID of a host. Available when listing them, cannot use default")
			f.StringL("addresses", "", "One or several IPv4/IPv6 Addresses, comma-separated (192.168.1.15,230.16.13.15)")
			f.StringL("hostname", "", "One or several hostnames")
			f.StringL("os_name", "", "OS name of a host (Windows 10/Linux)")
			f.StringL("os_family", "", "OS family of a host")
			f.StringL("os_flavor", "", "OS flavor of a host")
			f.StringL("os_sp", "", "OS Service Pack (windows) or kernel (Unix/Apple) of a host")
			f.StringL("arch", "", "CPU architecture of a host")
			f.StringL("name", "", "Name given to a host")
			f.StringL("info", "", "Informations on a host")
			f.StringL("comment", "", "Comment about a host")
			f.BoolL("up", false, "Hosts that are up")
		},
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			updateHost(cctx, gctx)
			fmt.Println()
			return nil
		},
		HelpGroup: consts.DataServiceHelpGroup,
	})

	// Finally, register root command
	app.AddCommand(hostsCommand)
}

func hosts(cctx *context.Context, gctx *grumble.Context) {

	var hosts []models.Host
	var err error
	opts := parseFilters(gctx)
	if len(opts) == 0 {
		hosts, err = remote.Hosts(*cctx, nil)
		if err != nil {
			fmt.Printf("%s[!]%s Error: %s\n",
				tui.RED, tui.RESET, err.Error())
			return
		}
	} else {
		hosts, err = remote.Hosts(*cctx, opts)
		if err != nil {
			fmt.Printf("%s[!]%s Error: %s\n",
				tui.RED, tui.RESET, err.Error())
			return
		}
	}

	// Table
	table := util.Table()
	table.SetHeader([]string{"ID", "Addresses", "Name", "OS Name", "OS Flavor", "OS SP", "Arch", "Purpose", "Info", "Comments"})
	table.SetColWidth(40)
	table.SetColMinWidth(1, 15)
	table.SetColMinWidth(8, 20)
	table.SetColMinWidth(9, 40)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	data := [][]string{}
	for _, h := range hosts {
		addrsList := []string{}

		for _, a := range h.Addresses {
			addrsList = append(addrsList, a.String())
		}
		addrs := ""
		if len(addrsList) == 1 {
			addrs = addrsList[0]
		} else {
			addrs = strings.Join(addrsList, " ")
		}

		hostID := uint64(h.ID)
		data = append(data, []string{strconv.FormatUint(hostID, 10), addrs, h.Name,
			h.OSName, h.OSFlavor, h.OSSp, h.Arch, h.Purpose, h.Info, h.Comment})
	}
	table.AppendBulk(data)
	table.Render()

}

func addHost(cctx *context.Context, gctx *grumble.Context) {

	opts := parseFilters(gctx)

	// Get existing host for comparing ReportHost() results
	hosts, err := remote.Hosts(*cctx, nil)

	// If an IP address has been given as option, and that it matches one
	// of already existing hosts, it will report this existing host
	host, err := remote.ReportHost(*cctx, opts)
	if err != nil {
		fmt.Printf("%s[!]%s Error: %s\n",
			tui.RED, tui.RESET, err.Error())
		return
	}

	// Check if the hosts that have been returned before reporting a host
	// match with it. Means none has been created
	for _, h := range hosts {
		if h.ID == host.ID {
			fmt.Printf("%s*%s Host already exists at: %s\n",
				tui.YELLOW, tui.RESET, host.Addresses)
			return
		}
	}
	// Else notify a new host has been created
	fmt.Printf("%s*%s New host at: %s\n",
		tui.BLUE, tui.RESET, host.Addresses)
}

func deleteHosts(cctx *context.Context, gctx *grumble.Context) {

	var hosts []models.Host
	var err error

	opts := parseFilters(gctx)

	// Get a list of hosts matching filters given
	switch len(opts) {
	case 0:
		fmt.Printf("%s[!]%s Provide filters for host selection \n",
			tui.RED, tui.RESET)
	default:
		hosts, err = remote.Hosts(*cctx, opts)
		if err != nil {
			fmt.Printf("%s[!]%s Error: %s\n",
				tui.RED, tui.RESET, err.Error())
			return
		}
		if len(hosts) == 0 {
			fmt.Printf("%s*%s No hosts match the given filters \n",
				tui.YELLOW, tui.RESET)
			return
		}

		for i, _ := range hosts {
			opts["host_id"] = hosts[i].ID
			err = remote.DeleteHost(*cctx, opts)
			if err != nil {
				fmt.Printf("%s[!]%s Error: %s\n",
					tui.RED, tui.RESET, err.Error())
				continue
			} else {
				fmt.Printf("%s*%s Deleted host at: %s\n",
					tui.BLUE, tui.RESET, hosts[i].Addresses)
			}
		}
	}
}

func updateHost(cctx *context.Context, gctx *grumble.Context) {

	var host *models.Host

	// Parse host_id
	if gctx.Flags.Uint("host_id") != 0 {
		hosts, _ := remote.Hosts(*cctx, nil)
		for i, _ := range hosts {
			if hosts[i].ID == gctx.Flags.Uint("host_id") {
				host = &hosts[i]
			}
		}
	} else {
		// If not, return and ask for one
		fmt.Printf("%s[!]%s Provide a host ID (--host_id 2)\n",
			tui.RED, tui.RESET)
		return
	}

	// Parse options
	if gctx.Flags.String("os_name") != "" {
		host.OSName = gctx.Flags.String("os_name")
	}
	if gctx.Flags.String("os_family") != "" {
		host.OSFamily = gctx.Flags.String("os_family")
	}
	if gctx.Flags.String("os_flavor") != "" {
		host.OSFlavor = gctx.Flags.String("os_flavor")
	}
	if gctx.Flags.String("os_sp") != "" {
		host.OSSp = gctx.Flags.String("os_sp")
	}
	if gctx.Flags.String("arch") != "" {
		host.Arch = gctx.Flags.String("arch")
	}
	if gctx.Flags.String("addresses") != "" {
		addrsString := strings.Split(gctx.Flags.String("addresses"), ",")
		var addrs []models.Address
		exists := false
		for _, a := range addrsString {
			addr := models.Address{Addr: a, HostID: host.ID, AddrType: "IPv4"}
			for i, _ := range host.Addresses {
				if host.Addresses[i].Addr == addr.Addr {
					exists = true
				}
			}
			if exists == false {
				host.Addresses = append(host.Addresses, addr)
			} else {
				exists = false
			}
		}
		host.Addresses = addrs
	}
	// if gctx.Flags.String("hostnames") != "" {
	//         opts["hostnames"] = strings.Split(gctx.Flags.String("hostnames"), ",")
	// }
	if gctx.Flags.String("name") != "" {
		host.Name = gctx.Flags.String("name")
	}
	if gctx.Flags.String("info") != "" {
		host.Info = gctx.Flags.String("info")
	}
	if gctx.Flags.String("comment") != "" {
		host.Comment = gctx.Flags.String("comment")
	}

	// Update host
	host, err := remote.UpdateHost(host)
	if err != nil {
		fmt.Printf("%s[!]%s Error: %s\n",
			tui.RED, tui.RESET, err.Error())
	} else {
		fmt.Printf("%s*%s Updated host at: %s\n",
			tui.BLUE, tui.RESET, host.Addresses)
	}

}

func parseFilters(ctx *grumble.Context) (opts map[string]interface{}) {
	opts = make(map[string]interface{}, 0)

	if ctx.Flags.Uint("host_id") != 0 {
		opts["host_id"] = ctx.Flags.Uint("host_id")
	}
	if ctx.Flags.String("os_name") != "" {
		opts["os_name"] = ctx.Flags.String("os_name")
	}
	if ctx.Flags.String("os_family") != "" {
		opts["os_family"] = ctx.Flags.String("os_family")
	}
	if ctx.Flags.String("os_flavor") != "" {
		opts["os_flavor"] = ctx.Flags.String("os_flavor")
	}
	if ctx.Flags.String("os_sp") != "" {
		opts["os_sp"] = ctx.Flags.String("os_sp")
	}
	if ctx.Flags.String("arch") != "" {
		opts["arch"] = ctx.Flags.String("arch")
	}
	if ctx.Flags.String("addresses") != "" {
		opts["addresses"] = strings.Split(ctx.Flags.String("addresses"), ",")
	}
	if ctx.Flags.String("hostnames") != "" {
		opts["hostnames"] = strings.Split(ctx.Flags.String("hostnames"), ",")
	}
	if ctx.Flags.String("name") != "" {
		opts["name"] = ctx.Flags.String("name")
	}
	if ctx.Flags.String("info") != "" {
		opts["info"] = ctx.Flags.String("info")
	}
	if ctx.Flags.String("comment") != "" {
		opts["comment"] = ctx.Flags.String("comment")
	}
	if ctx.Flags.Bool("up") == true {
		opts["alive"] = ctx.Flags.Bool("up")
	}

	return opts
}
