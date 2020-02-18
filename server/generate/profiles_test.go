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
	"testing"
)

func TestProfileByName(t *testing.T) {
	name := "foobar"
	config := &GhostConfig{
		GOOS:   "windows",
		GOARCH: "amd64",
	}
	err := ProfileSave(name, config)
	if err != nil {
		t.Errorf("%v", err)
	}

	profile, err := ProfileByName(name)
	if err != nil {
		t.Errorf("%v", err)
	}

	if profile.GOOS != config.GOOS {
		t.Errorf("Fetched data does not match saved data")
	}
}
