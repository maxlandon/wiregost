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

package constants

// [ Meta ] -----------------------------------------------------------------//
const (
	KeepAliveStr = "keepalive"
)

// [ Events ] ---------------------------------------------------------------//
const (
	EventStr = "event"

	ServerErrorStr = "server"

	StackEvent  = "stack pop"
	ModuleEvent = "module"

	// ConnectedEvent - Ghost Connected
	ConnectedEvent = "connected"
	// DisconnectedEvent - Ghost disconnected
	DisconnectedEvent = "disconnected"

	// JoinedEvent - User joined server
	JoinedEvent = "joined"
	// LeftEvent - Player disconnected from server
	LeftEvent = "left"

	// CanaryEvent - A DNS canary was triggered
	CanaryEvent = "canary"

	// StartedEvent - Job was started
	StartedEvent = "started"
	// StoppedEvent - Job was stopped
	StoppedEvent = "stopped"
)

// [ Commands ] -----------------------------------------------------------//
const (
	// Core ----------------------------------

	Help     = "help"
	Core     = "core"
	Shell    = "!"
	Cd       = "cd"
	Resource = "resource"

	// Server
	Server        = "server"
	ServerConnect = "connect"

	// User
	User = "user"

	// Data Service --------------------------

	// Workspace
	Workspace       = "workspace"
	WorkspaceSwitch = "switch"
	WorkspaceAdd    = "add"
	WorkspaceDelete = "delete"
	WorkspaceUpdate = "update"

	// Hosts
	Hosts       = "hosts"
	HostsAdd    = "add"
	HostsDelete = "delete"
	HostsUpdate = "update"

	// Services
	Services       = "services"
	ServicesAdd    = "add"
	ServicesDelete = "delete"
	ServicesUpdate = "update"

	// Groups
	DataServiceHelpGroup = "Data Service:"

	// Stack & Modules -----------------------
	Stack              = "stack"
	StackList          = "list"
	StackUse           = "use"
	StackPop           = "pop"
	Module             = "module"
	ModuleUse          = "use"
	ModuleInfo         = "info"
	ModuleOptions      = "options"
	ModuleSetOption    = "set"
	ModuleList         = "list"
	ModuleRun          = "run"
	ModuleToListener   = "to-listener"
	ModuleParseProfile = "parse-profile"
	ModuleToProfile    = "to-profile"
	ModuleBack         = "back"

	// Jobs
	Jobs        = "jobs"
	JobsKill    = "kill"
	JobsKillAll = "kill-all"

	// Profiles
	Profiles       = "profiles"
	ProfilesDelete = "delete"

	// GhostBuilds
	Ghosts = "ghosts"

	// Canaries
	Canaries = "canaries"

	// Sessions
	Sessions           = "sessions"
	SessionsInteract   = "interact"
	SessionsKill       = "kill"
	SessionsKillAll    = "kill-all"
	SessionsBackground = "background"

	// Ghost Implants -------------------------
	// FileSystem
	GhostCd       = "cd"
	GhostLs       = "ls"
	GhostPwd      = "pwd"
	GhostCat      = "cat"
	GhostMkdir    = "mkdir"
	GhostRm       = "rm"
	GhostDownload = "download"
	GhostUpload   = "upload"

	// Info
	GetUID = "getuid"
	GetPID = "getpid"
	GetGID = "getgid"
	Whoami = "whoami"

	// Network
	IfConfig = "ifconfig"
	Netstat  = "netstat"

	// Proc
	Ps        = "ps"
	ProcDump  = "procdump"
	Terminate = "terminate"
	Migrate   = "migrate"

	// Priv
	RunAs       = "run-as"
	Impersonate = "impersonate"
	Elevate     = "elevate"
	GetSystem   = "getsystem"
	Rev2Self    = "rev-to-self"

	// Execute
	Execute          = "execute"
	ExecuteShellcode = "execute-shellcode"
	ExecuteAssembly  = "execute-assembly"
	Sideload         = "sideload"
	SpawnDll         = "spawn-dll"
	MsfInject        = "msf-inject"

	// Shell
	SystemShell = "shell"

	// Help
	AgentHelp         = "agent-help"
	CompleteAgentHelp = "complete-agent-help"
)
