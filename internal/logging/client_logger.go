package logging

import (
	"fmt"

	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/sirupsen/logrus"
)

var levels = map[string]int{
	"debug":   1,
	"info":    2,
	"warning": 3,
	"error":   4,
	"fatal":   5,
	"panic":   6,
}

// ClientLogger is used to forward log events to its respective client.
// It filters what event to send based on its log level filter, which can be changed
// for and from each client shell.
type ClientLogger struct {
	// Client
	ClientID           int
	CurrentWorkspaceID *int
	// Level
	Level int
	// Channels with client
	ForwardClients chan<- messages.Message
}

// NewClientLogger instantiates a new Logger, dedicated to one client only.
func NewClientLogger(clientID int, workspaceID *int, channel chan<- messages.Message) *ClientLogger {
	logger := &ClientLogger{
		ClientID:           clientID,
		CurrentWorkspaceID: workspaceID,
		ForwardClients:     channel,
	}
	// Setup default level at startup
	logger.Level = levels["info"]

	return logger
}

// Forward log items to a given client
func (cl *ClientLogger) Forward(entry *logrus.Entry) error {

	if levels[entry.Level.String()] >= cl.Level {
		event := messages.LogEvent{
			ClientID:    cl.ClientID,
			WorkspaceID: *cl.CurrentWorkspaceID,
			Level:       entry.Level.String(),
			Message:     entry.Message,
		}
		msg := messages.Message{
			ClientID: cl.ClientID,
			Type:     "logEvent",
			Content:  event,
		}
		cl.ForwardClients <- msg
	}

	return nil
}

func (cl *ClientLogger) SetLevel(request messages.ClientRequest) {
	cl.Level = levels[request.Command[2]]
	// Return response
	status := fmt.Sprintf(" => %s", request.Command[2])
	res := messages.LogResponse{Log: status}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "log",
		Content:  res,
	}
	messages.Responses <- msg
}
