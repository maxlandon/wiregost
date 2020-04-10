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

// ListDirectories - Make a remote directory
type MkdirCmd struct {
	Positional struct {
		Path string `description:"Remote directory to create" required:"1"`
	} `positional-args:"yes"`
}

var Mkdir MkdirCmd

func RegisterGhostMkdir() {
	GhostParser.AddCommand(constants.GhostMkdir, "", "", &Mkdir)

	m := GhostParser.Find(constants.GhostMkdir)
	m.ShortDescription = "Make a remote directory"
	m.Args()[0].RequiredMaximum = 1
}

// Execute - Command
func (m *MkdirCmd) Execute(args []string) error {

	rpc := Context.Server.RPC

	data, _ := proto.Marshal(&ghostpb.MkdirReq{
		GhostID: Context.Ghost.ID,
		Path:    m.Positional.Path,
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgMkdirReq,
		Data: data,
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	mkdir := &ghostpb.Mkdir{}
	err := proto.Unmarshal(resp.Data, mkdir)
	if err != nil {
		fmt.Printf(Error+"Unmarshaling envelope error: %v\n", err)
		return nil
	}
	if mkdir.Success {
		fmt.Printf(Info+"%s\n", mkdir.Path)
	} else {
		fmt.Printf(Warn+"%s\n", mkdir.Err)
	}

	return nil
}
