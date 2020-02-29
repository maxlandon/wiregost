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

package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/server/assets"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	// Server address
	LHost string
	LPort int

	// Metasploit
	MsfDirPath string
}

func LoadServerConfig() *ServerConfig {
	// Load a default console config, eventually parse one if found
	conf := &ServerConfig{
		LHost:      "localhost",
		LPort:      1708,
		MsfDirPath: "",
	}

	// Load config
	file := filepath.Join(assets.GetRootAppDir(), "server.yaml")

	if _, err := os.Stat(file); os.IsNotExist(err) {
		err = SaveServerConfig(conf)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		err = SaveServerConfig(conf)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(tui.Red("[!] Error: failed to unmarshal server.yaml file."))
	}

	return conf

}

func SaveServerConfig(config *ServerConfig) error {
	saveTo := assets.GetRootAppDir()
	configYAML, _ := yaml.Marshal(config)

	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		err = os.MkdirAll(saveTo, os.ModePerm)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot write to wiregost root directory %s", err))
		}
	}

	fi, err := os.Stat(saveTo)
	if fi.IsDir() {
		filename := "server.yaml"
		saveTo = filepath.Join(saveTo, filename)
	}

	err = ioutil.WriteFile(saveTo, configYAML, 0600)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write config to: %s (%v) \n", saveTo, err))
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

var helpConfig = `

# [ HELP ] ------------------------------------------------------//

# SERVER ---
# lhost:            IP Address to listen client connections on. 
# lport:            Listen port 

# MSF -----
# msfdirpath:       If Metasploit is installed from source and not accessible from path,
#                   specify the path to the repo here. Ignore if not needed
`
