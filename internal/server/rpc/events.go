package rpc

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

import (
	"github.com/maxlandon/wiregost/internal/proto/clientpb"
	"github.com/maxlandon/wiregost/internal/proto/commonpb"
	"github.com/maxlandon/wiregost/internal/proto/rpcpb"
	"github.com/maxlandon/wiregost/internal/server/core"
	"github.com/maxlandon/wiregost/internal/server/log"
)

var rpcEventsLog = log.NamedLogger("rpc", "events")

// Events - Stream events to client
func (rpc *Server) Events(_ *commonpb.Empty, stream rpcpb.Core_EventsServer) error {
	commonName := rpc.getClientCommonName(stream.Context())
	client := core.NewClient(commonName)
	core.Clients.Add(client)
	events := core.EventBroker.Subscribe()

	defer func() {
		rpcEventsLog.Infof("%d client disconnected", client.ID)
		core.EventBroker.Unsubscribe(events)
		core.Clients.Remove(client.ID)
	}()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case event := <-events:
			pbEvent := &clientpb.Event{
				EventType: event.EventType,
				Data:      event.Data,
			}

			if event.Job != nil {
				pbEvent.Job = event.Job.ToProtobuf()
			}
			if event.Client != nil {
				pbEvent.Client = event.Client.ToProtobuf()
			}
			if event.Err != nil {
				pbEvent.Err = event.Err.Error()
			}

			err := stream.Send(pbEvent)
			if err != nil {
				rpcEventsLog.Warnf(err.Error())
				return err
			}
		}
	}
}
