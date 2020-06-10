package log

import "github.com/sirupsen/logrus"

var (
	// UsrLogger - Logs all events pertaining to a user (commands to implant, modules run, etc)
	UsrLogger = userLogger()
)

// UserLogger - Logger for a user's action. This is the function actually called by other packages
func UserLogger(userName string, userID uint32, pkg string, stream string) (entry *logrus.Entry) {
	// Set apropriate output files, depending on user name/ID/both
	SetUserLogFiles(userName, userID)

	// Return the entry once all settings are correct
	return
}

// userLogger - Logs all events related to a user
func userLogger() (logger *logrus.Logger) {
	return
}

// SetUserLogFiles - Sets the apropriate log file for the logger.
func SetUserLogFiles(userName string, userID uint32) {

}
