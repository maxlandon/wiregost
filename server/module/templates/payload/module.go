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

// CHANGE THE NAME OF THE PACKAGE WITH THE NAME OF YOUR MODULE/DIRECTORY
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/module/templates"
)

// metadataFile - Full path to module metadata
var metadataFile = filepath.Join(assets.GetModulesDir(), "payload/path/to/metadata.json")

// Change the ModuleTypeName to the same name as your Module name above
type PayloadModule struct {
	templates.Module
}

// Instantiates a base module - Do not modify
func New() *PayloadModule {
	return &PayloadModule{templates.Module{}}
}

// Init - Module initialization, loads metadata. ** DO NOT ERASE **
func (s *PayloadModule) Init() error {
	file, err := os.Open(metadataFile)
	if err != nil {
		return err
	}

	metadata, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(metadata, s)
	if err != nil {
		return err
	}

	return nil
}

// Run - Module entrypoint. ** DO NOT ERASE **
func (s *PayloadModule) Run() error {

	return nil
}
