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
	"path/filepath"
)

const (
	// WiregostClientDirName - Directory storing all of the client configs/logs
	WiregostClientDirName = ".wiregost-client"

	resourceDirName = "resource"
)

// GetRootAppDir - Get the Wiregost client app dir ~/.wiregost-client/
func GetRootAppDir() string {
	user, _ := user.Current()
	dir := path.Join(user.HomeDir, WiregostClientDirName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dir
}

// GetConfigDir - Returns the path to the config dir
func GetResourceDir() string {
	rootDir, _ := filepath.Abs(GetRootAppDir())
	dir := path.Join(rootDir, resourceDirName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dir
}
