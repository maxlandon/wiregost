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

// ListDirectories - Remove a remote directory/file
type RmCmd struct {
	Positional struct {
		Path string `description:"Remote directory/file" required:"1"`
	} `positional-args:"yes"`
}

var Rm RmCmd

func RegisterGhostRm() {
	GhostParser.AddCommand(constants.GhostRm, "", "", &Rm)

	rm := GhostParser.Find(constants.GhostRm)
	rm.ShortDescription = "Remove a remote directory/file"
	rm.Args()[0].RequiredMaximum = 1
}

// Execute - Command
func (r *RmCmd) Execute(args []string) error {

	rpc := Context.Server.RPC

	data, _ := proto.Marshal(&ghostpb.RmReq{
		GhostID: Context.Ghost.ID,
		Path:    r.Positional.Path,
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgRmReq,
		Data: data,
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	rm := &ghostpb.Rm{}
	err := proto.Unmarshal(resp.Data, rm)
	if err != nil {
		fmt.Printf(Error+"Unmarshaling envelope error: %v\n", err)
		return nil
	}
	if rm.Success {
		fmt.Printf(Info+"%s\n", rm.Path)
	} else {
		fmt.Printf(Warn+"%s\n", rm.Err)
	}

	return nil
}
