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

package clientpb

const (
	clientOffset = 10000 // To avoid duplications with ghostpb constants

	// MsgEvent - Initial message from sliver with metadata
	MsgEvent = uint32(clientOffset + iota)

	// MESSAGES ---------------------------------------------------------------//

	// MsgTcp - TCP message
	MsgTcp
	// MsgMtls - MTLS message
	MsgMtls
	// MsgDns - DNS message
	MsgDns

	// USERS ------------------------------------------------------------------//

	// MsgUserReq - Add a user and create its config
	MsgUserReq
	// MsgDeleteUserReq -  Delete a user from the server
	MsgDeleteUserReq
	// MsgUser - Success/failure of user creation request
	MsgUser

	// MODULES ----------------------------------------------------------------//

	// MsgStackReq loads a module (server-side) so client can use it
	MsgStackUse
	// MsgStackPop pops one or all modules from a stack
	MsgStackPop
	// MsgStackList request a list of modules loaded onto a stack
	MsgStackList

	// MsgModuleReq - Performs an action on a module
	MsgModuleReq
	// MsgModule - Success/failure of MsgModuleReq
	MsgModule
	// MsgModuleList - Get a list of all modules in Wiregost
	MsgModuleList

	// MsgOptionReq - Set a module option
	MsgOptionReq
	// MsgOption - Success/failure of option request
	MsgOption

	// NON-LISTENER JOBS ------------------------------------------------------//

	// MsgJobs - Jobs message
	MsgJobs
	MsgJobKill

	// LISTENERS -------------------------------------------------------------//

	// MsgHttp - HTTP(S) listener  message
	MsgHttp
	// MsgHttps - HTTP(S) listener  message
	MsgHttps

	// MsgProxy - Start a server-side Proxy for routing traffic
	MsgProxy

	// COMMS -----------------------------------------------------------------//

	// MsgTunnelCreate - Create tunnel message
	MsgTunnelCreate
	MsgTunnelClose

	// COMPILER -------------------------------------------------------------//

	// MsgGenerate - Generate message
	MsgGenerate
	MsgNewProfile
	MsgProfiles
	MsgDeleteProfile
	MsgTask
	MsgMigrate
	MsgGetSystemReq
	MsgEggReq
	MsgExecuteAssemblyReq
	MsgSideloadReq

	MsgRegenerate

	MsgListGhostBuilds
	MsgListCanaries

	// GHOSTS ---------------------------------------------------------------//

	// MsgSessions - Sessions message
	MsgSessions

	// METASPLOIT ----------------------------------------------------------//

	// MsgMsf - MSF message
	MsgMsf
	// MsgMsfInject - MSF injection message
	MsgMsfInject

	// WEBSITES -----------------------------------------------------------//
	// Website related messages
	MsgWebsiteList
	MsgWebsiteAddContent
	MsgWebsiteRemoveContent
)
