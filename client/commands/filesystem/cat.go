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
	"github.com/maxlandon/wiregost/util"
)

// ListDirectories - Print a file to std output (downloads the file first)
type CatCmd struct {
	Positional struct {
		Path string `description:"Remote directory/file" required:"1"`
	} `positional-args:"yes"`
}

var Cat CatCmd

func RegisterGhostCat() {
	GhostParser.AddCommand(constants.GhostCat, "", "", &Cat)

	c := GhostParser.Find(constants.GhostCat)
	c.ShortDescription = "Print a file to std output (downloads the file first)"
	c.Args()[0].RequiredMaximum = 1
}

// Execute - Command
func (r *CatCmd) Execute(args []string) error {

	rpc := Context.Server.RPC

	data, _ := proto.Marshal(&ghostpb.DownloadReq{
		GhostID: Context.Ghost.ID,
		Path:    r.Positional.Path,
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgDownloadReq,
		Data: data,
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	download := &ghostpb.Download{}
	proto.Unmarshal(resp.Data, download)
	if download.Encoder == "gzip" {
		download.Data, _ = new(util.Gzip).Decode(download.Data)
	}
	fmt.Printf(string(download.Data))

	return nil
}
