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

	"github.com/maxlandon/wiregost/data_service/remote"
	"github.com/sirupsen/logrus"
)

var (
	// GhostLoggerName - Ghost logger name, contains all log data for a ghost
	GhostLoggerName = "ghost"
)

// [ Loggers ] --------------------------------------------------------//

// GhostLogger - returns a logger for a ghost
func GhostLogger(workspaceID uint, ghostName string) *logrus.Entry {
	var gLogger = ghostLogger(workspaceID, ghostName)
	return gLogger.WithFields(logrus.Fields{})
}

func ghostLogger(workspaceID uint, ghostName string) *logrus.Logger {
	ghostLogger := logrus.New()
	ghostLogger.Formatter = &logrus.JSONFormatter{}
	jsonFilePath := path.Join(GetGhostDir(workspaceID, ghostName), "ghost-log.json")
	jsonFile, err := os.OpenFile(jsonFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open ghost log file: %s", err))
	}

	ghostLogger.Out = jsonFile
	ghostLogger.SetLevel(logrus.DebugLevel)
	ghostLogger.SetReportCaller(true)
	ghostLogger.AddHook(NewTxtHook("ghost", serverTxtLogger()))

	return ghostLogger
}

func ghostTxtLogger(workspaceID uint, ghostName string) *logrus.Logger {
	ghostLogger := logrus.New()
	ghostLogger.Formatter = &logrus.TextFormatter{ForceColors: true}
	txtFilePath := path.Join(GetGhostDir(workspaceID, ghostName), "ghost.log")
	txtFile, err := os.OpenFile(txtFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open ghost log file: %s", err))
	}

	ghostLogger.Out = txtFile
	ghostLogger.SetLevel(logrus.DebugLevel)

	return ghostLogger
}

// [ Directories ] ------------------------------------------------------//

// GetGhostDir - Get directory for ghost implants
func GetGhostDir(workspaceID uint, ghostName string) string {
	rootDir := GetRootAppDir()
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		err = os.MkdirAll(rootDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	workspaces, _ := remote.Workspaces(nil)
	var workspaceName string
	for _, w := range workspaces {
		if w.ID == workspaceID {
			workspaceName = w.Name
		}
	}
	if workspaceName == "" {
		workspaceName = "default"
	}

	logDir := path.Join(rootDir, "workspaces", workspaceName, ghostName)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return logDir
}
