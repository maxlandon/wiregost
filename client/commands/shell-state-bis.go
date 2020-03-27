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

	"github.com/lmorg/readline"
	"github.com/maxlandon/wiregost/client/config"
	"github.com/maxlandon/wiregost/client/core"
	"github.com/maxlandon/wiregost/data-service/models"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
)

type ShellContext struct {
	Shell     *readline.Instance   // Shell object
	Config    *config.Config       // Shell configuration
	DBContext context.Context      // DB queries context
	Menu      *string              // Current shell menu
	Module    *clientpb.Module     // Current module
	UserID    int32                // Unique user ID for module requests
	Workspace *models.Workspace    // Current workspace
	Server    *core.WiregostServer // Wiregost Server
	Jobs      *int                 // Number of jobs
	Ghosts    *int                 // Number of connected implants
	Ghost     *clientpb.Ghost      // Current implant
	GhostPwd  *string              // Current working
}
