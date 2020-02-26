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
	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/util"
	. "github.com/maxlandon/wiregost/client/util"
	"github.com/maxlandon/wiregost/data_service/remote"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// AutoCompleter is the autocompletion engine
type OptionCompleter struct {
	Command *commands.Command
}

// Do is the completion function triggered at each line
func (oc *OptionCompleter) Do(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {

	// Complete command args
	splitLine := strings.Split(string(line), " ")
	line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))

	moduleOptions := []string{}
	moduleOptions = append(moduleOptions, util.SortGenerateOptionKeys(ctx.Module.Options)...)
	moduleOptions = append(moduleOptions, util.SortListenerOptionKeys(ctx.Module.Options)...)

	switch word := splitLine[0]; {
	case wordInOptions(word, moduleOptions):
		return oc.yieldOptionValues(ctx, word, line, pos)
	default:
		return oc.yieldOptionNames(ctx, line, pos)
	}
}

// Do is the completion function triggered at each line
func (oc *OptionCompleter) yieldOptionNames(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {
	switch *ctx.MenuContext {
	case "module":
		for _, v := range util.SortListenerOptionKeys(ctx.Module.Options) {
			search := v + " "
			if !hasPrefix(line, []rune(search)) {
				sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}

		for _, v := range util.SortGenerateOptionKeys(ctx.Module.Options) {
			search := v + " "
			if !hasPrefix(line, []rune(search)) {
				sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}
	}
	return options, offset
}

// Do is the completion function triggered at each line
func (oc *OptionCompleter) yieldOptionValues(ctx *commands.ShellContext, optionName string, line []rune, pos int) (options [][]rune, offset int) {

	switch optionName {
	case "StageImplant", "StageConfig":
		// Get ghost builds
		rpc := ctx.Server.RPC
		resp := <-rpc(&ghostpb.Envelope{
			Type: clientpb.MsgListGhostBuilds,
		}, defaultTimeout)
		if resp.Err != "" {
			fmt.Printf(RPCError+"%s\n", resp.Err)
			return
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
			search := c.Name
			if !hasPrefix(line, []rune(search)) {
				sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}
	case "Workspace":
		workspaces, _ := remote.Workspaces(nil)
		for _, w := range workspaces {
			search := w.Name
			if !hasPrefix(line, []rune(search)) {
				sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}
	case "Session":
		rpc := ctx.Server.RPC

		resp := <-rpc(&ghostpb.Envelope{
			Type: clientpb.MsgSessions,
			Data: []byte{},
		}, defaultTimeout)
		if resp.Err != "" {
			fmt.Printf(RPCError+"%s\n", resp.Err)
			return
		}
		sessions := &clientpb.Sessions{}
		proto.Unmarshal(resp.Data, sessions)

		ghosts := map[uint32]*clientpb.Ghost{}
		for _, ghost := range sessions.Ghosts {
			ghosts[ghost.ID] = ghost
		}

		for _, g := range ghosts {
			search := g.Name
			if !hasPrefix(line, []rune(search)) {
				sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}
	}

	return options, offset
}

func wordInOptions(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
