package log

import (
	"io/ioutil"

	"github.com/maxlandon/wiregost/ghost/assets"
	"github.com/sirupsen/logrus"
)

var (
	// GhostLogger - Logs all events of a ghost implant process
	GhostLogger = ghostLogger()
)

// GhostLog - Logs all events related to the server. This is the function actually called by other packages
func GhostLog(pkg string, stream string) *logrus.Entry {

	// Return the logger once everything settings are correct
	return GhostLogger.WithFields(logrus.Fields{
		"pkg":    pkg,
		"stream": stream,
	})
}

// ghostLogger - All settings for the server JSON logger
func ghostLogger() (logger *logrus.Logger) {

	// Format settings (dates, formats, etc)
	logger = logrus.New()
	logger.Formatter = &logrus.TextFormatter{}

	// Desactivate local print
	// Change output sources, make them nil by default:
	// We must be sure that no log is output to a source
	// we do not control, or that we did not explictly allowed.
	if assets.DebugLocal != "true" {
		logger.Out = ioutil.Discard
	}

	// Add local and remote hooks
	logger.AddHook(NewRemoteLogger("ghost", logger))

	return
}

// ghostTxtLogger - Settings for the server Text logger
func ghostTxtLogger() (logger *logrus.Logger) {
	return
}

// SetupLogging - Inits all logging infrastructure for implant
func SetupLogging() {

	// We make a temporary log for logging ourselves
	tmpLog := GhostLog("log", "setup")
	tmpLog.Info("Initialized logging infrastructure")
}
