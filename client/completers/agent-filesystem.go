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
	"path/filepath"
	"strings"

	"github.com/gogo/protobuf/proto"

	"github.com/maxlandon/wiregost/client/commands"
	. "github.com/maxlandon/wiregost/client/util"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// AutoCompleter is the autocompletion engine
type ImplantPathCompleter struct {
	Command *commands.Command
}

// Do is the completion function triggered at each line
func (pc *ImplantPathCompleter) Do(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {

	splitLine := strings.Split(string(line), " ")
	line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))

	// 1) Get the absolute path. There are two cases:
	//      - The path is "rounded" with a slash: no filter to keep.
	//      - The path is not a slash: a filter to keep for later.
	// We keep a boolean for remembering which case we found
	linePath := ""
	// path := ""
	lastPath := ""
	if strings.HasSuffix(string(line), "/") {
		// Trim the non needed slash
		linePath = strings.TrimSuffix(string(line), "/")
		// linePath = filepath.Dir(string(line))
		// Get absolute path
		// path, _ = fs.Expand(string(linePath))

	} else if string(line) == "" {
		linePath = "."
	} else {
		linePath = string(line)
		// linePath = filepath.Dir(string(line))
		// Get absolute path
		// path, _ = fs.Expand(string(linePath))
		// Save filter
		lastPath = filepath.Base(string(line))
	}

	// 2) We take the absolute path we found, and get all dirs in it.
	var dirs []string

	rpc := ctx.Server.RPC
	data, _ := proto.Marshal(&ghostpb.LsReq{
		GhostID: ctx.CurrentAgent.ID,
		Path:    linePath,
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgLsReq,
		Data: data,
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}

	dirList := &ghostpb.Ls{}
	err := proto.Unmarshal(resp.Data, dirList)
	if err != nil {
		fmt.Printf(Error+"Unmarshaling envelope error: %v\n", err)
		return
	}

	for _, fileInfo := range dirList.Files {
		if fileInfo.IsDir {
			dirs = append(dirs, fileInfo.Name)
		}
	}

	switch lastPath {
	case "":
		for _, dir := range dirs {
			search := ""
			if ctx.CurrentAgent.OS == "windows" {
				search = dir + "\\"
			} else {
				search = dir + "/"
			}
			if !hasPrefix([]rune(lastPath), []rune(search)) {
				sLine, sOffset := doInternal([]rune(lastPath), pos, len([]rune(lastPath)), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}
	default:
		filtered := []string{}
		for _, dir := range dirs {
			if strings.HasPrefix(dir, lastPath) {
				filtered = append(filtered, dir)
			}
		}

		for _, dir := range filtered {
			search := dir + "/"
			if !hasPrefix([]rune(lastPath), []rune(search)) {
				sLine, sOffset := doInternal([]rune(lastPath), pos, len([]rune(lastPath)), []rune(search))
				options = append(options, sLine...)
				offset = sOffset
			}
		}
	}

	return options, offset
}
