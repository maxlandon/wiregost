package log

import "github.com/sirupsen/logrus"

var (
	// GhostLogger - Logs all events sent back by ghost implants
	GhostLogger = ghostLogger()
)

// GhostLog - Logs all events related to the server. This is the function actually called by other packages
func GhostLog(ghostID uint32, name string, pkg string, stream string) *logrus.Entry {

	// Set apropriate output files, depending on ghost
	SetGhostLogFiles(ghostID, name)

	// Return the logger once everything settings are correct
	return GhostLogger.WithFields(logrus.Fields{
		"pkg":    pkg,
		"stream": stream,
	})
}

// SetGhostLogFiles - Sets the apropriate files for a given ghost logger.
func SetGhostLogFiles(ghostID uint32, name string) {

}

// ghostLogger - All settings for the server JSON logger
func ghostLogger() (logger *logrus.Logger) {
	return
}

// ghostTxtLogger - Settings for the server Text logger
func ghostTxtLogger() (logger *logrus.Logger) {
	return
}
