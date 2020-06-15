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
	DBName     string
	DBUser     string
	DBPassword string
	// Database as Service
	DatabaseRPCHost     string
	DatabaseRPCPort     uint
	DatabaseRESTAddress string
}

// LoadServerConfiguration - Loads config from the config file, handling any other cases.
func LoadServerConfiguration() error {
	return nil
}

// SaveServerConfiguration - Write the config file to disk.
func SaveServerConfiguration() error {
	return nil
}
