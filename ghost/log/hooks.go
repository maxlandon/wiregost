package log

import (
	"github.com/maxlandon/wiregost/ghost/assets"
	"github.com/sirupsen/logrus"
)

// RemoteLogger - A hook for logging with text formatting
type RemoteLogger struct {
	Name   string
	logger *logrus.Logger
}

// NewRemoteLogger - New hook
func NewRemoteLogger(name string, logger *logrus.Logger) (hook *RemoteLogger) {
	return
}

// Fire - Function needed to implement the logrus.TxtLogger interface
func (hook *RemoteLogger) Fire(entry *logrus.Entry) (err error) {

	// If DebugRemote, send log to appropriate function, which
	// will handle details for sending logs back to server.
	// (Timing of reports and route strategies)
	if assets.DebugRemote == "true" {
		return HandleRemoteLog(entry)
	}

	return
}

// Levels - Function needed to implement the logrus.TxtLogger interface
func (hook *RemoteLogger) Levels() (levels []logrus.Level) {
	return logrus.AllLevels
}
