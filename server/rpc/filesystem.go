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

package rpc

import (
	"time"

	"github.com/golang/protobuf/proto"

	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
)

func rpcLs(req []byte, timeout time.Duration, resp RPCResponse) {
	dirList := &ghostpb.LsReq{}
	err := proto.Unmarshal(req, dirList)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(dirList.GhostID)

	data, _ := proto.Marshal(&ghostpb.LsReq{
		Path: dirList.Path,
	})
	data, err = ghost.Request(ghostpb.MsgLsReq, timeout, data)
	resp(data, err)
}

func rpcRm(req []byte, timeout time.Duration, resp RPCResponse) {
	rmReq := &ghostpb.RmReq{}
	err := proto.Unmarshal(req, rmReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(rmReq.GhostID)

	data, _ := proto.Marshal(&ghostpb.RmReq{
		Path: rmReq.Path,
	})
	data, err = ghost.Request(ghostpb.MsgRmReq, timeout, data)
	resp(data, err)
}

func rpcMkdir(req []byte, timeout time.Duration, resp RPCResponse) {
	mkdirReq := &ghostpb.MkdirReq{}
	err := proto.Unmarshal(req, mkdirReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(mkdirReq.GhostID)

	data, _ := proto.Marshal(&ghostpb.MkdirReq{
		Path: mkdirReq.Path,
	})
	data, err = ghost.Request(ghostpb.MsgMkdirReq, timeout, data)
	resp(data, err)
}

func rpcCd(req []byte, timeout time.Duration, resp RPCResponse) {
	cdReq := &ghostpb.CdReq{}
	err := proto.Unmarshal(req, cdReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(cdReq.GhostID)

	data, _ := proto.Marshal(&ghostpb.CdReq{
		Path: cdReq.Path,
	})
	data, err = ghost.Request(ghostpb.MsgCdReq, timeout, data)
	resp(data, err)
}

func rpcPwd(req []byte, timeout time.Duration, resp RPCResponse) {
	pwdReq := &ghostpb.PwdReq{}
	err := proto.Unmarshal(req, pwdReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := (*core.Wire.Ghosts)[pwdReq.GhostID]

	data, _ := proto.Marshal(&ghostpb.PwdReq{})
	data, err = ghost.Request(ghostpb.MsgPwdReq, timeout, data)
	resp(data, err)
}

func rpcDownload(req []byte, timeout time.Duration, resp RPCResponse) {
	downloadReq := &ghostpb.DownloadReq{}
	err := proto.Unmarshal(req, downloadReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(downloadReq.GhostID)

	data, _ := proto.Marshal(&ghostpb.DownloadReq{
		Path: downloadReq.Path,
	})
	data, err = ghost.Request(ghostpb.MsgDownloadReq, timeout, data)
	resp(data, err)
}

func rpcUpload(req []byte, timeout time.Duration, resp RPCResponse) {
	uploadReq := &ghostpb.UploadReq{}
	err := proto.Unmarshal(req, uploadReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(uploadReq.GhostID)

	data, _ := proto.Marshal(&ghostpb.UploadReq{
		Encoder: uploadReq.Encoder,
		Path:    uploadReq.Path,
		Data:    uploadReq.Data,
	})
	data, err = ghost.Request(ghostpb.MsgUploadReq, timeout, data)
	resp(data, err)
}
