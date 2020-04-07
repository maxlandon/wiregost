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
	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
)

// Back from the current module but keep it on the stack
type ModuleBackCmd struct{}

var ModuleBack ModuleBackCmd

func RegisterModuleBack() {
	CommandParser.AddCommand(constants.ModuleBack, "", "", &ModuleBack)

	mb := CommandParser.Find(constants.ModuleBack)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], mb)
	mb.ShortDescription = "Back from the current module but keep it on the stack"
}

func (mb *ModuleBackCmd) Execute(args []string) error {

	// We just need to make a new, empty module.
	// Everything, from context to prompts, automatically adapt to this empty/non-empty module
	*Context.Module = clientpb.Module{}

	return nil
}
