package log

import "github.com/sirupsen/logrus"

// HandleRemoteLog - When remote logging is activated, this function handles how to store the logs,
// send them, find the good strategies and timing of reports.
func HandleRemoteLog(entry *logrus.Entry) (err error) {

	// Check connection and check report timing strategies/constraints

	// If not connected yet, or constraints, store the log event

	// If connected and constraints ok, send the log or group of logs

	return
}
