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
	"os"
	"strings"
	"text/tabwriter"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/util"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/util"
)

func RegisterFileSystemCommands() {

	ls := &Command{
		Name: "ls",
		Handle: func(r *Request) error {
			rpc := r.context.Server.RPC
			fmt.Println(tui.RESET)
			if len(r.Args) < 1 {
				r.Args = append(r.Args, ".")
			}

			data, _ := proto.Marshal(&ghostpb.LsReq{
				GhostID: r.context.CurrentAgent.ID,
				Path:    r.Args[0],
			})
			resp := <-rpc(&ghostpb.Envelope{
				Type: ghostpb.MsgLsReq,
				Data: data,
			}, defaultTimeout)
			if resp.Err != "" {
				fmt.Printf(RPCError+"%s\n", resp.Err)
				return nil
			}

			dirList := &ghostpb.Ls{}
			err := proto.Unmarshal(resp.Data, dirList)
			if err != nil {
				fmt.Printf(Error+"Unmarshaling envelope error: %v\n", err)
				return nil
			}
			printDirList(dirList)
			return nil
		},
	}
	AddCommand("agent", ls)

	cd := &Command{
		Name: "cd",
		Handle: func(r *Request) error {
			fmt.Println(tui.RESET)
			rpc := r.context.Server.RPC
			if len(r.Args) == 0 {
				fmt.Printf(Warn + "Missing parameter: file or directory name\n")
				return nil
			}

			data, _ := proto.Marshal(&ghostpb.CdReq{
				GhostID: r.context.CurrentAgent.ID,
				Path:    r.Args[0],
			})
			resp := <-rpc(&ghostpb.Envelope{
				Type: ghostpb.MsgCdReq,
				Data: data,
			}, defaultTimeout)
			if resp.Err != "" {
				fmt.Printf(RPCError+"%s\n", resp.Err)
				return nil
			}

			pwd := &ghostpb.Pwd{}
			err := proto.Unmarshal(resp.Data, pwd)
			if err != nil {
				fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
				return nil
			}
			fmt.Printf(Info+"%s\n", pwd.Path)

			// Update prompt
			*r.context.AgentPwd = pwd.Path
			return nil
		},
	}
	AddCommand("agent", cd)

	rm := &Command{
		Name: "rm",
		Handle: func(r *Request) error {
			fmt.Println(tui.RESET)
			rpc := r.context.Server.RPC
			if len(r.Args) == 0 {
				fmt.Printf(Warn + "Missing parameter: file or directory name\n")
				return nil
			}

			data, _ := proto.Marshal(&ghostpb.RmReq{
				GhostID: r.context.CurrentAgent.ID,
				Path:    r.Args[0],
			})
			resp := <-rpc(&ghostpb.Envelope{
				Type: ghostpb.MsgRmReq,
				Data: data,
			}, defaultTimeout)
			if resp.Err != "" {
				fmt.Printf(RPCError+"%s\n", resp.Err)
				return nil
			}

			rm := &ghostpb.Rm{}
			err := proto.Unmarshal(resp.Data, rm)
			if err != nil {
				fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
				return nil
			}
			if rm.Success {
				fmt.Printf(Info+"%s\n", rm.Path)
			} else {
				fmt.Printf(Warn+"%s\n", rm.Err)
			}

			return nil
		},
	}
	AddCommand("agent", rm)

	mkdir := &Command{
		Name: "mkdir",
		Handle: func(r *Request) error {

			return nil
		},
	}
	AddCommand("agent", mkdir)

	pwd := &Command{
		Name: "pwd",
		Handle: func(r *Request) error {
			fmt.Println(tui.RESET)
			rpc := r.context.Server.RPC
			data, _ := proto.Marshal(&ghostpb.PwdReq{
				GhostID: r.context.CurrentAgent.ID,
			})
			resp := <-rpc(&ghostpb.Envelope{
				Type: ghostpb.MsgPwdReq,
				Data: data,
			}, defaultTimeout)
			if resp.Err != "" {
				fmt.Printf(RPCError+"%s\n", resp.Err)
				return nil
			}

			pwd := &ghostpb.Pwd{}
			err := proto.Unmarshal(resp.Data, pwd)
			if err != nil {
				fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
				return nil
			}
			fmt.Printf(Info+"%s\n", pwd.Path)
			return nil
		},
	}
	AddCommand("agent", pwd)

	cat := &Command{
		Name: "cat",
		Handle: func(r *Request) error {

			return nil
		},
	}
	AddCommand("agent", cat)

	download := &Command{
		Name: "download",
		Handle: func(r *Request) error {

			return nil
		},
	}
	AddCommand("agent", download)

	upload := &Command{
		Name: "upload",
		Handle: func(r *Request) error {

			return nil
		},
	}
	AddCommand("agent", upload)
}

func printDirList(dirList *ghostpb.Ls) {
	fmt.Printf("%s\n", dirList.Path)
	fmt.Printf("%s\n", strings.Repeat("=", len(dirList.Path)))

	table := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	for _, fileInfo := range dirList.Files {
		if fileInfo.IsDir {
			fmt.Fprintf(table, "%s\t<dir>\t\n", fileInfo.Name)
		} else {
			fmt.Fprintf(table, "%s\t%s\t\n", fileInfo.Name, util.ByteCountBinary(fileInfo.Size))
		}
	}
	table.Flush()
}

func agentPwd(name string, rpc RPCServer) string {
	ghost := getGhost(name, rpc)
	data, _ := proto.Marshal(&ghostpb.PwdReq{
		GhostID: ghost.ID,
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgPwdReq,
		Data: data,
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return ""
	}

	pwd := &ghostpb.Pwd{}
	err := proto.Unmarshal(resp.Data, pwd)
	if err != nil {
		fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
		return ""
	}

	return pwd.Path
}
