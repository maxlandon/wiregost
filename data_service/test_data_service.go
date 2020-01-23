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

package main

import (
	"context"
	"fmt"

	"github.com/maxlandon/wiregost/data_service/remote"
)

func main() {
	// --------------------------------------------------------------------
	// WORkSPACE TESTING
	// --------------------------------------------------------------------

	// List workspaces
	// workspaces, _ := remote.Workspaces()
	// fmt.Println(workspaces)

	// Add worspace
	// err := remote.AddWorkspaces(nil, []string{"Macron"})
	// if err != nil {
	//         log.Fatal("[!] error in request: " + err.Error())
	// }

	// Delete workspaces
	// err := remote.DeleteWorkspaces([]int{8674665223082153551})
	// if err != nil {
	//         log.Fatal("[!] error in request: " + err.Error())
	// }

	// update workspace
	// MacronWorkspace := workspaces[3]
	// MacronWorkspace.Description = "Test update in Macron's workspace"
	// err := remote.UpdateWorkspace(MacronWorkspace)
	// if err != nil {
	//         fmt.Println("Change return error values")
	// }

	// List workspaces
	workspaces, _ := remote.Workspaces(nil)
	fmt.Println(workspaces)

	// --------------------------------------------------------------------
	// HOST TESTING
	// --------------------------------------------------------------------
	// Create host
	// h := &models.Host{
	//         OSName:      "ArchLinux",
	//         WorkspaceID: 8674665223082153551,
	// }

	// List hosts
	parentCtx := context.Background()
	ctx := context.WithValue(parentCtx, "workspace_id", 3)

	opts := map[string]interface{}{"os_name": "Windows", "addresses": []string{"192.168.2.10"}}
	newHost, err := remote.ReportHost(ctx, opts)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(newHost)
	}
	//         newHost.OSName = "Linux"
	//         newHost.Addresses = append(newHost.Addresses, models.Address{HostID: newHost.ID, Addr: "220.198.230.14", AddrType: "IPv4"})
	//         newHost, err = remote.UpdateHost(newHost)
	// }

	// List hosts with search option filters
	// opts := map[string]interface{}{"os_name": "linux", "addresses": []string{"192.168.1.49", "192.168.1.18", "220.198.230.14"}}
	// // opts := map[string]interface{}{"addresses": []string{"192.168.1.49"}}
	// hosts, err := remote.Hosts(ctx, opts)
	// if len(hosts) == 0 {
	//         fmt.Println(err)
	// } else {
	//         updated := &hosts[0]
	//         fmt.Println(updated)
	// }
	//
	// Get a single host with HostAPIPath/search

	// singleHost, _ := remote.GetHost(ctx, map[string]interface{}{"host_id": 5})
	// fmt.Println(singleHost)

	// Update Host
	// updated.OSFlavor = "Vista"
	// updated.Arch = "x86_64"
	// updated, _ = remote.UpdateHost(updated)

	// updated.Addresses = append(updated.Addresses, models.Address{HostID: updated.ID, Addr: "220.198.230.14", AddrType: "IPv4"})

	// updated, _ = remote.UpdateHost(updated)
	// Delete Host
	// remote.DeleteHost(hosts[3].ID)

	// List hosts
	// newhosts, _ := remote.Hosts()
	// fmt.Println(newhosts)

	// Get a single host
	// opts := map[string]string{"host_id": "5577006791947779410"}
	// singleHost, _ := remote.GetHost(opts)
	//
	// fmt.Println(*singleHost)
}
