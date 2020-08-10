package log

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/events"
)

// ModuleLogger - The module.Module base type embeds a Logger, for giving logging capabilities for all module
// subtypes. This means any base Module type, whether it be an Exploit, a Transport, a Payload, etc. will derive
// its logging from the same Logger instance, merely adding fields to it.
// Potentially, we could analyse fields passed in those loggers for varying the output, from files to consoles.
// Each time a Module logger is instantiated, it should be passed a full Client object, containing its user.
// This client/user information is propagated down all module subtypes.
func ModuleLogger(remote bool, cli *clientpb.Client, rpc serverpb.EventsClient) (entry *logrus.Entry) {

	// New logger instance
	logger := logrus.New()

	// Formatting
	logger.Formatter = &logrus.TextFormatter{}
	logger.Out = ioutil.Discard // Don't need to output anything to stdout

	// Register to event system
	logger.AddHook(newModEvent(cli, rpc, remote))

	// The "module" key:value will be changed by module subtypes.
	return logger.WithField("module", "module")
}

// modEvent - A hook for logging with text formatting
type modEvent struct {
	client   *clientpb.Client
	rpc      serverpb.EventsClient
	isRemote bool
}

// newModEvent - Adds the capability to push log events back to the console, while retaining logrus capabilities.
func newModEvent(client *clientpb.Client, rpc serverpb.EventsClient, remote bool) (hook *modEvent) {
	return &modEvent{
		client:   client,
		rpc:      rpc,    // There is no gRPC client connected yet.
		isRemote: remote, // Therefore, by default the module is local.
	}
}

// Fire - Function needed to implement the logrus.TxtLogger interface
func (hook *modEvent) Fire(entry *logrus.Entry) (err error) {

	// Source file fields
	file := "<no caller>"
	if entry.HasCaller() {
		idx := strings.Index(entry.Caller.File, "wiregost")
		file = entry.Caller.File
		if idx != -1 {
			file = file[idx:]
		}
	}

	// Log this to source files, with more fields if needed
	// Many things to do here, because module content needs variable logging:
	// It might have to be blended in other log files, such as users' ones, sessions', etc.

	// Create event
	var event = serverpb.Event{}
	event.Client = hook.client // First things first, add Client and User (should be in Client)
	event.Type = serverpb.EventType_MODULE
	event.Level = EventLevels[entry.Level]
	event.Message = entry.Message
	event.Module = entry.Data["module"].(string) // Passed when a module has been instantiated.

	// Depending on the remote status of the logger, push it remote or locally
	if hook.isRemote {
		// Push event through gRPC client. It has been wired by the stack binary
		hook.rpc.EventPush(context.Background(), &event)
		return
	}

	// Or push event locally
	events.Push(event)

	return

}

// Levels - Function needed to implement the logrus.TxtLogger interface
func (hook *modEvent) Levels() (levels []logrus.Level) {
	return logrus.AllLevels
}

// EventLevels - A map of levels (as passed by the Event system) to logrus.Levels.
var EventLevels = map[logrus.Level]serverpb.Level{
	logrus.TraceLevel: 0,
	logrus.DebugLevel: 1,
	logrus.InfoLevel:  2,
	logrus.WarnLevel:  3,
	logrus.ErrorLevel: 4,
}
