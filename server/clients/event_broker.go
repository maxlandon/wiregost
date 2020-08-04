package clients

import (
	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/events"
)

// eventBroker - A broker attached to a console, for pushing events.
type eventBroker struct {
}

// Events - This makes the broker act as a GRPC service, which pushes events to a console
func (b *eventBroker) Events(req *clientpb.Empty, stream serverpb.EventRPC_EventsServer) error {

	// Subscribe to event broker in events package
	incoming := events.Broker.Subscribe()
	defer events.Broker.Unsubscribe(incoming)

	// For each event coming in, check event type,
	for event := range incoming {

		// Depending on event type, we might have to push to several users/clients
		switch event.Type {
		case serverpb.EventType_USER:
			// We push anyway

		case serverpb.EventType_MODULE:
			// For modules, we push only if consoles matches the token

		case serverpb.EventType_SESSION:
			// For sessions, we push anyway.

		case serverpb.EventType_LISTENER:
			// We push anyway.

		case serverpb.EventType_JOB:
			// We push anyway.

		case serverpb.EventType_STACK:
			// We push only if user matches AND console DOES NOT match

		case serverpb.EventType_CANARY:
			// We push anyway

		}
	}

	return nil
}
