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
	"fmt"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"
	"github.com/lmorg/readline"

	. "github.com/maxlandon/wiregost/client/commands"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func CompleteProfileNames(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get last path
	splitLine := strings.Split(string(line), " ")
	last := splitLine[len(splitLine)-1]

	// Get profiles
	rpc := Context.Server.RPC

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgProfiles,
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError, "%s\n", resp.Err)
		return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
	}

	pbProfiles := &clientpb.Profiles{}
	err := proto.Unmarshal(resp.Data, pbProfiles)
	if err != nil {
		fmt.Println()
		fmt.Printf(Error, "%s", err.Error())
		return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
	}

	for _, p := range pbProfiles.List {
		if strings.HasPrefix(p.Name, string(last)) {
			suggestions = append(suggestions, p.Name[(len(last)):])

			os := osPad(p.Config, pbProfiles)
			format := formatPad(p.Config, pbProfiles)
			c2s := c2Pad(p.Config, pbProfiles)
			desc := tui.Dim(os + " - " + format + " - " + c2s)
			listSuggestions[p.Name[(len(last)):]] = desc
		}
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayList
}

func osPad(p *clientpb.GhostConfig, profs *clientpb.Profiles) string {
	var max int
	for _, prof := range profs.List {
		if len(prof.Config.GOOS+"/"+prof.Config.GOARCH) > max {
			max = len(prof.Config.GOOS + "/" + prof.Config.GOARCH)
		}
	}
	var pad string
	for i := 0; i < max-len(p.GOOS+"/"+p.GOARCH); i++ {
		pad += " "
	}

	return p.GOOS + "/" + p.GOARCH + pad
}

func formatPad(p *clientpb.GhostConfig, profs *clientpb.Profiles) string {
	var max int
	var pFormat string
	for _, prof := range profs.List {
		var format string
		switch prof.Config.Format {
		case clientpb.GhostConfig_EXECUTABLE:
			format = "exe"
		case clientpb.GhostConfig_SHARED_LIB:
			format = "shared"
		case clientpb.GhostConfig_SHELLCODE:
			format = "shellcode"
		}
		if len(format) > max {
			max = len(format)
		}
	}
	switch p.Format {
	case clientpb.GhostConfig_EXECUTABLE:
		pFormat = "exe"
	case clientpb.GhostConfig_SHARED_LIB:
		pFormat = "shared"
	case clientpb.GhostConfig_SHELLCODE:
		pFormat = "shellcode"
	}

	var pad string
	for i := 0; i < max-len(pFormat); i++ {
		pad += " "
	}

	return pFormat + pad
}

func c2Pad(p *clientpb.GhostConfig, profs *clientpb.Profiles) string {
	var max int
	for _, prof := range profs.List {
		var c2s string
		for _, c2 := range prof.Config.C2 {
			c2s += "| " + c2.URL
		}
		if len(c2s) > max {
			max = len(c2s)
		}
	}

	var pc2s string
	for _, c2 := range p.C2 {
		pc2s += "| " + c2.URL
	}

	var pad string
	for i := 0; i < max-len(pc2s); i++ {
		pad += " "
	}

	return pc2s + pad
}
