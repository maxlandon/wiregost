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
	"log"
	"os"
	"os/user"
	"path"
)

const (
	// GoDirName - The directory to store the go compiler/toolchain files in
	envVarName          = "SLIVER_ROOT_DIR"
	dataServiceDir      = "data-service"
	dataServiceCertsDir = "certs"
)

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
			log.Fatalf("Cannot write to wiregost root directory %s", err)
		}
	}
	return dir
}

// GetDataServiceDir - Returns the full path to the data service directory
func GetDataServiceDir() string {
	dir := path.Join(GetRootAppDir(), dataServiceDir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatalf("Cannot write to Wiregost Data Service directory %s", err)
		}
	}

	return dir
}

// GetDataServiceCertsDir - Returns the full path to the data service certs directory
func GetDataServiceCertsDir() string {
	dir := path.Join(GetDataServiceDir(), dataServiceCertsDir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatalf("Cannot write to Wiregost Data Service directory %s", err)
		}
	}

	return dir
}
