package assets

import (
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

// WiregostClientDirName - Personal directory for the console
const WiregostClientDirName = ".wiregost-client"

// ConfigDirName - The configs directory
const ConfigDirName = "configs"

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

// GetConfigDir - Returns the path to the configuration directory
func GetConfigDir() string {
	rootDir, _ := filepath.Abs(GetRootAppDir())
	dir := path.Join(rootDir, ConfigDirName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dir
}
