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

import "github.com/gobuffalo/packr"

// SetupAssets - General function for managing all assets necessary to Wiregost
func SetupAssets() error {
	return nil
}

// AssetVersion - Returns the current version (git commit hash) of extracted assets, if any.
func AssetVersion() (version string) {
	return
}

// SetupGoToolchain - Downloads and/or unpacks the Go toolchain for major OS/arch.
func SetupGoToolchain() error {
	return nil
}

// SetupCodenames - Generates codenames for ghost implants
func SetupCodenames() error {
	return nil
}

// UnzipGoDependency - Unzip a Go toolchain into its appropriate directory.
func UnzipGoDependency(filename, targetPath string, assetsBox packr.Box) error {
	return nil
}

// SetupDataPath - Sets the data directory up, for things like .NET hosting DLLs
func SetupDataPath(appDir string) error {
	return nil
}

// SetupGoPath - Extracts dependencies to Wiregost GoPath.
// Might not be needed though, as in Sliver this function mainly extracts
// go constant files from proto package, assumingly for template compilation needs.
func SetupGoPath(goPathSrc string) error {
	return nil
}

// unzip - Utility unzip function
func unzip(src, dest string) (filenames []string, err error) {
	return
}
