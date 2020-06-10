package log

import "github.com/sirupsen/logrus"

// TxtHook - A hook for logging with text formatting
type TxtHook struct {
	Name   string
	logger *logrus.Logger
}

// NewTxtHook - New hook
func NewTxtHook(name string, logger *logrus.Logger) (hook *TxtHook) {
	return
}

// Fire - Function needed to implement the logrus.TxtLogger interface
func (hook *TxtHook) Fire(entry *logrus.Entry) (err error) {
	return
}

// Levels - Function needed to implement the logrus.TxtLogger interface
func (hook *TxtHook) Levels() (levels []logrus.Level) {
	return logrus.AllLevels
}
