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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/spin"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/util"
)

// ListDirectories - Upload a local directory/file onto the ghost target system
type UploadCmd struct {
	Positional struct {
		LocalPath  string `description:"Local path to file to upload" required:"1"`
		RemotePath string `description:"Remote directory/file in which to upload"`
	} `positional-args:"yes"`
}

var Upload UploadCmd

func RegisterGhostUpload() {
	GhostParser.AddCommand(constants.GhostUpload, "", "", &Upload)

	u := GhostParser.Find(constants.GhostUpload)
	u.ShortDescription = "Upload a local directory/file onto the ghost target system"
	u.Args()[0].RequiredMaximum = 1
	u.Args()[1].RequiredMaximum = 1
}

// Execute - Command
func (u *UploadCmd) Execute(args []string) error {

	rpc := Context.Server.RPC

	src, _ := filepath.Abs(u.Positional.LocalPath)
	_, err := os.Stat(src)
	if err != nil {
		fmt.Printf(Error+"%v\n", err)
		return nil
	}

	var dst string
	if u.Positional.RemotePath == "" {
		dst = filepath.Base(src)
	} else {
		dst = u.Positional.RemotePath
	}

	fileBuf, err := ioutil.ReadFile(src)
	uploadGzip := bytes.NewBuffer([]byte{})
	new(util.Gzip).Encode(uploadGzip, fileBuf)

	ctrl := make(chan bool)
	go spin.Until(fmt.Sprintf("%s -> %s", src, dst), ctrl)
	data, _ := proto.Marshal(&ghostpb.UploadReq{
		GhostID: Context.Ghost.ID,
		Path:    dst,
		Data:    uploadGzip.Bytes(),
		Encoder: "gzip",
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgUploadReq,
		Data: data,
	}, DefaultTimeout)
	ctrl <- true
	<-ctrl
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	upload := &ghostpb.Upload{}
	err = proto.Unmarshal(resp.Data, upload)
	if err != nil {
		fmt.Printf(Error+"Unmarshaling envelope error: %v\n", err)
		return nil
	}
	if upload.Success {
		fmt.Printf(Info+"Written to %s\n", upload.Path)
	} else {
		fmt.Printf(Warn+"Error %s\n", upload.Err)
	}

	return nil
}
