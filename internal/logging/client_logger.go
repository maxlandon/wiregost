package logging

import (
	"fmt"

	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/sirupsen/logrus"
)

var Levels = map[string]int{
	"debug":   1,
	"info":    2,
	"warning": 3,
	"error":   4,
	"fatal":   5,
	"panic":   6,
}

type ClientLogger struct {
	// Client
	ClientId           int
	CurrentWorkspaceId *int
	// Level
	Level int
	// Channels with client
	ForwardClients chan<- messages.Message
}

func NewClientLogger(clientId int, workspaceId *int, channel chan<- messages.Message) *ClientLogger {
	logger := &ClientLogger{
		ClientId:           clientId,
		CurrentWorkspaceId: workspaceId,
		ForwardClients:     channel,
	}
	// Setup default level at startup
	logger.Level = Levels["info"]

	return logger
}

// Forward log items to a given client
func (cl *ClientLogger) Forward(entry *logrus.Entry) error {

	if Levels[entry.Level.String()] >= cl.Level {
		event := messages.LogEvent{
			ClientId:    cl.ClientId,
			WorkspaceId: *cl.CurrentWorkspaceId,
			Level:       entry.Level.String(),
			Message:     entry.Message,
		}
		msg := messages.Message{
			ClientId: cl.ClientId,
			Type:     "logEvent",
			Content:  event,
		}
		cl.ForwardClients <- msg
	}

	return nil
}

func (cl *ClientLogger) SetLevel(request messages.ClientRequest) {
	cl.Level = Levels[request.Command[2]]
	// Return response
	status := fmt.Sprintf(" => %s", request.Command[2])
	res := messages.LogResponse{Log: status}
	msg := messages.Message{
		ClientId: request.ClientId,
		Type:     "log",
		Content:  res,
	}
	messages.Responses <- msg
}
