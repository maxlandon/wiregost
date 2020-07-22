package events

import (
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

var (
	// EventBroker - Distributes event messages
	EventBroker = newBroker()
)

// broker - An object handling all events happening in Wiregost,
// and dispatching them to the appropriate client consoles.
type broker struct {
	stop        chan struct{}
	subscribe   chan chan serverpb.Event
	unsubscribe chan chan serverpb.Event
	publish     chan serverpb.Event
	send        chan serverpb.Event
}

// Push - This function is used by outside packages, in order to push events to consoles.
// These packages fill an Event object, with details, intended user and/or console.
// Then, this function determines which to client/user the event should be pushed.
func Push(event serverpb.Event) {
	EventBroker.publish <- event
}

func newBroker() *broker {
	b := &broker{
		stop:        make(chan struct{}),
		publish:     make(chan serverpb.Event, eventBufSize),
		subscribe:   make(chan chan serverpb.Event, eventBufSize),
		unsubscribe: make(chan chan serverpb.Event, eventBufSize),
		send:        make(chan serverpb.Event, eventBufSize),
	}
	go b.Start()
	return b
}

// Start processing event messages
func (b *broker) Start() {
	subscribers := map[chan serverpb.Event]struct{}{}
	for {
		select {
		case <-b.stop:
			for sub := range subscribers {
				close(sub)
			}
			return
		case sub := <-b.subscribe:
			subscribers[sub] = struct{}{}
		case sub := <-b.unsubscribe:
			delete(subscribers, sub)
		case event := <-b.publish:
			for sub := range subscribers {
				sub <- event
			}
		}
	}
}

// Subscribe - Generate a new subscription channel
func (b *broker) Subscribe() chan serverpb.Event {
	events := make(chan serverpb.Event, eventBufSize)
	b.subscribe <- events
	return events
}

// Unsubscribe - Remove a subscription channel
func (b *broker) Unsubscribe(events chan serverpb.Event) {
	b.unsubscribe <- events
	close(events)
}

const (
	// Size is arbitrary, just want to avoid weird cases where we'd block on channel sends
	eventBufSize = 5
)
