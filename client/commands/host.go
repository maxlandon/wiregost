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

// import (
//         "context"
//         "fmt"
//         "regexp"
//         "strconv"
//         "strings"
//
//         "github.com/olekukonko/tablewriter"
//
//         "github.com/maxlandon/wiregost/client/constants"
//         "github.com/maxlandon/wiregost/client/help"
//         "github.com/maxlandon/wiregost/client/util"
//         "github.com/maxlandon/wiregost/data-service/models"
//         "github.com/maxlandon/wiregost/data-service/remote"
// )
//
// func registerHostCommands() {
//
//         // Declare all commands, subcommands and arguments
//         hosts := &Command{
//                 Name: "hosts",
//                 Help: "Manage database hosts (search/add/update/delete)",
//                 // SubCommands: []string{
//                 //         "search",
//                 //         "add",
//                 //         "update",
//                 //         "delete",
//                 // },
//                 Args: []*Arg{
//                         &Arg{Name: "host-id", Type: "uint"},
//                         &Arg{Name: "addresses", Type: "string"},
//                         &Arg{Name: "hostnames", Type: "string"},
//                         &Arg{Name: "os-name", Type: "string"},
//                         &Arg{Name: "os-family", Type: "string"},
//                         &Arg{Name: "os-flavor", Type: "string"},
//                         &Arg{Name: "os-sp", Type: "string"},
//                         &Arg{Name: "arch", Type: "string"},
//                         &Arg{Name: "name", Type: "string"},
//                         &Arg{Name: "info", Type: "string"},
//                         &Arg{Name: "comment", Type: "string"},
//                         &Arg{Name: "up", Type: "boolean"},
//                 },
//                 Handle: func(r *Request) error {
//                         switch length := len(r.Args); {
//                         // No arguments: Print hosts
//                         case length == 0:
//                                 fmt.Println()
//                                 hosts(&r.context.DBContext, nil)
//                         // Arguments: commands entered
//                         case length >= 1:
//                                 switch r.Args[0] {
//                                 case "search":
//                                         fmt.Println()
//                                         hosts(&r.context.DBContext, r.Args[1:])
//                                 case "add":
//                                         fmt.Println()
//                                         addHost(&r.context.DBContext, r.Args[1:])
//                                         fmt.Println()
//                                 case "delete":
//                                         fmt.Println()
//                                         deleteHosts(&r.context.DBContext, r.Args[1:])
//                                         fmt.Println()
//                                 case "update":
//                                         fmt.Println()
//                                         updateHost(&r.context.DBContext, r.Args[1:])
//                                         fmt.Println()
//                                 // No actions were asked for, list hosts with filters
//                                 default:
//                                         fmt.Println()
//                                         hosts(&r.context.DBContext, r.Args)
//                                 }
//
//                         }
//
//                         return nil
//                 },
//         }
//
//         // Add commands for each context
//         AddCommand("main", hosts)
//         AddCommand("module", hosts)
//         AddCommand("ghost", hosts)
// }
//
// func hosts(ctx *context.Context, options []string) {
//         var hosts []models.Host
//         var err error
//         opts := hostFilters(options)
//
//         if len(opts) == 0 {
//                 hosts, err = remote.Hosts(*ctx, nil)
//                 if err != nil {
//                         fmt.Printf(DBError+"%s", err.Error())
//                         return
//                 }
//         } else {
//                 hosts, err = remote.Hosts(*ctx, opts)
//                 if err != nil {
//                         fmt.Printf(DBError+"%s", err.Error())
//                         return
//                 }
//         }
//         hostsTable(ctx, &hosts)
// }
//
// func addHost(cctx *context.Context, args []string) {
//
//         opts := hostFilters(args)
//         if len(opts) == 0 {
//                 fmt.Printf(help.GetHelpFor(constants.HostsAdd))
//                 return
//         }
//
//         // Get existing host for comparing ReportHost() results
//         hosts, err := remote.Hosts(*cctx, nil)
//
//         host, err := remote.ReportHost(*cctx, opts)
//         if err != nil {
//                 fmt.Printf(DBError+"%s", err.Error())
//                 return
//         }
//
//         for _, h := range hosts {
//                 if h.ID == host.ID {
//                         fmt.Printf(Warn+"Host already exists at: %s", host.Addresses)
//                         return
//                 }
//         }
//
//         fmt.Printf(Info, "New host at: %s", host.Addresses)
// }
//
// func deleteHosts(cctx *context.Context, args []string) {
//         opts := hostFilters(args)
//         if len(opts) == 0 {
//                 fmt.Printf(help.GetHelpFor(constants.HostsDelete))
//                 return
//         }
//
//         // Get a list of hosts matching filters given
//         switch len(opts) {
//         case 0:
//                 fmt.Printf(Error + "Provide filters for host selection ")
//         default:
//                 var hosts []models.Host
//                 var err error
//
//                 hosts, err = remote.Hosts(*cctx, opts)
//                 if err != nil {
//                         fmt.Printf(DBError+"%s", err.Error())
//                         return
//                 }
//                 if len(hosts) == 0 {
//                         fmt.Printf(Warn + "No hosts match the given filters ")
//                         return
//                 }
//
//                 for i := range hosts {
//                         opts["host_id"] = []uint{hosts[i].ID}
//                         err = remote.DeleteHosts(*cctx, opts)
//                         if err != nil {
//                                 fmt.Printf(DBError, "%s", err.Error())
//                                 continue
//                         } else {
//                                 fmt.Printf(Info+"Deleted host at: %s", hosts[i].Addresses)
//                         }
//                 }
//         }
// }
//
// func updateHost(cctx *context.Context, args []string) {
//
//         var host *models.Host
//         opts := hostFilters(args)
//         if len(opts) == 0 {
//                 fmt.Printf(help.GetHelpFor(constants.HostsUpdate))
//                 return
//         }
//
//         if opts["host_id"] != nil {
//                 ids := opts["host_id"].([]uint)
//
//                 hosts, _ := remote.Hosts(*cctx, nil)
//                 for i := range hosts {
//                         for _, u := range ids {
//                                 if hosts[i].ID == u {
//                                         host = &hosts[i]
//                                 }
//                         }
//                 }
//         } else {
//                 fmt.Printf(Error + "Provide a host ID (host-id=2)")
//                 return
//         }
//
//         // Populate host with filters
//         osName, found := opts["os_name"]
//         if found {
//                 host.OSName = osName.(string)
//         }
//         osFamily, found := opts["os_family"]
//         if found {
//                 host.OSFamily = osFamily.(string)
//         }
//         osFlavor, found := opts["os_flavor"]
//         if found {
//                 host.OSFlavor = osFlavor.(string)
//         }
//         osSp, found := opts["os_sp"]
//         if found {
//                 host.OSSp = osSp.(string)
//         }
//         arch, found := opts["arch"]
//         if found {
//                 host.Arch = arch.(string)
//         }
//         name, found := opts["name"]
//         if found {
//                 host.Name = name.(string)
//         }
//         info, found := opts["info"]
//         if found {
//                 host.Info = info.(string)
//         }
//         comment, found := opts["comment"]
//         if found {
//                 host.Comment = comment.(string)
//         }
//         addrs, found := opts["addresses"]
//         if found {
//                 addrsString := strings.Split(addrs.(string), ",")
//                 var addrs []models.Address
//                 for _, a := range addrsString {
//                         addr := models.Address{Addr: a, HostID: host.ID, AddrType: "IPv4"}
//                         addrs = append(addrs, addr)
//                 }
//                 host.Addresses = addrs
//         }
//
//         // Update host
//         updated, err := remote.UpdateHost(host)
//         if err != nil {
//                 fmt.Printf(DBError+"%s", err.Error())
//         } else {
//                 fmt.Printf(Info+"Updated host at: %s", updated.Addresses)
//         }
// }
//
// func hostsTable(cctx *context.Context, hosts *[]models.Host) {
//         table := util.Table()
//         table.SetHeader([]string{"ID", "Addresses", "Name", "OS Name", "OS Flavor", "OS SP", "Arch", "Purpose", "Info", "Comments"})
//         table.SetColWidth(60)
//         table.SetColMinWidth(1, 15)
//         // table.SetColMinWidth(8, 20)
//         table.SetColMinWidth(9, 60)
//         table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//         )
//
//         data := [][]string{}
//         for _, h := range *hosts {
//                 addrsList := []string{}
//
//                 for _, a := range h.Addresses {
//                         addrsList = append(addrsList, a.String())
//                 }
//                 addrs := ""
//                 if len(addrsList) == 1 {
//                         addrs = addrsList[0]
//                 } else {
//                         addrs = strings.Join(addrsList, " ")
//                 }
//
//                 hostID := uint64(h.ID)
//                 data = append(data, []string{strconv.FormatUint(hostID, 10), addrs, h.Name,
//                         h.OSName, h.OSFlavor, h.OSSp, h.Arch, h.Purpose, h.Info, h.Comment})
//         }
//         table.AppendBulk(data)
//         table.Render()
// }
//
// func hostFilters(args []string) (opts map[string]interface{}) {
//         opts = make(map[string]interface{}, 0)
//
//         for _, arg := range args {
//
//                 // Host ID
//                 if strings.Contains(arg, "host-id") {
//                         vals := strings.Split(arg, "=")[1]
//                         ids := strings.Split(vals, ",")
//                         var uIds []uint
//                         for _, id := range ids {
//                                 uID, _ := strconv.Atoi(id)
//                                 uIds = append(uIds, uint(uID))
//                         }
//                         opts["host_id"] = uIds
//                 }
//
//                 // Host OS Properties
//                 if strings.Contains(arg, "os-name") {
//                         vals := strings.Split(arg, "=")
//                         // names := strings.Split(vals, ",")
//                         opts["os_name"] = vals[1]
//                 }
//                 if strings.Contains(arg, "os-family") {
//                         vals := strings.Split(arg, "=")
//                         opts["os_family"] = vals[1]
//                 }
//                 if strings.Contains(arg, "os-flavor") {
//                         vals := strings.Split(arg, "=")
//                         opts["os_flavor"] = vals[1]
//                 }
//                 if strings.Contains(arg, "os-sp") {
//                         vals := strings.Split(arg, "=")
//                         opts["os_sp"] = vals[1]
//                 }
//
//                 // Host CPU
//                 if strings.Contains(arg, "arch") {
//                         vals := strings.Split(arg, "=")
//                         opts["arch"] = vals[1]
//                 }
//
//                 // Host addresses
//                 if strings.Contains(arg, "addresses") {
//                         vals := strings.Split(arg, "=")[1]
//                         addrs := strings.Split(vals, ",")
//                         opts["addresses"] = addrs
//                 }
//
//                 // Host names, users
//                 if strings.Contains(arg, "hostnames") {
//                         vals := strings.Split(arg, "=")[1]
//                         hostnames := strings.Split(vals, ",")
//                         opts["hostnames"] = hostnames
//                 }
//
//                 if strings.Contains(arg, "name") {
//                         vals := strings.Split(arg, "=")[1]
//                         names := strings.Split(vals, ",")
//                         opts["name"] = names
//                 }
//
//                 // Host state
//                 if strings.Contains(arg, "up") {
//                         vals := strings.Split(arg, "=")
//                         opts["alive"], _ = strconv.ParseBool(vals[1])
//                 }
//
//                 // Host info
//                 if strings.Contains(arg, "info") {
//                         desc := regexp.MustCompile(`\b(info){1}.*"`)
//                         result := desc.FindStringSubmatch(strings.Join(args, " "))
//                         opts["info"] = strings.Trim(strings.TrimPrefix(result[0], "info="), "\"")
//                 }
//                 if strings.Contains(arg, "comment") {
//                         desc := regexp.MustCompile(`\b(comment){1}.*"`)
//                         result := desc.FindStringSubmatch(strings.Join(args, " "))
//                         opts["comment"] = strings.Trim(strings.TrimPrefix(result[0], "comment="), "\"")
//                 }
//         }
//
//         return opts
// }
