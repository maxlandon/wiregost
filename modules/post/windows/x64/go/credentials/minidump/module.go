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

package minidump

import (
	"errors"
	"fmt"
	"path/filepath"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	pb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module/templates"
)

var rpcLog = log.ServerLogger("windows/x64/go/credentials/minidump", "module")

// metadataFile - Full path to module metadata
var metadataFile = filepath.Join(assets.GetModulesDir(), "post/windows/x64/go/credentials/minidump/metadata.json")

// [ Base Methods ] ------------------------------------------------------------------------//

// Minidump - A single stage DNS implant
type Minidump struct {
	Base *templates.Module
}

// New - Instantiates a reverse DNS module, empty.
func New() *Minidump {
	return &Minidump{Base: &templates.Module{}}
}

// Init - Module initialization, loads metadata. ** DO NOT ERASE **
func (s *Minidump) Init() error {
	return s.Base.Init(metadataFile)
}

// ToProtobuf - Returns protobuf version of module
func (s *Minidump) ToProtobuf() *pb.Module {
	return s.Base.ToProtobuf()
}

// SetOption - Sets a module option through its base object.
func (s *Minidump) SetOption(option, name string) {
	s.Base.SetOption(option, name)
}

// Run - Module entrypoint. ** DO NOT ERASE **
func (s *Minidump) Run(command string) (result string, err error) {

	// Check session
	if ok, err := s.checkSession(); !ok {
		return "", err
	}

	return "Minidump module executed", nil
}

func (s *Minidump) checkSession() (ok bool, err error) {

	// Check empty session
	if s.Base.Options["Session"].Value == "" {
		return false, errors.New("Provide a Session to run this module on.")
	}

	// Check valid session
	sessions := &clientpb.Sessions{}
	if 0 < len(*core.Wire.Ghosts) {
		for _, ghost := range *core.Wire.Ghosts {
			sessions.Ghosts = append(sessions.Ghosts, ghost.ToProtobuf())
		}
	}
	found := false
	for _, sess := range sessions.Ghosts {
		if sess.Name == s.Base.Options["Session"].Value {
			found = true
		}
	}

	if found {
		return true, nil
	} else {
		invalid := fmt.Sprintf("Invalid or non-connected session: %s", s.Base.Options["Session"].Value)
		return false, errors.New(invalid)
	}
}
