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

package ghost

// Message Name Constants

const (
	// BASE CONTROL -----------------------------------------------------------------------//

	// MsgRegister - Initial message from ghost with metadata
	MsgRegister = uint32(1 + iota)

	// MsgPing - Confirm connection is open used as req/resp
	MsgPing

	// MsgKill - Kill request to the ghost process
	MsgKill

	// BASE COMMANDS -----------------------------------------------------------------------//

	// MsgLsReq - Request a directory listing from the remote system
	MsgLsReq
	//MsgLs - Directory listing (resp to MsgDirListReq)
	MsgLs

	// MsgDownloadReq - Request to download a file from the remote system
	MsgDownloadReq
	// MsgDownload - File contents for download (resp to DownloadReq)
	MsgDownload

	// MsgUploadReq - Upload a file to the remote file system
	MsgUploadReq
	// MsgUpload - Confirms success/failure of the file upload
	MsgUpload

	// MsgCdReq - Request a change directory on the remote system
	MsgCdReq
	// MsgCd - Confirms the success/failure of the `cd` request (resp to MsgCdReq)
	MsgCd

	// MsgPwdReq - A request to get the CWD from the remote process
	MsgPwdReq
	// MsgPwd - The CWD of the remote process (resp to MsgPwdReq)
	MsgPwd

	// MsgRmReq - Request to delete remote file
	MsgRmReq
	// MsgRm - Confirms the success/failure of delete request (resp to MsgRmReq)
	MsgRm

	// MsgMkdirReq - Request to create a directory on the remote system
	MsgMkdirReq
	// MsgMkdir - Confirms the success/failure of the mkdir request (resp to MsgMkdirReq)
	MsgMkdir

	// MsgPsReq - List processes req
	MsgPsReq
	// MsgPs - List processes resp
	MsgPs

	// MsgShellReq - Starts an interactive shell
	MsgShellReq
	// MsgShell - Response on starting shell
	MsgShell

	// MODULES ----------------------------------------------------------------------------//

	// MsgExecuteReq - Execute a command on the remote system, like Merlin's Modules
	MsgExecuteReq

	// NETWORK ----------------------------------------------------------------------------//

	// MsgTunnelData - Data for duplex tunnels
	MsgTunnelData
	// MsgTunnelClose - Close a duplex tunnel
	MsgTunnelClose

	// Routing ----------------------------------------------------------------------------//

	// MsgRouteReq - Create a network route on an implant
	MsgRouteReq
	// MsgRoute - Success/failure of route creation
	MsgRoute

	// MsgPivotReq - Create/Modify a pivot (a route with forwarding rules)
	MsgPivotReq
	// MsgPivot - Success/failure of pivot creation
	MsgPivot

	// PROCESS ----------------------------------------------------------------------------//

	// MsgTerminate - Kill a remote process
	MsgTerminate
	// MsgProcessDumpReq - Request to create a process dump
	MsgProcessDumpReq
	// MsgProcessDump - Dump of process)
	MsgProcessDump
	// MsgImpersonateReq - Request for process impersonation
	MsgImpersonateReq
	// MsgImpersonate - Output of the impersonation command
	MsgImpersonate
	// MsgRunAs - Run process as user
	MsgRunAs
	// MsgRevToSelf - Revert to self
	MsgRevToSelf
	// MsgGetSystemReq - Elevate as SYSTEM user
	MsgGetSystemReq
	// MsgGetSystem - Response to getsystem request
	MsgGetSystem
	// MsgElevateReq - Request to run a new sliver session in an elevated context
	MsgElevateReq
	//MsgElevate - Response to the elevation request
	MsgElevate
	// MsgMigrateReq - Spawn a new sliver in a designated process
	MsgMigrateReq

	// INJECTION ---------------------------------------------------------------------------//

	// MsgTask - A local shellcode injection task
	MsgTask

	// MsgRemoteTask - Remote thread injection task
	MsgRemoteTask

	// MsgExecuteAssemblyReq - Request to load and execute a .NET assembly
	MsgExecuteAssemblyReq
	// MsgExecuteAssembly - Output of the assembly execution
	MsgExecuteAssembly
	// MsgSideloadReq - request to sideload a binary
	MsgSideloadReq
	// MsgSideload - output of the binary
	MsgSideload
	// MsgSpawnDllReq - Reflective DLL injection request
	MsgSpawnDllReq
	// MsgSpawnDll - Reflective DLL injection output
	MsgSpawnDll
	// MsgIfconfigReq - Ifconfig (network interface config) request
	MsgIfconfigReq
)
