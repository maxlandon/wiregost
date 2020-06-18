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

package assets

import (
	"os"
	"os/user"
	"path"
)

const (
	// GoDirName - The directory to store the go compiler/toolchain files in
	GoDirName       = "go"
	goPathDirName   = "gopath"
	versionFileName = "version"
	dataDirName     = "data"
	envVarName      = "WIREGOST_ROOT_DIR"
	moduleDirPath   = "modules"
	stagersDirName  = "stagers"
)

// GetRootAppDir - Returns the root directory for Wiregost data. Creates it if needed.
func GetRootAppDir() (dir string) {

	value := os.Getenv(envVarName)

	if len(value) == 0 {
		user, _ := user.Current()
		dir = path.Join(user.HomeDir, ".wiregost")
	} else {
		dir = value
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			// setupLog.Fatalf("Cannot write to wiregost root directory %s", err)
		}
	}
	return
}

// GetStagersDir - Returns the directory where stager files are generated.
func GetStagersDir() (dir string) {
	return
}

// GetDataDir - Returns the directory for data
func GetDataDir() (dir string) {
	return
}

// GetDatabaseDir - Get the root directory where all DB-related files are. Creates it if needed.
func GetDatabaseDir() (dir string) {
	return
}

// GetModulesDir - Returns the directory where all unpacked module source code is stored.
func GetModulesDir() (dir string) {
	return
}

// GetGhostDir - Each ghost has its own directory for binaries, log and other data. Find it.
func GetGhostDir(workspaceID uint32, ghostName string) (dir string) {
	return
}
