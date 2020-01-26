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

package help

import (
	consts "github.com/maxlandon/wiregost/client/constants"
)

var (
	cmdHelp = map[string]string{
		// Data Service --------------------//

		// Workspace
		consts.WorkspaceStr: workspaceHelp,

		// Hosts
		consts.HostsStr:    hostsHelp,
		consts.HostsAdd:    hostsAdd,
		consts.HostsDelete: hostsDelete,
		consts.HostsUpdate: hostsUpdate,

		// Services
		consts.ServicesStr:    serviceHelp,
		consts.ServicesAdd:    servicesAdd,
		consts.ServicesDelete: servicesDelete,
		consts.ServicesUpdate: servicesUpdate,
	}
)

// GetHelpFor - Get help string for a command
func GetHelpFor(cmdName string) string {
	if 0 < len(cmdName) {
		if helpTempl, ok := cmdHelp[cmdName]; ok {
			return helpTempl
		}
	}
	return ""
}
