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

import (
	"context"
	"fmt"
	"io"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/connection"
	cliCtx "github.com/maxlandon/wiregost/client/context"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

// StartEventListener - Listens for events coming from the server/ghosts.
func (c *console) StartEventListener() {

	// Listen for RPC events
	events, err := connection.EventsRPC.Events(context.Background(), cliCtx.Context.Client)
	if err != nil {
		fmt.Println(RPCError + tui.Red("There was an error with the Event stream"))
	}

	for {
		event, err := events.Recv()
		if err == io.EOF || event == nil { // Safety checks
			return
		}
		// Switch event type
		switch event.Type {
		case serverpb.EventType_MODULE:
			ModuleEvent(event)
		}

	}
}

// ModuleEvent - Console behavior upon module event reception.
func ModuleEvent(event *serverpb.Event) {

	// Final status line printed to console
	line := Levels[event.Level]

	// Format the entry module type. We assume a max length of 10
	// We also print this module name in red if its an error message
	if event.Level != serverpb.Level_ERROR {
		line += fmt.Sprintf("%-10v", event.Module)
	} else {
		line += fmt.Sprintf("%s%-10v%s %s-%s ",
			tui.RED, event.Module, tui.RESET, tui.Dim, tui.RESET)
	}

	// This case should not happen normally, but handle it...
	// We add a second line with the message in addition, just in case
	if event.Err != "" && event.Message != "" {
		line += event.Err
		line += "\n"
		line += Levels[serverpb.Level_ERROR]
		line += fmt.Sprintf("%-10v %s-%s", " ", tui.DIM, tui.RESET)
		line += event.Message

		fmt.Println(line)
		return
	}

	// Else we deal with fields
	if event.Err != "" {
		line += event.Err
	}
	if event.Message != "" {
		line += event.Message
	}

	fmt.Println(line)
}

// ImplantEvent - Console behavior upon ghost implant event reception.
func ImplantEvent() {}

// CanaryEvent - Console behavior upon canary alert reception.
func CanaryEvent() {}

// UserEvent - Console behavior upon user event reception (connections, disconnections, etc)
func UserEvent() {}

// Levels - Maps Event levels with their associated string icon.
var Levels = map[serverpb.Level]string{
	serverpb.Level_TRACE:   "",
	serverpb.Level_DEBUG:   "",
	serverpb.Level_INFO:    "",
	serverpb.Level_WARNING: "",
	serverpb.Level_ERROR:   "",
}
