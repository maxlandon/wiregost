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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/evilsocket/islazy/tui"
	"gopkg.in/yaml.v2"
)

var (
	// ServerConfiguration - The server configuration object
	ServerConfiguration = &serverConfiguration{}
)

// ServerConfiguration - Holds all configuration information for the server and the database
type serverConfiguration struct {
	// Server
	ServerHost        string
	ServerPort        int
	MetasploitDirPath string
	// Database
	DBName        string
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        uint
	DBCertificate string
	DBPrivateKey  string
	// Database as Service
	DatabaseRPCHost     string
	DatabaseRPCPort     uint
	DatabaseRESTAddress string
}

// LoadServerConfiguration - Loads config from the config file, handling any other cases.
func LoadServerConfiguration() (conf *serverConfiguration) {

	// Load a default console config, eventually parse one if found
	ServerConfiguration = &serverConfiguration{
		ServerHost:          "localhost",
		ServerPort:          1708,
		MetasploitDirPath:   "",
		DBName:              "wiregost_db",
		DBUser:              "wiregost",
		DBPassword:          "wiregost",
		DatabaseRPCHost:     "localhost",
		DBHost:              "localhost",
		DBPort:              5432,
		DBCertificate:       "",
		DBPrivateKey:        "",
		DatabaseRPCPort:     1710,
		DatabaseRESTAddress: "localhost:1712",
	}

	// Load config
	file := filepath.Join(GetRootAppDir(), "server.yaml")

	if _, err := os.Stat(file); os.IsNotExist(err) {
		err = SaveServerConfiguration(ServerConfiguration)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		err = SaveServerConfiguration(ServerConfiguration)
	}

	err = yaml.Unmarshal(data, &ServerConfiguration)
	if err != nil {
		log.Fatal(tui.Red("[!] Error: failed to unmarshal server.yaml file."))
	}

	return ServerConfiguration
}

// SaveServerConfiguration - Write the config file to disk.
func SaveServerConfiguration(config *serverConfiguration) error {

	saveTo := GetRootAppDir()
	configYAML, _ := yaml.Marshal(config)

	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		err = os.MkdirAll(saveTo, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Cannot write to wiregost root directory %s", err)
		}
	}

	fi, err := os.Stat(saveTo)
	if fi.IsDir() {
		filename := "server.yaml"
		saveTo = filepath.Join(saveTo, filename)
	}

	err = ioutil.WriteFile(saveTo, configYAML, 0600)
	if err != nil {
		return fmt.Errorf("Failed to write config to: %s (%v) \n", saveTo, err)
	}

	f, err := os.OpenFile(saveTo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(helpConfig); err != nil {
		panic(err)
	}

	return nil
}

// helpConfig - A documentation string appended to each server configuration file.
var helpConfig = `

# [ HELP ] ------------------------------------------------------//

# SERVER ---
# ServerHost:            IP Address to listen client connections on. 
# ServerPort:            Listen port 

# MSF -----
# msfdirpath:       If Metasploit is installed from source and not accessible from path,
#                   specify the path to the repo here. Ignore if not needed

# DATABASE ---
#DBName     PostgreSQL Database name used by Wiregost
#DBUser     PostgreSQL Database user owning the Wiregost DB
#DBPassword Password needed to access the DB
#DBHost     Host on which Wiregost can connect to PostgreSQL 
#DBPort     Port on which Wiregost can connect to PostgreSQL 

# DATABASE SERVICE ---
# DatabaseRPCHost     gRPC DB Host
# DatabaseRPCPort     gRPC DB Port
# PublicKeyDB         Not used through this file 
# PrivateKeyDB        Not used through this file 
# DatabaseRESTAddress gRPC REST Gateway 
`
