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
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// ProfilesCmd - List ghost profiles
type ProfilesCmd struct{}

var Profiles ProfilesCmd

func RegisterProfiles() {
	MainParser.AddCommand(constants.Profiles, "", "", &Profiles)

	pro := MainParser.Find(constants.Profiles)
	CommandMap[MAIN_CONTEXT] = append(CommandMap[MAIN_CONTEXT], pro)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], pro)
	pro.ShortDescription = "List ghost profiles"
	pro.SubcommandsOptional = true
}

// Execute - List ghost profiles
func (p *ProfilesCmd) Execute(args []string) error {
	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgProfiles,
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	pbProfiles := &clientpb.Profiles{}
	err := proto.Unmarshal(resp.Data, pbProfiles)
	if err != nil {
		fmt.Printf(Error+"%s \n", err.Error())
		return nil
	}

	profiles := &map[string]*clientpb.Profile{}
	for _, profile := range pbProfiles.List {
		(*profiles)[profile.Name] = profile
	}

	headers := []string{"Name", "Platform", "Format", "Command & Control", "Limitations", "Debug"}
	widths := []int{15, 15, 10, 30, 30, 5}

	tab := util.NewTable()
	tab.SetColumns(headers, widths)
	tab.SetColWidth(40)

	for k, p := range *profiles {
		platform := fmt.Sprintf("%s/%s", p.Config.GOOS, p.Config.GOARCH)
		c2s := []string{}
		for _, c := range p.Config.C2 {
			c2s = append(c2s, c.URL)
		}
		limits := getLimitsString(p.Config)
		tab.Append([]string{k, platform, p.Config.Format.String(), strings.Join(c2s, ","), limits, strconv.FormatBool(p.Config.Debug)})
	}

	tab.Render()

	return nil
}

func getLimitsString(config *clientpb.GhostConfig) string {
	limits := []string{}
	if config.LimitDatetime != "" {
		limits = append(limits, fmt.Sprintf("datetime=%s", config.LimitDatetime))
	}
	if config.LimitDomainJoined {
		limits = append(limits, fmt.Sprintf("domainjoined=%v", config.LimitDomainJoined))
	}
	if config.LimitUsername != "" {
		limits = append(limits, fmt.Sprintf("username=%s", config.LimitUsername))
	}
	if config.LimitHostname != "" {
		limits = append(limits, fmt.Sprintf("hostname=%s", config.LimitHostname))
	}
	return strings.Join(limits, "; ")
}
