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

package sessions

import (
	"fmt"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
)

// SessionBackgroundCmd - Put the current ghost session in the background
type SessionBackgroundCmd struct{}

var SessionBackground SessionBackgroundCmd

func RegisterSessionBackground() {
	GhostParser.AddCommand(constants.SessionsBackground, "", "", &SessionBackground)

	back := GhostParser.Find(constants.SessionsBackground)
	CommandMap[GHOST_CONTEXT] = append(CommandMap[GHOST_CONTEXT], back)
	back.ShortDescription = "Put the current ghost session in the background"
}

// Execute - Put the current ghost session in the background
func (sb *SessionBackgroundCmd) Execute(args []string) error {

	fmt.Printf(Info+"Background session %d (%s) ... \n", Context.Ghost.ID, Context.Ghost.Name)
	*Context.Ghost = clientpb.Ghost{}

	return nil
}
