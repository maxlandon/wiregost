package assets

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
	"os"
	"os/user"
	"path/filepath"

	"github.com/maxlandon/wiregost/internal/server/log"
)

const (
	// GoDirName - The directory to store the go compiler/toolchain files in
	GoDirName = "go"

	versionFileName = "version"
	envVarName      = "WIREGOST_ROOT_DIR"
)

var setupLog = log.NamedLogger("assets", "setup")

// GetRootAppDir - Get the Wiregost app dir, default is: ~/.wiregost/
func GetRootAppDir() string {
	value := os.Getenv(envVarName)
	var dir string
	if len(value) == 0 {
		user, _ := user.Current()
		dir = filepath.Join(user.HomeDir, ".wiregost")
	} else {
		dir = value
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0o700)
		if err != nil {
			setupLog.Fatalf("Cannot write to wiregost root dir %s", err)
		}
	}
	return dir
}
