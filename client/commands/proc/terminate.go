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

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// PsCmd - "List processes running on the target, with (--optional) filters"
type TerminateCmd struct {
	*ProcOptions
	Positional struct {
		PID int `description:"Process ID to terminate" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var Terminate TerminateCmd

func RegisterTerminate() {
	GhostParser.AddCommand(constants.Terminate, "", "", &Terminate)

	t := GhostParser.Find(constants.Terminate)
	t.ShortDescription = "Terminate a process running on the target"
}

// Execute - Command
func (p *TerminateCmd) Execute(args []string) error {

	data, _ := proto.Marshal(&ghostpb.TerminateReq{
		GhostID: Context.Ghost.ID,
		Pid:     int32(p.Positional.PID),
	})
	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: ghostpb.MsgTerminate,
		Data: data,
	}, DefaultTimeout)

	termResp := &ghostpb.Terminate{}
	err := proto.Unmarshal(resp.Data, termResp)
	if err != nil {
		fmt.Printf(Warn+"Error: %v\n", err)
		return nil
	}
	if termResp.Err != "" {
		fmt.Printf(Warn+"Error: %s\n", termResp.Err)
		return nil
	}
	fmt.Printf(Info+"Process %d has been terminated\n", p.Positional.PID)

	return nil
}
