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

package handlers

import (
	"io/ioutil"
	"log"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"gopkg.in/yaml.v2"

	"github.com/maxlandon/wiregost/data_service/models"
)

// Env contains all configuration options for DB access and data web service. It passes its DB connection
// pool to HTTP handlers that need it, and its HTTP service parameters to remote/ functions.
type Env struct {
	DB *models.DB
	// Database
	Database struct {
		DbName     string `yaml:"db_name"`
		DbUser     string `yaml:"db_user"`
		DbPassword string `yaml:"db_password"`
	}

	// Web service
	Service struct {
		Address     string `yaml:"address"`
		Port        int    `yaml:"port"`
		URL         string `yaml:"url"`
		Certificate string `yaml:"certificate"`
		Key         string `yaml:"key"`
	}
}

// LoadEnv instantiates an Env object and populates it with the config.yaml file
func LoadEnv() *Env {
	env := &Env{}

	// Load config
	file, _ := fs.Expand("~/pentest/wiregost/data_service/config.yaml")

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(tui.Red("[!] Error: failed to read config.yaml file."))
	}

	err = yaml.Unmarshal(data, &env)
	if err != nil {
		log.Fatal(tui.Red("[!] Error: failed to unmarshal config.yaml file."))
	}

	// Adjust for certificate and key file paths
	env.Service.Certificate, err = fs.Expand(env.Service.Certificate)
	env.Service.Key, err = fs.Expand(env.Service.Key)

	// Connect to postgreSQL
	env.DB = models.New(env.Database.DbName, env.Database.DbUser, env.Database.DbPassword)

	return env
}
