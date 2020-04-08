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
	"github.com/maxlandon/wiregost/client/commands/module"
	"github.com/maxlandon/wiregost/data-service/remote"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func CompleteOptionNames(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get last path
	splitLine := strings.Split(string(line), " ")
	last := splitLine[len(splitLine)-1]

	switch Context.Module.Type {
	case "payload":
		for _, v := range module.SortListenerOptionKeys(Context.Module.Options) {
			if strings.HasPrefix(v, string(last)) {
				suggestions = append(suggestions, v[len(last):]+" ")
			}
		}

		for _, v := range module.SortGenerateOptionKeys(Context.Module.Options) {
			if strings.HasPrefix(v, string(last)) {
				suggestions = append(suggestions, v[len(last):]+" ")
			}
		}
	case "post":
		for _, v := range module.SortPostOptions(Context.Module.Options) {
			if strings.HasPrefix(v, string(last)) {
				suggestions = append(suggestions, v[len(last):]+" ")
			}
		}
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
}

func CompleteOptionValues(optionName string, line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get last path
	splitLine := strings.Split(string(line), " ")
	last := splitLine[len(splitLine)-1]

	switch optionName {
	case "StageImplant", "StageConfig":
		// Get ghost builds
		rpc := Context.Server.RPC
		resp := <-rpc(&ghostpb.Envelope{
			Type: clientpb.MsgListGhostBuilds,
		}, DefaultTimeout)
		if resp.Err != "" {
			fmt.Printf(RPCError+"%s\n", resp.Err)
			return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
		}

		builds := &clientpb.GhostBuilds{}
		proto.Unmarshal(resp.Data, builds)
		shellcodeBuilds := []*clientpb.GhostConfig{}
		for _, c := range builds.Configs {
			if (c.Format == clientpb.GhostConfig_SHARED_LIB) || (c.Format == clientpb.GhostConfig_SHELLCODE) {
				shellcodeBuilds = append(shellcodeBuilds, c)
			}
		}

		for _, c := range shellcodeBuilds {
			if strings.HasPrefix(c.Name, string(last)) {
				suggestions = append(suggestions, c.Name[len(last):])
			}
		}
	case "Workspace":
		workspaces, _ := remote.Workspaces(nil)
		for _, w := range workspaces {
			if strings.HasPrefix(w.Name, string(last)) {
				suggestions = append(suggestions, w.Name[len(last):])
			}
		}
	case "Session":
		rpc := Context.Server.RPC

		resp := <-rpc(&ghostpb.Envelope{
			Type: clientpb.MsgSessions,
			Data: []byte{},
		}, DefaultTimeout)
		if resp.Err != "" {
			fmt.Printf(RPCError+"%s\n", resp.Err)
			return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
		}
		sessions := &clientpb.Sessions{}
		proto.Unmarshal(resp.Data, sessions)

		ghosts := map[uint32]*clientpb.Ghost{}
		for _, ghost := range sessions.Ghosts {
			ghosts[ghost.ID] = ghost
		}

		for _, g := range ghosts {
			if strings.HasPrefix(g.Name, string(last)) {
				suggestions = append(suggestions, g.Name[len(last):])
			}
		}
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
}

func wordInOptions(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
