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
//         "fmt"
//         "strconv"
//         "strings"
//
//         "github.com/evilsocket/islazy/tui"
//         "github.com/gogo/protobuf/proto"
//
//         clientpb "github.com/maxlandon/wiregost/protobuf/client"
//         ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
// )
//
// func registerAgentInfoCommands() {
//
//         info := &Command{
//                 Name: "info",
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         info(r.Args, *r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", info)
//
//         getpid := &Command{
//                 Name: "getpid",
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         getPID(*r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", getpid)
//
//         getuid := &Command{
//                 Name: "getuid",
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         getUID(*r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", getuid)
//
//         getgid := &Command{
//                 Name: "getgid",
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         getGID(*r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", getgid)
//
//         whoami := &Command{
//                 Name: "whoami",
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         whoami(*r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", whoami)
//
//         ifconfig := &Command{
//                 Name: "ifconfig",
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         ifconfig(*r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", ifconfig)
// }
//
// func info(args []string, ctx ShellContext, rpc RPCServer) {
//
//         var ghost *clientpb.Ghost
//         if ctx.Ghost.Name != "" {
//                 ghost = ctx.Ghost
//         } else if 0 < len(args) {
//                 ghost = getGhost(args[0], rpc)
//         }
//
//         fmt.Printf(" %s%sGhost Implant %s%s \n", tui.BOLD, tui.BLUE, ghost.Name, tui.RESET)
//         fmt.Println()
//
//         if ghost.Name != "" {
//                 fmt.Printf(tui.BLUE+"            ID: %s%d\n", tui.FOREWHITE, ghost.ID)
//                 fmt.Printf(tui.BLUE+"          Name: %s%s\n", tui.FOREWHITE, ghost.Name)
//                 fmt.Printf(tui.BLUE+"      Hostname: %s%s\n", tui.FOREWHITE, ghost.Hostname)
//                 fmt.Printf(tui.BLUE+"      Username: %s%s\n", tui.FOREWHITE, ghost.Username)
//                 fmt.Printf(tui.BLUE+"           UID: %s%s\n", tui.FOREWHITE, ghost.UID)
//                 fmt.Printf(tui.BLUE+"           GID: %s%s\n", tui.FOREWHITE, ghost.GID)
//                 fmt.Printf(tui.BLUE+"           PID: %s%d\n", tui.FOREWHITE, ghost.PID)
//                 fmt.Printf(tui.BLUE+"            OS: %s%s\n", tui.FOREWHITE, ghost.OS)
//                 fmt.Printf(tui.BLUE+"       Version: %s%s\n", tui.FOREWHITE, ghost.Version)
//                 fmt.Printf(tui.BLUE+"          Arch: %s%s\n", tui.FOREWHITE, ghost.Arch)
//                 fmt.Printf(tui.BLUE+"Remote Address: %s%s\n", tui.FOREWHITE, ghost.RemoteAddress)
//                 fmt.Printf(tui.BLUE+"  Last Checkin: %s%s\n", tui.FOREWHITE, ghost.LastCheckin)
//
//         } else {
//                 fmt.Printf(Warn + "No target Ghost, see `help`\n")
//         }
// }
//
// func ping(ctx ShellContext, rpc RPCServer) {
//         if ctx.Ghost.Name == "" {
//                 fmt.Printf(Warn + "Please select an active Ghost implant via `interact`\n")
//                 return
//         }
// }
//
// func getPID(ctx ShellContext, rpc RPCServer) {
//         if ctx.Ghost.Name == "" {
//                 fmt.Printf(Warn + "Please select an active Ghost implant via `interact`\n")
//                 return
//         }
//         fmt.Printf("%d\n", ctx.Ghost.PID)
// }
//
// func getUID(ctx ShellContext, rpc RPCServer) {
//         if ctx.Ghost.Name == "" {
//                 fmt.Printf(Warn + "Please select an active Ghost implant via `interact`\n")
//                 return
//         }
//         fmt.Printf("%s\n", ctx.Ghost.UID)
// }
//
// func getGID(ctx ShellContext, rpc RPCServer) {
//         if ctx.Ghost.Name == "" {
//                 fmt.Printf(Warn + "Please select an active Ghost implant via `interact`\n")
//                 return
//         }
//         fmt.Printf("%s\n", ctx.Ghost.GID)
// }
//
// func whoami(ctx ShellContext, rpc RPCServer) {
//         if ctx.Ghost.Name == "" {
//                 fmt.Printf(Warn + "Please select an active Ghost implant via `interact`\n")
//                 return
//         }
//         fmt.Printf("%s\n", ctx.Ghost.Username)
// }
//
// func ifconfig(ctx ShellContext, rpc RPCServer) {
//
//         data, _ := proto.Marshal(&ghostpb.IfconfigReq{GhostID: ctx.Ghost.ID})
//         resp := <-rpc(&ghostpb.Envelope{
//                 Type: ghostpb.MsgIfconfigReq,
//                 Data: data,
//         }, defaultTimeout)
//         if resp.Err != "" {
//                 fmt.Printf(RPCError+"%s\n", resp.Err)
//                 return
//         }
//
//         ifaceConfigs := &ghostpb.Ifconfig{}
//         err := proto.Unmarshal(resp.Data, ifaceConfigs)
//         if err != nil {
//                 fmt.Printf(Error + "Failed to decode response\n")
//                 return
//         }
//
//         for ifaceIndex, iface := range ifaceConfigs.NetInterfaces {
//                 fmt.Printf("%s%s%s (%d)\n", tui.BOLD, iface.Name, tui.RESET, ifaceIndex)
//                 if 0 < len(iface.MAC) {
//                         fmt.Printf("   MAC Address: %s\n", iface.MAC)
//                 }
//                 for _, ip := range iface.IPAddresses {
//
//                         // Try to find local IPs and colorize them
//                         subnet := -1
//                         if strings.Contains(ip, "/") {
//                                 parts := strings.Split(ip, "/")
//                                 subnetStr := parts[len(parts)-1]
//                                 subnet, err = strconv.Atoi(subnetStr)
//                                 if err != nil {
//                                         subnet = -1
//                                 }
//                         }
//
//                         if 0 < subnet && subnet <= 32 && !isLoopback(ip) {
//                                 fmt.Printf(tui.BLUE+"    IP Address: %s%s%s\n", tui.FOREWHITE, ip, tui.RESET)
//                         } else if 32 < subnet && !isLoopback(ip) {
//                                 fmt.Printf(tui.BLUE+"    IP Address: %s%s%s\n", tui.FOREWHITE, ip, tui.RESET)
//                         } else {
//                                 fmt.Printf("    IP Address: %s\n", ip)
//                         }
//                 }
//         }
// }
//
// func isLoopback(ip string) bool {
//         if strings.HasPrefix(ip, "127") || strings.HasPrefix(ip, "::1") {
//                 return true
//         }
//         return false
// }
