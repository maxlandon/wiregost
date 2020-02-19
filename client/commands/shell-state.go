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

package commands

import (
	"context"

	ccore "github.com/maxlandon/wiregost/client/core"
	"github.com/maxlandon/wiregost/data_service/models"
	score "github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/module/templates"
)

// ShellContext - Passes client shell variable pointers to command for read/write access
type ShellContext struct {
	// Context
	Context     context.Context
	MenuContext *string
	Mode        *string

	CurrentModule    *string
	Module           *templates.Module
	CurrentWorkspace *models.Workspace

	// Server state
	Server *ccore.WiregostServer

	// Jobs
	Listeners *int

	// Agents
	Ghosts *int
	// Keep for prompt, until not needed anymore
	CurrentAgent *score.Ghost
}
