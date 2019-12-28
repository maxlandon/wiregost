package logging

import (
	"time"

	"github.com/go-playground/log"
)

// All framework components need to have their own Message() function:
// This function routes all their logs to the HandleLogs() function here,
// via a channel.
// func Message(lvl log.Level, msg string) {
//      fields := specific to the component/workspace/client/module
//      log := log.Entry{ Message: msg, Component: module, LogLevel: lvl}
//      logging.LogReqs <- log
//}

// Example: Message(logging.debug, "Received connection request")

// Log Levels
var (
	debug  = log.DebugLevel
	info   = log.InfoLevel
	notice = log.NoticeLevel
	warn   = log.WarnLevel
	alert  = log.AlertLevel
	err    = log.ErrorLevel
	pan    = log.PanicLevel
)

// Channels. All channels convey log entries, with all information needed
var moduleLogs = make(chan log.Entry)
var workspaceLogs = make(chan log.Entry)
var endpointLogs = make(chan log.Entry)
var LogReqs = make(chan log.Entry)

// SET TIME FORMAT
// cLog.SetTimestampFormat("01-02 15:04:05")

// LogManager is a handler taking care of :
//    - dispatching logs to their respective client loggers     -> HandleLogs()
//    - saving all framework logs to their respective files.    -> SaveLogs()
// It's always available at debug level, so that everything is logged.
type LogManager struct {
}

// Each client has its own dedicated ClientLogger instance, because it can then
// modulate its log level independently of others.
type ClientLogger struct {
	ClientId           int
	CurrentWorkspaceId int

	// ClientLogger forwards logs to their appropriate destinations, and keeps track of
	// which clients/workspaces ask for which log levels.
}

// All logs in WireGost are centralized and processed by this function.
// They are then dispatched to their respective client loggers.
func (lm *LogManager) HandleLogs() {
	// Loop indefinitely to receive logs
	for {
		req := <-LogReqs
		entry := log.Entry{
			Timestamp: time.Now(),
			Message:   req.Message,
			Fields:    req.Fields,
			Level:     req.Level,
		}
		// Send request to client loggers and file logger.
		log.HandleEntry(entry)
	}
}

// All logs from all components will be handled by this function.
// The function filters the ones pertaining to the general framework,
// and saves them to a log file.
func (lm *LogManager) Log(e log.Entry) {

}

// Log() method needs to be implemented for satisfying the Logger interface
func (cl *ClientLogger) Log(e log.Entry) {
	// Once the entry has been received by all client loggers satisfying the debug level
}
