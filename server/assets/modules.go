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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/gobuffalo/packr"
)

var (
	// ModuleBox - Contains all static assets for modules
	ModuleBox = packr.NewBox("../../modules/")
)

func GetModulesDir() string {
	dir := path.Join(GetRootAppDir(), moduleDirPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			setupLog.Fatalf("Cannot write to wiregost root directory %s", err)
		}
	}
	return dir
}

func setupModules() error {

	setupLog.Infof("Unpacking modules data to '%s'", GetModulesDir())
	modulesPath := GetModulesDir()

	// All modules files are extracted from box and written to disk
	files := ModuleBox.List()

	for _, file := range files {
		script, _ := ModuleBox.Find(file)
		modFilePath := filepath.Join(modulesPath, file)
		if _, err := os.Stat(modFilePath); os.IsNotExist(err) {
			os.RemoveAll(modFilePath)
		}
		dirPath := filepath.Dir(modFilePath)
		os.MkdirAll(dirPath, os.ModePerm)
		ioutil.WriteFile(modFilePath, script, 0644)
	}

	return nil
}
