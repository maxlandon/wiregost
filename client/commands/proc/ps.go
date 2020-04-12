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

package proc

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// ProcOptions - Filters available to proc commands
type ProcOptions struct {
	PID   int    `long:"pid" description:"Process ID" default:"-1"`
	Owner string `long:"owner" description:"Process owner"`
	Exe   string `long:"exe" description:"Process name"`
}

// PsCmd - "List processes running on the target, with (--optional) filters"
type PsCmd struct {
	ProcOptions `group:"Process filters"`
}

var Ps PsCmd

func RegisterPs() {
	GhostParser.AddCommand(constants.Ps, "", "", &Ps)

	ps := GhostParser.Find(constants.Ps)
	ps.ShortDescription = "List processes running on the target, with (--optional) filters"
}

// Execute - Command
func (p *PsCmd) Execute(args []string) error {

	data, _ := proto.Marshal(&ghostpb.PsReq{GhostID: Context.Ghost.ID})
	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: ghostpb.MsgPsReq,
		Data: data,
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}
	ps := &ghostpb.Ps{}
	err := proto.Unmarshal(resp.Data, ps)
	if err != nil {
		fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
		return nil
	}

	outputBuf := bytes.NewBufferString("")
	table := tabwriter.NewWriter(outputBuf, 0, 2, 2, ' ', 0)

	fmt.Fprintf(table, "%sPID\tPPID\tExecutable\tOwner%s\t\n", tui.YELLOW, tui.RESET)
	fmt.Fprintf(table, "%s%s\t%s\t%s\t%s%s\t\n",
		tui.YELLOW,
		strings.Repeat("-", len("pid")),
		strings.Repeat("-", len("ppid")),
		strings.Repeat("-", len("executable")),
		strings.Repeat("-", len("owner")),
		tui.RESET,
	)

	// filters
	pidFilter := p.PID
	exeFilter := p.Exe
	ownerFilter := p.Owner

	lineColors := []string{}
	for _, proc := range ps.Processes {
		var lineColor = ""
		if pidFilter != -1 && proc.Pid == int32(pidFilter) {
			lineColor = printProcInfo(table, Context, proc)
		}
		if exeFilter != "" && strings.Contains(proc.Executable, exeFilter) {
			lineColor = printProcInfo(table, Context, proc)
		}
		if ownerFilter != "" && strings.HasPrefix(proc.Owner, ownerFilter) {
			lineColor = printProcInfo(table, Context, proc)
		}
		if pidFilter == -1 && exeFilter == "" && ownerFilter == "" {
			lineColor = printProcInfo(table, Context, proc)
		}

		// Should be set to normal/green if we rendered the line
		if lineColor != "" {
			lineColors = append(lineColors, lineColor)
		}
	}
	table.Flush()

	for index, line := range strings.Split(outputBuf.String(), "\n") {
		if len(line) == 0 {
			continue
		}
		// We need to account for the two rows of column headers
		if 0 < len(line) && 2 < index {
			lineColor := lineColors[index-2]
			fmt.Printf("%s%s%s\n", lineColor, line, tui.FOREWHITE)
		} else {
			fmt.Printf("%s\n", line)
		}
	}

	// Reset options
	p.PID = 0
	p.Owner = ""
	p.Exe = ""

	return nil
}

// // printProcInfo - Stylizes the process information
func printProcInfo(table *tabwriter.Writer, Context ShellContext, proc *ghostpb.Process) string {
	color := tui.FOREWHITE
	if modifyColor, ok := knownProcs[proc.Executable]; ok {
		color = modifyColor
	}
	if Context.Ghost.Name != "" && proc.Pid == Context.Ghost.PID {
		color = tui.GREEN
	}
	fmt.Fprintf(table, "%d\t%d\t%s\t%s\t\n", proc.Pid, proc.Ppid, proc.Executable, proc.Owner)
	return color
}
