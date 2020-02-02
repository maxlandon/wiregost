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

	// MODULES ----------------------------------------------------------------//

	// MsgModuleReq loads a module (server-side) so client can use it
	MsgModuleReq
	// MsgModule sends back necessary module information to the client
	MsgModule

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
