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
	"path/filepath"

	"github.com/gogo/protobuf/proto"
	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// ChangeDirectory - Change the working directory of the client console
type ChangeDirectoryCmd struct {
	Positional struct {
		Path string `description:"Remote path" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var ChangeDirectory ChangeDirectoryCmd

func RegisterGhostCd() {
	// cd
	GhostParser.AddCommand(constants.GhostCd, "", "", &ChangeDirectory)

	cd := GhostParser.Find(constants.GhostCd)
	cd.ShortDescription = "Change the ghost implant working directory"

	cd.Args()[0].RequiredMaximum = 1
}

// Execute - Handler for ChangeDirectory
func (cd *ChangeDirectoryCmd) Execute(args []string) error {

	rpc := Context.Server.RPC

	path := cd.Positional.Path
	if (path == "~" || path == "~/") && Context.Ghost.OS == "linux" {
		path = filepath.Join("/home", Context.Ghost.Username)
	}

	data, _ := proto.Marshal(&ghostpb.CdReq{
		GhostID: Context.Ghost.ID,
		Path:    path,
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgCdReq,
		Data: data,
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	pwd := &ghostpb.Pwd{}
	err := proto.Unmarshal(resp.Data, pwd)
	if err != nil {
		fmt.Printf(Error+"Unmarshaling envelope error: %v\n", err)
		return nil
	}
	fmt.Printf(Info+"%s\n", pwd.Path)

	// Update prompt
	*Context.GhostPwd = pwd.Path

	return nil
}
