// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package log

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

var (
	// ServerLoggerName - Server logger name, contains all server log data
	ServerLoggerName = "server"

	// ServLogger - Server logger
	ServLogger = serverLogger()
)

func ServerLogger(pkg, stream string) *logrus.Entry {
	return ServLogger.WithFields(logrus.Fields{
		"pkg":    pkg,
		"stream": stream,
	})
}

func serverLogger() *logrus.Logger {
	serverLogger := logrus.New()
	serverLogger.Formatter = &logrus.JSONFormatter{}
	jsonFilePath := path.Join(GetLogDir(), "server.json")
	jsonFile, err := os.OpenFile(jsonFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open server log file: %s", err))
	}

	serverLogger.Out = jsonFile
	serverLogger.SetLevel(logrus.DebugLevel)
	serverLogger.SetReportCaller(true)
	serverLogger.AddHook(NewTxtHook("server", serverTxtLogger()))

	return serverLogger
}

func serverTxtLogger() *logrus.Logger {
	serverLogger := logrus.New()
	serverLogger.Formatter = &logrus.TextFormatter{ForceColors: true}
	txtFilePath := path.Join(GetLogDir(), "server.log")
	txtFile, err := os.OpenFile(txtFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open server log file: %s", err))
	}

	serverLogger.Out = txtFile
	serverLogger.SetLevel(logrus.DebugLevel)

	return serverLogger
}
