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

package console

// StartEventListener - Listens for events coming from the server/ghosts.
func (c *console) StartEventListener() {

	// Listen for RPC events

	for {
		// Switch event type

		// Call function for this event
	}
}

// ModuleEvent - Console behavior upon module event reception.
func ModuleEvent() {
	// If pending

	// If non-pending
}

// ImplantEvent - Console behavior upon ghost implant event reception.
func ImplantEvent() {}

// CanaryEvent - Console behavior upon canary alert reception.
func CanaryEvent() {}

// UserEvent - Console behavior upon user event reception (connections, disconnections, etc)
func UserEvent() {}
