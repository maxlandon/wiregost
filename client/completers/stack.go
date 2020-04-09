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

	"github.com/gogo/protobuf/proto"
	"github.com/lmorg/readline"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func completeStackModulePath(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get last path
	splitLine := strings.Split(string(line), " ")
	last := splitLine[len(splitLine)-1]

	stack, _ := proto.Marshal(&clientpb.StackReq{
		Action:      constants.StackList,
		WorkspaceID: uint32(Context.Workspace.ID),
		User:        Context.Server.Config.User,
	})

	rpc := Context.Server.RPC

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgStackList,
		Data: stack,
	}, DefaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError, "%s\n", resp.Err)
		return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
	}

	stackList := &clientpb.Stack{}
	proto.Unmarshal(resp.Data, stackList)
	if stackList.Err != "" {
		fmt.Printf(Error, "%s\n", stackList.Err)
		return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
	}

	for _, mod := range stackList.Modules {
		name := strings.Join(mod.Path, "/")
		if strings.HasPrefix(name, string(last)) {
			suggestions = append(suggestions, name[len(last):])
		}
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
}
