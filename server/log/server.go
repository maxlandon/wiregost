package log

import "github.com/sirupsen/logrus"

var (
	// ServLogger - Server logger
	ServLogger = serverLogger()

	// servTxtLogger - Used internally to change txt logging settings on the fly, such as file location
	servTxtLogger = serverTxtLogger()
)

// ServerLogger - Logs all events related to the server. This is the function actually called by other packages
func ServerLogger(logType string, pkg string, stream string) *logrus.Entry {

	// Set the correct file in which to write the event, depending on event type
	SetLogFiles(logType)

	// Return the logger once everything settings are correct
	return ServLogger.WithFields(logrus.Fields{
		"pkg":    pkg,
		"stream": stream,
	})
}

// SetLogFiles - Depending on the type of log event, set the apropriate output files (txt & json)
func SetLogFiles(logType string) {
	// TODO: Check if this does not induce bugs when multiple logs at the same time
	switch logType {
	case "listener":
		// ServLogger.Out = listenerLogFile
	case "compilation":
		// ServLogger.Out = compilationLogFile
	case "server":
		// ServLogger.Out = serverLogFile
	case "client":
		// ServLogger.Out = clientLogFile
	}
}

// serverLogger - All settings for the server JSON logger
func serverLogger() (logger *logrus.Logger) {
	return
}

// serverTxtLogger - Settings for the server Text logger
func serverTxtLogger() (logger *logrus.Logger) {
	return
}
