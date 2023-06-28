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
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	ver "github.com/maxlandon/wiregost/internal/client/version"
)

const (
	// wiregostClientDirName - Directory storing all of the client configs/logs
	wiregostClientDirName = ".wiregost-client"

	versionFileName = "version"
)

// GetRootAppDir - Get the wiregost app dir ~/.wiregost-client/
func GetRootAppDir() string {
	user, _ := user.Current()
	dir := filepath.Join(user.HomeDir, wiregostClientDirName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0o700)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dir
}

// GetClientLogsDir - Get the wiregost client logs dir ~/.wiregost-client/logs/
func GetClientLogsDir() string {
	logsDir := filepath.Join(GetRootAppDir(), "logs")
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		err = os.MkdirAll(logsDir, 0o700)
		if err != nil {
			log.Fatal(err)
		}
	}
	return logsDir
}

// GetConsoleLogsDir - Get the wiregost client console logs dir ~/.wiregost-client/logs/console/
func GetConsoleLogsDir() string {
	consoleLogsDir := filepath.Join(GetClientLogsDir(), "console")
	if _, err := os.Stat(consoleLogsDir); os.IsNotExist(err) {
		err = os.MkdirAll(consoleLogsDir, 0o700)
		if err != nil {
			log.Fatal(err)
		}
	}
	return consoleLogsDir
}

func assetVersion() string {
	appDir := GetRootAppDir()
	data, err := os.ReadFile(filepath.Join(appDir, versionFileName))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func saveAssetVersion(appDir string) {
	versionFilePath := filepath.Join(appDir, versionFileName)
	fVer, _ := os.Create(versionFilePath)
	defer fVer.Close()
	fVer.Write([]byte(ver.GitCommit))
}

// Setup - Extract or create local assets
func Setup(force bool, echo bool) {
	appDir := GetRootAppDir()
	localVer := assetVersion()
	if force || localVer == "" || localVer != ver.GitCommit {
		saveAssetVersion(appDir)
	}
	if _, err := os.Stat(filepath.Join(appDir, settingsFileName)); os.IsNotExist(err) {
		SaveSettings(nil)
	}
}
