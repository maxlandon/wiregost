package log

import "github.com/sirupsen/logrus"

var (
	// ClientLogger - Audits all clients connections/disconnections
	ClientLogger = newAuditLogger()
)

// newAuditLogger - Instantiates a logger for client connections/disconnections
func newAuditLogger() (logger *logrus.Logger) {
	return
}
