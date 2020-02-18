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

package generate

import (
	"os"
	"runtime"
)

func getCCompiler(arch string) string {
	var found bool // meh, ugly
	var compiler string
	if arch == "amd64" {
		compiler = os.Getenv(GhostCC64EnvVar)
	}
	if arch == "386" {
		compiler = os.Getenv(GhostCC32EnvVar)
	}
	if compiler == "" {
		if compiler, found = defaultMingwPath[arch]; !found {
			compiler = defaultMingwPath["amd64"] // should not happen, but just in case ...
		}
	}
	if _, err := os.Stat(compiler); os.IsNotExist(err) {
		buildLog.Warnf("CC path %v does not exist", compiler)
		return ""
	}
	if runtime.GOOS == "windows" {
		compiler = "" // TODO: Add windows mingw support
	}
	buildLog.Infof("CC = %v", compiler)
	return compiler
}
