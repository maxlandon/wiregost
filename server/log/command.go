package log

import "github.com/sirupsen/logrus"

// CommandLogger - Logs all commands typed by a user. This is the function actually called by other packages.
// filepath   -> GetUserConsoleHistory(user) should yield a path to pass as argument to this function.
func CommandLogger(user string, clientID string, filepath string) (logg *logrus.Entry) {

	// A logger with default settings
	var logger = defLogger()
	logger.Formatter = &logrus.TextFormatter{}

	// Return the logger once everything settings are correct
	return logger.WithFields(logrus.Fields{
		"user":     user,
		"clientID": clientID,
	})
}

// defaultLogger - All logger in Wiregost will use an instance of a
// default logger, that they will modify depending on their needs.
func defLogger() (logger *logrus.Logger) {
	logger = logrus.New()
	return
}

// commandLogger - All settings for the command JSON logger
func commandLogger() (logger *logrus.Logger) {

	// Format settings (dates, formats, etc)
	logger = logrus.New()
	logger.Formatter = &logrus.TextFormatter{}
	return
}
