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

package rpc

import (
	"time"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/log"
)

var (
	rpcLog = log.ServerLogger("rpc", "server")
)

// Response - Called with response data, mapped back to ReqID
type Response func([]byte, error)

// Handler - RPC handlers accept bytes and return bytes
type Handler func([]byte, time.Duration, Response)

// TunnelHandler - Tunnel handlers join tunnels from client and server
type TunnelHandler func(*core.Client, []byte, Response)

var (
	rpcHandlers = &map[uint32]Handler{

		// CLIENT -------------------------------------//
		// Users
		clientpb.MsgUser:          rpcListUsers,
		clientpb.MsgUserReq:       rpcAddUser,
		clientpb.MsgDeleteUserReq: rpcDeleteUser,

		// Stack
		clientpb.MsgStackUse:  rpcStackUse,
		clientpb.MsgStackPop:  rpcStackPop,
		clientpb.MsgStackList: rpcStackList,

		// Module
		clientpb.MsgOptionReq:  rpcModuleSetOption,
		clientpb.MsgModuleReq:  rpcModuleRun,
		clientpb.MsgModuleList: rpcModuleList,

		// Jobs
		clientpb.MsgJobs:    rpcJobs,
		clientpb.MsgJobKill: rpcJobKill,

		// Profiles
		clientpb.MsgProfiles:      rpcListProfiles,
		clientpb.MsgDeleteProfile: rpcDeleteProfile,

		// Builds & Canaries
		clientpb.MsgListGhostBuilds: rpcGhostBuilds,
		clientpb.MsgListCanaries:    rpcListCanaries,

		// Sessions
		clientpb.MsgSessions: rpcSessions,

		// GHOST -------------------------------------//
		// Special
		ghostpb.MsgKill: rpcKill,

		// File system
		ghostpb.MsgCdReq:       rpcCd,
		ghostpb.MsgPwdReq:      rpcPwd,
		ghostpb.MsgLsReq:       rpcLs,
		ghostpb.MsgRmReq:       rpcRm,
		ghostpb.MsgMkdirReq:    rpcMkdir,
		ghostpb.MsgDownloadReq: rpcDownload,
		ghostpb.MsgUploadReq:   rpcUpload,

		// Info
		ghostpb.MsgIfconfigReq: rpcIfconfig,
		ghostpb.MsgNetstatReq:  rpcNetstat,

		// Priv
		ghostpb.MsgElevate:       rpcElevate,
		ghostpb.MsgImpersonate:   rpcImpersonate,
		ghostpb.MsgRunAs:         rpcRunAs,
		ghostpb.MsgRevToSelf:     rpcRevToSelf,
		clientpb.MsgGetSystemReq: rpcGetSystem,

		// Proc
		ghostpb.MsgPsReq:          rpcPs,
		ghostpb.MsgTerminate:      rpcTerminate,
		ghostpb.MsgProcessDumpReq: rpcProcdump,
		clientpb.MsgMigrate:       rpcMigrate,

		// Shell
		ghostpb.MsgShellReq: rpcShell,

		// Execute
		clientpb.MsgTask:               rpcTask,
		ghostpb.MsgExecuteReq:          rpcExecute,
		clientpb.MsgExecuteAssemblyReq: rpcExecuteAssembly,
		clientpb.MsgMsfInject:          rpcMsfInject,
	}

	tunnelHandlers = &map[uint32]TunnelHandler{
		clientpb.MsgTunnelCreate: tunnelCreate,
		ghostpb.MsgTunnelData:    tunnelData,
		ghostpb.MsgTunnelClose:   tunnelClose,
	}
)

// GetRPCHandlers - Returns a map of server-side msg handlers
func GetRPCHandlers() *map[uint32]Handler {
	return rpcHandlers
}

// GetTunnelHandlers - Returns a map of tunnel handlers
func GetTunnelHandlers() *map[uint32]TunnelHandler {
	return tunnelHandlers
}
