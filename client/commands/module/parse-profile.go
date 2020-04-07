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

// ModuleParseProfileCmd - parse a ghost profile into the current module settings
type ModuleParseProfileCmd struct {
	Positional struct {
		Name string `description:"Ghost profile to parse" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var ModuleParseProfile ModuleParseProfileCmd

func RegisterModuleParseProfile() {
	CommandParser.AddCommand(constants.ModuleParseProfile, "", "", &ModuleParseProfile)

	pp := CommandParser.Find(constants.ModuleParseProfile)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], pp)
	pp.ShortDescription = "parse a ghost profile into the current module settings"
	pp.Args()[0].RequiredMaximum = 1

}

// Execute - Parse a ghost profile into the current module settings
func (pp *ModuleParseProfileCmd) Execute(args []string) error {
	m := Context.Module

	run, _ := proto.Marshal(&clientpb.ModuleActionReq{
		WorkspaceID: uint32(Context.Workspace.ID),
		User:        Context.Server.Config.User,
		Path:        m.Path,
		Action:      constants.ModuleParseProfile,
		Profile:     pp.Positional.Name,
	})

	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgModuleReq,
		Data: run,
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s", resp.Err)
		return nil
	}

	result := &clientpb.ModuleAction{}
	proto.Unmarshal(resp.Data, result)

	if result.Success == false {
		fmt.Printf(Error+"%s", result.Err)
	} else {
		*m = *result.Updated
		fmt.Printf(Info+"%s", result.Result)
	}

	return nil
}
