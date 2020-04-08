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
	"sort"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/lmorg/readline"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func completeModulePath(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get last path
	splitLine := strings.Split(string(line), " ")
	last := splitLine[len(splitLine)-1]

	stack, _ := proto.Marshal(&clientpb.ModuleActionReq{
		Action: constants.ModuleList,
	})

	rpc := Context.Server.RPC

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgModuleList,
		Data: stack,
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError, "%s\n", resp.Err)
		return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
	}

	modList := &clientpb.ModuleAction{}
	proto.Unmarshal(resp.Data, modList)

	sort.Strings(modList.Modules)

	for _, mod := range modList.Modules {
		if strings.HasPrefix(mod, string(last)) {
			suggestions = append(suggestions, mod[len(last):])
		}
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
}
