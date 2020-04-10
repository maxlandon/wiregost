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
	"fmt"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// ModuleToProfileCmd - Generate a ghost profile with module current settings
type ModuleToProfileCmd struct {
	Positional struct {
		Name string `description:"Ghost profile name to generate" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var ModuleToProfile ModuleToProfileCmd

func RegisterModuleToProfile() {
	MainParser.AddCommand(constants.ModuleToProfile, "", "", &ModuleToProfile)

	tp := MainParser.Find(constants.ModuleToProfile)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], tp)
	tp.ShortDescription = "Generate a ghost profile with module current settings"
	tp.Args()[0].RequiredMaximum = 1

}

// Execute - Generate a ghost profile with module current settings
func (tp *ModuleToProfileCmd) Execute(args []string) error {
	m := Context.Module

	run, _ := proto.Marshal(&clientpb.ModuleActionReq{
		WorkspaceID: uint32(Context.Workspace.ID),
		User:        Context.Server.Config.User,
		Path:        m.Path,
		Action:      constants.ModuleToProfile,
		Profile:     tp.Positional.Name,
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
		fmt.Printf(Error+"%s \n", result.Err)
	} else {
		fmt.Printf(Info+"%s \n", result.Result)
	}

	return nil
}
