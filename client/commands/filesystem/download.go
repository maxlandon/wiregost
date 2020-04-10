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
	"os"
	"path"
	"path/filepath"

	"github.com/gogo/protobuf/proto"
	"gopkg.in/AlecAivazis/survey.v1"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/spin"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/util"
)

// ListDirectories - Make a remote directory
type DownloadCmd struct {
	Positional struct {
		RemotePath string `description:"Remote directory/file to download" required:"1"`
		LocalPath  string `description:"Local path in which to save the file" required:"1"`
	} `positional-args:"yes"`
}

var Download DownloadCmd

func RegisterGhostDownload() {
	GhostParser.AddCommand(constants.GhostDownload, "", "", &Download)

	d := GhostParser.Find(constants.GhostDownload)
	d.ShortDescription = "Download a remote directory/file"
	d.Args()[0].RequiredMaximum = 1
	d.Args()[1].RequiredMaximum = 1
}

// Execute - Command
func (d *DownloadCmd) Execute(args []string) error {

	rpc := Context.Server.RPC

	src := d.Positional.RemotePath
	fileName := filepath.Base(src)
	dst, _ := filepath.Abs(d.Positional.LocalPath)
	fi, err := os.Stat(dst)
	if err != nil {
		fmt.Printf(Warn+"%v\n", err)
		return nil
	}
	if fi.IsDir() {
		dst = path.Join(dst, fileName)
	}

	if _, err := os.Stat(dst); err == nil {
		overwrite := false
		prompt := &survey.Confirm{Message: "Overwrite local file?"}
		survey.AskOne(prompt, &overwrite, nil)
		if !overwrite {
			return nil
		}
	}

	ctrl := make(chan bool)
	go spin.Until(fmt.Sprintf("%s -> %s", fileName, dst), ctrl)
	data, _ := proto.Marshal(&ghostpb.DownloadReq{
		GhostID: Context.Ghost.ID,
		Path:    d.Positional.RemotePath,
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgDownloadReq,
		Data: data,
	}, DefaultTimeout)
	ctrl <- true
	<-ctrl
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	download := &ghostpb.Download{}
	proto.Unmarshal(resp.Data, download)
	if download.Encoder == "gzip" {
		download.Data, _ = new(util.Gzip).Decode(download.Data)
	}
	f, err := os.Create(dst)
	if err != nil {
		fmt.Printf(Error+"Failed to open local file %s: %v\n", dst, err)
	}
	defer f.Close()
	n, err := f.Write(download.Data)
	if err != nil {
		fmt.Printf(Error+"Failed to write data %v\n", err)
	} else {
		fmt.Printf(Info+"Wrote %d bytes to %s\n", n, dst)
	}

	return nil
}
