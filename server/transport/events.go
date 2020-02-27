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

package transport

import (
	"fmt"
	"net"

	"github.com/gogo/protobuf/proto"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
)

func socketEventLoop(conn net.Conn, events chan core.Event) {
	for event := range events {
		pbEvent := &clientpb.Event{
			EventType:       event.EventType,
			EventSubType:    event.EventSubType,
			ModuleRequestID: event.ModuleRequestID,
			Data:            event.Data,
		}

		if event.Job != nil {
			pbEvent.Job = event.Job.ToProtobuf()
		}
		if event.Client != nil {
			pbEvent.Client = event.Client.ToProtobuf()
		}
		if event.Ghost != nil {
			pbEvent.Ghost = event.Ghost.ToProtobuf()
		}
		if event.Err != nil {
			pbEvent.Err = fmt.Sprintf("%v", event.Err)
		}

		data, _ := proto.Marshal(pbEvent)
		envelope := &ghostpb.Envelope{
			Type: clientpb.MsgEvent,
			Data: data,
		}
		err := socketWriteEnvelope(conn, envelope)
		if err != nil {
			clientLog.Errorf("Socket write failed %v", err)
			return
		}
	}
}
