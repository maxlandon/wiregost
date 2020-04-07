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

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// ModuleSetOptionCmd - Set a module's option
type ModuleSetOptionCmd struct {
	Positional struct {
		Option string `description:"Option name" required:"1"`
		Value  string `description:"Option value" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var ModuleSetOption ModuleSetOptionCmd

func RegisterModuleSetOption() {
	CommandParser.AddCommand(constants.ModuleSetOption, "", "", &ModuleSetOption)

	set := CommandParser.Find(constants.ModuleSetOption)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], set)
	set.ShortDescription = "Set a module option"
	set.Args()[0].RequiredMaximum = 1
	set.Args()[1].RequiredMaximum = 1
}

// Execute - Set a module's option
func (so *ModuleSetOptionCmd) Execute(args []string) error {

	if _, found := Context.Module.Options[so.Positional.Option]; !found {
		fmt.Printf(Error+"Invalid option: %s", so.Positional.Option)
		return nil
	}

	opt, _ := proto.Marshal(&clientpb.SetOptionReq{
		WorkspaceID: uint32(Context.Workspace.ID),
		User:        Context.Server.Config.User,
		Path:        Context.Module.Path,
		Name:        so.Positional.Option,
		Value:       so.Positional.Value,
	})

	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgOptionReq,
		Data: opt,
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError, "%s\n", resp.Err)
		return nil
	}

	changed := Context.Module.Options[so.Positional.Option]
	changed.Value = so.Positional.Value

	fmt.Printf("%s*%s %s => %s \n",
		tui.BLUE, tui.RESET, so.Positional.Option, so.Positional.Value)

	return nil
}
