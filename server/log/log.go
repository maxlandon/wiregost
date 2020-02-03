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
	"errors"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	envVarName = "WIREGOST_ROOT_DIR"
)

// [ Logger Hooks ] -----------------------------------------------------//

type TxtHook struct {
	Name   string
	logger *logrus.Logger
}

func NewTxtHook(name string, logger *logrus.Logger) *TxtHook {
	hook := &TxtHook{
		Name:   name,
		logger: logger,
	}
	return hook
}

func (hook *TxtHook) Fire(entry *logrus.Entry) error {
	if hook.logger == nil {
		return errors.New("No text logger")
	}

	// Determine the caller (filename/line number)
	srcFile := "<no caller>"
	if entry.HasCaller() {
		ghostIndex := strings.Index(entry.Caller.File, "wiregost")
		srcFile = entry.Caller.File
		if ghostIndex != -1 {
			srcFile = srcFile[ghostIndex:]
		}
	}

	switch entry.Level {
	case logrus.PanicLevel:
		hook.logger.Panicf("[%s:%d] %s", srcFile, entry.Caller.Line, entry.Message)
	case logrus.FatalLevel:
		hook.logger.Fatalf("[%s:%d] %s", srcFile, entry.Caller.Line, entry.Message)
	case logrus.ErrorLevel:
		hook.logger.Errorf("[%s:%d] %s", srcFile, entry.Caller.Line, entry.Message)
	case logrus.WarnLevel:
		hook.logger.Warnf("[%s:%d] %s", srcFile, entry.Caller.Line, entry.Message)
	case logrus.InfoLevel:
		hook.logger.Infof("[%s:%d] %s", srcFile, entry.Caller.Line, entry.Message)
	case logrus.DebugLevel:
		hook.logger.Debugf("[%s:%d] %s", srcFile, entry.Caller.Line, entry.Message)
	}

	return nil
}

func (hook *TxtHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// [ Directories ] ------------------------------------------------------//

// GetRootAppDir - Get the Wiregost app dir, default is: ~/.wiregost/
func GetRootAppDir() string {

	value := os.Getenv(envVarName)

	var dir string
	if len(value) == 0 {
		user, _ := user.Current()
		dir = path.Join(user.HomeDir, ".wiregost")
	} else {
		dir = value
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic("Cannot write to Wiregost root directory")
		}
	}
	return dir
}

// GetLogDir - Return the log dir
func GetLogDir() string {
	rootDir := GetRootAppDir()
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		err = os.MkdirAll(rootDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	logDir := path.Join(rootDir, "logs")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return logDir
}
