package logging

import (
	// Standard
	"fmt"

	// 3rd party
	"github.com/evilsocket/islazy/tui"
	"github.com/sirupsen/logrus"

	// Wiregost
	"github.com/maxlandon/wiregost/internal/messages"
)

var levels = map[string]int{
	"trace":   1,
	"debug":   2,
	"info":    3,
	"warning": 4,
	"error":   5,
	"fatal":   6,
	"panic":   7,
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

// SetLevel is used by clients to regulate the level of log events that are forwarded to them.
func (cl *ClientLogger) SetLevel(request messages.ClientRequest) {
	cl.Level = levels[request.Command[2]]
	// Return response
	status := fmt.Sprintf("[-] level => %s%s", tui.YELLOW, request.Command[2])
	res := messages.LogResponse{Log: status}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "log",
		Content:  res,
	}
	messages.Responses <- msg
}
