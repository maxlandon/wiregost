package db

// Wiregost - Post-Exploitation & Implant Framework
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

import (
	"strings"
	"time"

	"gorm.io/gorm/logger"

	"github.com/maxlandon/wiregost/internal/server/configs"
	"github.com/maxlandon/wiregost/internal/server/log"
)

var gormLog = log.NamedLogger("db", "gorm")

type gormWriter struct{}

func (w gormWriter) Printf(format string, args ...interface{}) {
	gormLog.Printf(format, args...)
}

func getGormLogger(dbConfig *configs.DatabaseConfig) logger.Interface {
	logConfig := logger.Config{
		SlowThreshold: time.Second,
		Colorful:      true,
		LogLevel:      logger.Info,
	}
	switch strings.ToLower(dbConfig.LogLevel) {
	case "silent":
		logConfig.LogLevel = logger.Silent
	case "err":
		fallthrough
	case "error":
		logConfig.LogLevel = logger.Error
	case "warning":
		fallthrough
	case "warn":
		logConfig.LogLevel = logger.Warn
	case "info":
		fallthrough
	default:
		logConfig.LogLevel = logger.Info
	}

	return logger.New(gormWriter{}, logConfig)
}
