package log

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
)

var (
	// ModLogger - Logs all events pertaining to a module (commands to implant, modules run, etc)
	ModLogger = &logrus.Logger{}
)

// ModuleLogger - Logger for a module's action. This is the function actually called by other packages
func ModuleLogger(modulePath string, cli *clientpb.Client) (entry *logrus.Entry) {

	// Set apropriate output files, depending on module name/ID/both
	// SetModuleLogFiles(moduleName, moduleID)

	// Return the entry once all settings are correct
	return
}

// moduleLogger - Logs all events related to a module
func moduleLogger() (logger *logrus.Logger) {
	// Format settings (dates, formats, etc)
	logger = logrus.New()
	logger.Formatter = &logrus.TextFormatter{}
	logger.Out = ioutil.Discard

	// logger.AddHook(NewModEvent())
	return
}

// SetModuleLogFiles - Sets the apropriate log file for the logger.
func SetModuleLogFiles(moduleName string, moduleID uint32) {

}

// modEvent - A hook for logging with text formatting
type modEvent struct {
	client *clientpb.Client
	logger *logrus.Logger
}

// NewModEvent - Adds the capability to push log events back to the console, while retaining logrus capabilities.
func NewModEvent(client *clientpb.Client, logger *logrus.Logger) (hook *modEvent) {
	return
}

// Fire - Function needed to implement the logrus.TxtLogger interface
func (hook *modEvent) Fire(entry *logrus.Entry) (err error) {
	return
}

// Levels - Function needed to implement the logrus.TxtLogger interface
func (hook *modEvent) Levels() (levels []logrus.Level) {
	return logrus.AllLevels
}
