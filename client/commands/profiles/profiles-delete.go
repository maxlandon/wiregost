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

package profiles

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// ProfilesDeleteCmd - Delete a ghost profile
type ProfilesDeleteCmd struct {
	Positional struct {
		Name string `description:"Profile name" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var ProfilesDelete ProfilesDeleteCmd

func RegisterProfilesDelete() {
	pro := MainParser.Find(constants.Profiles)
	pro.AddCommand(constants.ProfilesDelete, "", "", &ProfilesDelete)
	del := pro.Find(constants.ProfilesDelete)

	del.ShortDescription = "Delete a ghost profile"
	del.Args()[0].RequiredMaximum = 1
}

// Execute - Delete a ghost profile
func (pd *ProfilesDeleteCmd) Execute(args []string) error {
	pReq, _ := proto.Marshal(&clientpb.Profile{
		Name: pd.Positional.Name,
	})

	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgDeleteProfile,
		Data: pReq,
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	pRes := &clientpb.Profile{}
	err := proto.Unmarshal(resp.Data, pRes)
	if err != nil {
		fmt.Printf(Error+"%s\n", err.Error())
		return nil
	}

	if pRes.Name == pd.Positional.Name {
		fmt.Printf(Success+"Deleted profile %s\n", pd.Positional.Name)
	} else {
		fmt.Printf(Error+"%s\n", pRes.Name)
	}

	return nil
}
