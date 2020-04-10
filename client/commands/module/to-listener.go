// wiregost - golang exploitation framework
// copyright © 2020 para
//
// this program is free software: you can redistribute it and/or modify
// it under the terms of the gnu general public license as published by
// the free software foundation, either version 3 of the license, or
// (at your option) any later version.
//
// this program is distributed in the hope that it will be useful,
// but without any warranty; without even the implied warranty of
// merchantability or fitness for a particular purpose.  see the
// gnu general public license for more details.
//
// you should have received a copy of the gnu general public license
// along with this program.  if not, see <http://www.gnu.org/licenses/>.

package module

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// ToListenerCmd - Spawn a listener based on the module's options
type ToListenerCmd struct{}

var ToListener ToListenerCmd

func RegisterToListener() {
	MainParser.AddCommand(constants.ModuleToListener, "", "", &ToListener)

	listen := MainParser.Find(constants.ModuleToListener)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], listen)
	listen.ShortDescription = "Spawn a listener based on the module's options"
}

// Execute - Spawn a listener based on the module's options
func (tl *ToListenerCmd) Execute(args []string) error {

	m := Context.Module

	run, _ := proto.Marshal(&clientpb.ModuleActionReq{
		WorkspaceID: uint32(Context.Workspace.ID),
		User:        Context.Server.Config.User,
		Path:        m.Path,
		Action:      constants.ModuleToListener,
	})

	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgModuleReq,
		Data: run,
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s \n", resp.Err)
		return nil
	}

	result := &clientpb.ModuleAction{}
	proto.Unmarshal(resp.Data, result)

	if result.Success == false {
		fmt.Printf(Error+"%s\n", result.Err)
	} else {
		fmt.Printf(Success+"%s \n", result.Result)
	}

	return nil
}
