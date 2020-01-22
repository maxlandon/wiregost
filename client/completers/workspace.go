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

package completers

import (
	"github.com/maxlandon/wiregost/data_service/remote"
)

func CompleteWorkspaces() func(prefix string, args []string) []string {

	return func(string, []string) []string {
		workspaces, _ := remote.Workspaces(nil)
		var names []string
		for _, w := range workspaces {
			names = append(names, w.Name)
		}
		return names
	}
}

func CompleteWorkspacesAndFlags() func(prefix string, args []string) []string {

	return func(string, []string) []string {
		workspaces, _ := remote.Workspaces(nil)
		var names []string
		for _, w := range workspaces {
			names = append(names, w.Name)
		}
		// Append flags
		for _, i := range []string{"--description", "--boundary", "--limit_to_network"} {
			names = append(names, i)
		}
		return names
	}
}
