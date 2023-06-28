package constants

// Wiregost - Post-Exploitation & Implant Framework
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

// Events
const (
	// UpdateStr - "update"
	UpdateStr = "update"
	// VersionStr - "version"
	VersionStr = "version"

	// EventStr - "event"
	EventStr = "event"

	// ServersStr - "server-error"
	ServerErrorStr = "server-error"

	// JoinedEvent - Player joined the game
	JoinedEvent = "client-joined"
	// LeftEvent - Player left the game
	LeftEvent = "client-left"

	// StartedEvent - Job was started
	JobStartedEvent = "job-started"
	// StoppedEvent - Job was stopped
	JobStoppedEvent = "job-stopped"
)
