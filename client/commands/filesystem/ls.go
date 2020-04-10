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

package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/util"
)

// ListDirectories - List directory contents
type ListDirectoriesCmd struct {
	Positional struct {
		Path string `description:"Remote directory/file"`
	} `positional-args:"yes"`
}

var ListDirectories ListDirectoriesCmd

func RegisterGhostLs() {
	GhostParser.AddCommand(constants.GhostLs, "", "", &ListDirectories)

	ls := GhostParser.Find(constants.GhostLs)
	ls.ShortDescription = "List remote directory contents"
}

// Execute - Command
func (ls *ListDirectoriesCmd) Execute(args []string) error {

	rpc := Context.Server.RPC

	path := ls.Positional.Path
	if (path == "~" || path == "~/") && Context.Ghost.OS == "linux" {
		path = filepath.Join("/home", Context.Ghost.Username)
	}
	if strings.HasPrefix(path, "~") {
		path = filepath.Join("/home", Context.Ghost.Username, strings.TrimPrefix(path, "~"))
	}

	data, _ := proto.Marshal(&ghostpb.LsReq{
		GhostID: Context.Ghost.ID,
		Path:    path,
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgLsReq,
		Data: data,
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}

	dirList := &ghostpb.Ls{}
	err := proto.Unmarshal(resp.Data, dirList)
	if err != nil {
		fmt.Printf(Error+"Unmarshaling envelope error: %v\n", err)
		return nil
	}
	printDirList(dirList)

	return nil
}

func printDirList(dirList *ghostpb.Ls) {
	fmt.Printf("Listing dir: %s%s%s\n", tui.BOLD, dirList.Path, tui.RESET)
	fmt.Printf("%s\n", strings.Repeat(tui.Dim("-"), len(dirList.Path)))

	table := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	for _, fileInfo := range dirList.Files {
		if fileInfo.IsDir {
			fmt.Fprintf(table, "%s\t<dir>\t\n", fileInfo.Name)
		} else {
			fmt.Fprintf(table, "%s\t%s\t\n", fileInfo.Name, util.ByteCountBinary(fileInfo.Size))
		}
	}
	table.Flush()
}
