package events

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

// Server - A gRPC server that is accessed by networked components of Wiregost:
// user consoles, module stacks and maybe remote sessions/implants.
type Server struct {
	// There is one gRPC server instance all components.
	*pb.UnimplementedEventsServer
}

// EventsSubscribe - A networked component requested to subscribe to the Event broker.
func (s *Server) EventsSubscribe(in context.Context, req *pb.SubscribeRequest) (res *pb.Subscribe, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method EventsSubscribe not implemented")
}

// EventsUnsubscribe - Before disconnecting, a component requested to unsubscribe from the Event broker.
func (s *Server) EventsUnsubscribe(in context.Context, req *pb.UnsubscribeRequest) (res *pb.Unsubscribe, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method EventsUnsubscribe not implemented")
}

// Events - Each connected client will get a single instance of this function allocated to him.
func (s *Server) Events(in *clientpb.Client, stream pb.Events_EventsServer) (err error) {

	// Check user auth first with gRPC TLS creds (see Sliver)

	// Register local event channel
	events := Broker.Subscribe()
	defer Broker.Unsubscribe(events)

	// Process incoming events
	for event := range events {
		fmt.Println(event)

		// Initialize a nil Event. If at the goto Send, this var is not nil anymore,
		// it means we have to push the Event through the stream. Along this function,
		// the message can be
		var send *pb.Event

		// Many filterings have to take place first before pushing to client, in order:

		// If the Event Client/User is not ours, continue to next event
		if WrongUserClient(event, in) {
			continue
		}

		// Check if user still subscribes

		// If the message is destined to client, send it no matter the rest
		if ClientConcerned(event, in) {
			goto Send
		}

		if WrongClient(event, in) {
			continue
		}

		// If message went through here, this means the there is user specified
		// for the message, there is no client either, so that this component is
		// concerned.
		goto Send

	Send:
		if send == nil {
			err := stream.Send(&event)
			if err != nil {
				// log error to server here
			}
		}

	}

	return
}

// EventPush - A client pushes a new event to the server. Process it.
func (s *Server) EventPush(in context.Context, req *pb.Event) (res *clientpb.Empty, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method EventPush not implemented")
}

// ClientConcerned - Returns if the client is directly concerned by this event, no matter its type.
// This defacto filters events like module events (concerning only console).
func ClientConcerned(event pb.Event, client *clientpb.Client) (ok bool) {

	if event.Client == nil {
		return false
	}

	if event.Client.ID == client.ID {
		return true
	}

	if event.Client.User != nil && event.Client.User.ID == client.User.ID {
		return true
	}

	return false
}

// WrongClient - Returns if the user is concerned by this event, but the client is wrong.
// This filters events like Stack events (where all consoles of a user must adapt to).
func WrongClient(event pb.Event, client *clientpb.Client) (ok bool) {

	if event.Client.ID != client.ID {
		return true
	}

	return false
}

// WrongUserClient - Returns if the user is directly concerned by this event, no matter its type.
// This filters events like Stack events (where all consoles of a user must adapt to).
func WrongUserClient(event pb.Event, client *clientpb.Client) (ok bool) {

	if event.Client != nil {

		if event.Client.ID != client.ID {
			return true
		}

		if event.Client.User != nil && event.Client.User.ID != client.User.ID {
			return true
		}
	}

	return false
}
