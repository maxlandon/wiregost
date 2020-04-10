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

package filesystem

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// ListDirectories - Print the ghost implant working directory
type PwdCmd struct{}

var Pwd PwdCmd

func RegisterGhostPwd() {
	GhostParser.AddCommand(constants.GhostPwd, "", "", &Pwd)

	rm := GhostParser.Find(constants.GhostPwd)
	rm.ShortDescription = "Print the ghost implant working directory"
}

// Execute - Command
func (r *PwdCmd) Execute(args []string) error {

	rpc := Context.Server.RPC
	data, _ := proto.Marshal(&ghostpb.PwdReq{
		GhostID: Context.Ghost.ID,
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgPwdReq,
		Data: data,
	}, DefaultTimeout)
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
}

// func agentPwd(name string, rpc RPCServer) string {
//         ghost := getGhost(name, rpc)
//         data, _ := proto.Marshal(&ghostpb.PwdReq{
//                 GhostID: ghost.ID,
//         })
//         resp := <-rpc(&ghostpb.Envelope{
//                 Type: ghostpb.MsgPwdReq,
//                 Data: data,
//         }, defaultTimeout)
//         if resp.Err != "" {
//                 fmt.Printf(RPCError+"%s\n", resp.Err)
//                 return ""
//         }
//
//         pwd := &ghostpb.Pwd{}
//         err := proto.Unmarshal(resp.Data, pwd)
//         if err != nil {
//                 fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
//                 return ""
//         }
//
//         return pwd.Path
// }
