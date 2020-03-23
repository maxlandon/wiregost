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
	"fmt"

	"github.com/evilsocket/islazy/tui"
)

// RegisterCommands - register all commands in Wiregost and maps them
// to their respective contexts
func RegisterCommands() {

	// Core
	registerCoreCommands()
	registerHelpCommands()
	registerUserCommands()

	// Server
	registerServerCommands()

	// Data Service
	registerWorkspaceCommands()
	registerHostCommands()

	// Stack
	registerStackCommands()

	// Module
	registerModuleCommands()

	// Jobs
	registerJobCommands()

	// Profiles
	registerProfileCommands()

	// Builds & Canaries
	registerGhostBuildsCommands()

	// Sessions
	registerSessionCommands()

	// Ghosts
	registerAgentHelpCommands()
	registerFileSystemCommands()
	registerAgentInfoCommands()
	registerPrivCommands()
	registerProcCommands()
	registerExecuteCommands()

	// Nmap
	registerNmapCommands()
}

var (
	// Info - All normal message
	Info = fmt.Sprintf("%s[-]%s ", tui.BLUE, tui.RESET)
	// Warn - Errors in parameters, notifiable events in modules/sessions
	Warn = fmt.Sprintf("%s[!]%s ", tui.YELLOW, tui.RESET)
	// Error - Error in commands, filters, modules and implants.
	Error = fmt.Sprintf("%s[!]%s ", tui.RED, tui.RESET)
	// Success - Success events
	Success = fmt.Sprintf("%s[*]%s ", tui.GREEN, tui.RESET)

	// Infof - formatted
	Infof = fmt.Sprintf("%s[-] ", tui.BLUE)
	// Warnf - formatted
	Warnf = fmt.Sprintf("%s[!] ", tui.YELLOW)
	// Errorf - formatted
	Errorf = fmt.Sprintf("%s[!] ", tui.RED)
	// Sucessf - formatted
	Sucessf = fmt.Sprintf("%s[*] ", tui.GREEN)

	//RPCError - Errors from the server
	RPCError = fmt.Sprintf("%s[RPC Error]%s ", tui.RED, tui.RESET)
	// CommandError - Command input error
	CommandError = fmt.Sprintf("%s[Command Error]%s ", tui.RED, tui.RESET)
	// DBError - Data Service error
	DBError = fmt.Sprintf("%s[DB Error]%s ", tui.RED, tui.RESET)
)
