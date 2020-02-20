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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/maxlandon/wiregost/client/spin"
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
				fmt.Printf(Error+"Unmarshaling envelope error: %v\n", err)
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
				fmt.Printf(Error+"Unmarshaling envelope error: %v\n", err)
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
			fmt.Println(tui.RESET)
			rpc := r.context.Server.RPC
			if len(r.Args) == 0 {
				fmt.Printf(Warn + "Missing parameter: directory name\n")
				return nil
			}

			data, _ := proto.Marshal(&ghostpb.MkdirReq{
				GhostID: r.context.CurrentAgent.ID,
				Path:    r.Args[0],
			})
			resp := <-rpc(&ghostpb.Envelope{
				Type: ghostpb.MsgMkdirReq,
				Data: data,
			}, defaultTimeout)
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
			fmt.Println(tui.RESET)
			rpc := r.context.Server.RPC
			if len(r.Args) == 0 {
				fmt.Printf(Warn + "Missing parameter: file name\n")
				return nil
			}

			data, _ := proto.Marshal(&ghostpb.DownloadReq{
				GhostID: r.context.CurrentAgent.ID,
				Path:    r.Args[0],
			})
			resp := <-rpc(&ghostpb.Envelope{
				Type: ghostpb.MsgDownloadReq,
				Data: data,
			}, defaultTimeout)
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
		},
	}

	AddCommand("agent", cat)

	download := &Command{
		Name: "download",
		Handle: func(r *Request) error {
			fmt.Println(tui.RESET)
			rpc := r.context.Server.RPC
			if len(r.Args) < 1 {
				fmt.Println(Warn + "Missing parameter(s), see `help download`\n")
				return nil
			}

			src := r.Args[0]
			fileName := filepath.Base(src)
			dst, _ := filepath.Abs(r.Args[1])
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
				GhostID: r.context.CurrentAgent.ID,
				Path:    r.Args[0],
			})
			resp := <-rpc(&ghostpb.Envelope{
				Type: ghostpb.MsgDownloadReq,
				Data: data,
			}, defaultTimeout)
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
		},
	}

	AddCommand("agent", download)

	upload := &Command{
		Name: "upload",
		Handle: func(r *Request) error {
			fmt.Println(tui.RESET)
			rpc := r.context.Server.RPC
			if len(r.Args) < 1 {
				fmt.Println(Warn + "Missing parameter, see `help upload`\n")
				return nil
			}

			src, _ := filepath.Abs(r.Args[0])
			_, err := os.Stat(src)
			if err != nil {
				fmt.Printf(Error+"%v\n", err)
				return nil
			}

			if len(r.Args) == 1 {
				fileName := filepath.Base(src)
				r.Args = append(r.Args, fileName)
			}
			dst := r.Args[1]

			fileBuf, err := ioutil.ReadFile(src)
			uploadGzip := bytes.NewBuffer([]byte{})
			new(util.Gzip).Encode(uploadGzip, fileBuf)

			ctrl := make(chan bool)
			go spin.Until(fmt.Sprintf("%s -> %s", src, dst), ctrl)
			data, _ := proto.Marshal(&ghostpb.UploadReq{
				GhostID: r.context.CurrentAgent.ID,
				Path:    dst,
				Data:    uploadGzip.Bytes(),
				Encoder: "gzip",
			})
			resp := <-rpc(&ghostpb.Envelope{
				Type: ghostpb.MsgUploadReq,
				Data: data,
			}, defaultTimeout)
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
		},
	}

	AddCommand("agent", upload)
}

func printDirList(dirList *ghostpb.Ls) {
	fmt.Printf("Listing dir: %s%s%s\n", tui.BOLD, dirList.Path, tui.RESET)
	fmt.Printf("%s\n", strings.Repeat(tui.Dim("-"), len(dirList.Path)))

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
