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
	"fmt"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"
	"github.com/maxlandon/wiregost/client/help"
	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/olekukonko/tablewriter"
)

func RegisterUserCommands() {

	users := &Command{
		Name: "user",
		Help: help.GetHelpFor("user"),
		SubCommands: []string{
			"add",
		},
		Args: []*CommandArg{
			&CommandArg{Name: "name", Type: "string"},
			&CommandArg{Name: "lhost", Type: "string"},
			&CommandArg{Name: "lport", Type: "uint"},
			&CommandArg{Name: "default", Type: "boolean"},
		},
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				users(r.context.Server.RPC)
			case length >= 1:
				switch r.Args[0] {
				case "add":
					addUser(r.Args[1:], r.context.Server.RPC)
				}

			}
			return nil
		},
	}

	AddCommand("main", users)
	AddCommand("module", users)
}

func users(rpc RPCServer) {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgUser,
		Data: []byte{},
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf("%s[!] RPC Error:%s %s\n", tui.RED, tui.RESET, resp.Err)
		return
	}

	users := &clientpb.Players{}
	err := proto.Unmarshal(resp.Data, users)
	if err != nil {
		fmt.Println()
		fmt.Printf("%s[!]%s %s", tui.RED, tui.RESET, err.Error())
		fmt.Println()
		return
	}

	if 0 < len(users.Players) {
		displayUsers(users.Players)
	} else {
		fmt.Printf("%s*%s No players currently registered", tui.BLUE, tui.RESET)
	}

}

func displayUsers(users []*clientpb.Player) {
	table := util.Table()
	table.SetHeader([]string{"Name", "Status"})
	table.SetColWidth(40)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	for _, u := range users {
		status := ""
		if u.Online {
			status = fmt.Sprintf("%sonline%s", tui.GREEN, tui.RESET)
		} else {
			status = fmt.Sprintf("%soffline%s", tui.RED, tui.RESET)
		}
		table.Append([]string{u.Client.User, status})
	}
	table.Render()
}

func addUser(args []string, rpc RPCServer) {

	var name string
	var lhost string
	var lport uint32
	isDefault := false

	for _, arg := range args {
		if strings.Contains(arg, "name") {
			vals := strings.Split(arg, "=")
			name = vals[1]
		}
		if strings.Contains(arg, "lhost") {
			vals := strings.Split(arg, "=")
			lhost = vals[1]
		}
		if strings.Contains(arg, "lport") {
			vals := strings.Split(arg, "=")
			port, _ := strconv.Atoi(vals[1])
			lport = uint32(port)
		}
		if strings.Contains(arg, "default") {
			vals := strings.Split(arg, "=")
			if vals[1] == "true" {
				isDefault = true
			}
		}
	}

	if name == "" {
		fmt.Println()
		fmt.Printf("%s[!]%s Provide a user name (name='name')",
			tui.RED, tui.RESET)
		fmt.Println()
		return
	}
	if lhost == "" {
		fmt.Println()
		fmt.Printf("%s[!]%s Provide a lhost (lhost=192.168.1.1)",
			tui.RED, tui.RESET)
		fmt.Println()
		return
	}
	if lport == 0 {
		fmt.Println()
		fmt.Printf("%s[!]%s Provide a lport (lport=8443)",
			tui.RED, tui.RESET)
		fmt.Println()
		return
	}

	userReq, _ := proto.Marshal(&clientpb.UserReq{
		User:    name,
		LHost:   lhost,
		LPort:   lport,
		Default: isDefault,
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgUserReq,
		Data: userReq,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf("%s[!] RPC Error:%s %s\n", tui.RED, tui.RESET, resp.Err)
		return
	}

	userRes := &clientpb.User{}
	err := proto.Unmarshal(resp.Data, userRes)
	if err != nil {
		fmt.Println()
		fmt.Printf("%s[!]%s %s", tui.RED, tui.RESET, err.Error())
		fmt.Println()
		return
	}

	if userRes.Success {
		fmt.Println()
		fmt.Printf("%s[*]%s Added user %s with server %s:%d", tui.GREEN, tui.RESET, name, lhost, lport)
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Printf("%s[!]%s %s", tui.RED, tui.RESET, userRes.Err)
		fmt.Println()
	}
}
