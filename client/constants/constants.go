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

	// Data Service --------------------------

	// Workspace
	WorkspaceStr = "workspace"

	// Hosts
	HostsStr    = "hosts"
	HostsAdd    = "hosts add"
	HostsDelete = "hosts delete"
	HostsUpdate = "hosts update"

	// Services
	ServicesStr    = "services"
	ServicesAdd    = "services add"
	ServicesDelete = "services delete"
	ServicesUpdate = "services update"

	// Groups
	DataServiceHelpGroup = "Data Service:"
)
