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

package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"

	"github.com/maxlandon/wiregost/client/spin"
	. "github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func RegisterProcCommands() {

	ps := &Command{
		Name: "ps",
		Args: []*CommandArg{
			&CommandArg{Name: "owner", Type: "string"},
			&CommandArg{Name: "pid", Type: "string"},
			&CommandArg{Name: "exe", Type: "boolean"},
		},
		Handle: func(r *Request) error {
			fmt.Println()
			ps(r.Args, *r.context, r.context.Server.RPC)
			return nil
		},
	}
	AddCommand("agent", ps)

	procdump := &Command{
		Name: "procdump",
		Args: []*CommandArg{
			&CommandArg{Name: "proc", Type: "string"},
			&CommandArg{Name: "pid", Type: "string"},
			&CommandArg{Name: "timeout", Type: "int"},
		},
		Handle: func(r *Request) error {
			fmt.Println()
			procdump(r.Args, *r.context, r.context.Server.RPC)
			return nil
		},
	}
	AddCommand("agent", procdump)

	terminate := &Command{
		Name: "terminate",
		Handle: func(r *Request) error {
			fmt.Println()
			terminate(r.Args, *r.context, r.context.Server.RPC)
			return nil
		},
	}
	AddCommand("agent", terminate)

	migrate := &Command{
		Name: "migrate",
		Handle: func(r *Request) error {
			fmt.Println()
			migrate(r.Args, *r.context, r.context.Server.RPC)
			return nil
		},
	}
	AddCommand("agent", migrate)
}

func ps(args []string, ctx ShellContext, rpc RPCServer) {

	opts := procFilters(args)

	var pidFilter int
	spid, found := opts["pid"]
	if found {
		pidFilter, _ = strconv.Atoi((spid.(string)))
	} else {
		pidFilter = -1
	}
	var ownerFilter string
	owner, found := opts["owner"]
	if found {
		ownerFilter = owner.(string)
	}
	var exeFilter string
	exe, found := opts["exe"]
	if found {
		exeFilter = exe.(string)
	}

	data, _ := proto.Marshal(&ghostpb.PsReq{GhostID: ctx.CurrentAgent.ID})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgPsReq,
		Data: data,
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}
	ps := &ghostpb.Ps{}
	err := proto.Unmarshal(resp.Data, ps)
	if err != nil {
		fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
		return
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

	lineColors := []string{}
	for _, proc := range ps.Processes {
		var lineColor = ""
		if pidFilter != -1 && proc.Pid == int32(pidFilter) {
			lineColor = printProcInfo(table, ctx, proc)
		}
		if exeFilter != "" && strings.HasPrefix(proc.Executable, exeFilter) {
			lineColor = printProcInfo(table, ctx, proc)
		}
		if ownerFilter != "" && strings.HasPrefix(proc.Owner, ownerFilter) {
			lineColor = printProcInfo(table, ctx, proc)
		}
		if pidFilter == -1 && exeFilter == "" && ownerFilter == "" {
			lineColor = printProcInfo(table, ctx, proc)
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
		if 0 < len(line) && 2 <= index {
			lineColor := lineColors[index-2]
			fmt.Printf("%s%s%s\n", lineColor, line, tui.FOREWHITE)
		} else {
			fmt.Printf("%s\n", line)
		}
	}

}

// printProcInfo - Stylizes the process information
func printProcInfo(table *tabwriter.Writer, ctx ShellContext, proc *ghostpb.Process) string {
	color := tui.FOREWHITE
	if modifyColor, ok := knownProcs[proc.Executable]; ok {
		color = modifyColor
	}
	if ctx.CurrentAgent.Name != "" && proc.Pid == ctx.CurrentAgent.PID {
		color = tui.GREEN
	}
	fmt.Fprintf(table, "%d\t%d\t%s\t%s\t\n", proc.Pid, proc.Ppid, proc.Executable, proc.Owner)
	return color
}

func procdump(args []string, ctx ShellContext, rpc RPCServer) {

	opts := procFilters(args)

	var pid int
	spid, found := opts["pid"]
	if found {
		pid, _ = strconv.Atoi((spid.(string)))
	} else {
		pid = -1
	}

	var name string
	proc, found := opts["proc"]
	if found {
		name = proc.(string)
	}

	var timeout int
	stimeout, found := opts["timeout"]
	if found {
		timeout, _ = strconv.Atoi(stimeout.(string))
	} else {
		timeout = 360
	}

	if pid == -1 && name != "" {
		pid = getPIDByName(name, ctx, rpc)
	}
	if pid == -1 {
		fmt.Printf(Warn + "Invalid process target\n")
		return
	}

	if timeout < 1 {
		fmt.Printf(Warn + "Invalid timeout argument\n")
		return
	}

	ctrl := make(chan bool)
	go spin.Until("Dumping remote process memory ...", ctrl)
	data, _ := proto.Marshal(&ghostpb.ProcessDumpReq{
		GhostID: ctx.CurrentAgent.ID,
		Pid:     int32(pid),
		Timeout: int32(timeout),
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgProcessDumpReq,
		Data: data,
	}, defaultTimeout)
	ctrl <- true
	<-ctrl

	procDump := &ghostpb.ProcessDump{}
	proto.Unmarshal(resp.Data, procDump)
	if procDump.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}

	hostname := ctx.CurrentAgent.Hostname
	temp := path.Base(fmt.Sprintf("procdump_%s_%d_*", hostname, pid))
	f, err := ioutil.TempFile("", temp)
	if err != nil {
		fmt.Printf(Warn+"Error creating temporary file: %v\n", err)
	}
	f.Write(procDump.GetData())
	fmt.Printf(Success+"Process dump stored in %s\n", f.Name())
}

func terminate(args []string, ctx ShellContext, rpc RPCServer) {

	opts := procFilters(args)
	if len(args) != 1 {
		fmt.Printf(Warn + "Please provide a PID\n")
		return
	}
	var pidStr string
	spid, found := opts["pid"]
	if found {
		pidStr = spid.(string)
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Printf(Warn+"Error: %v\n", err)
		return
	}
	data, _ := proto.Marshal(&ghostpb.TerminateReq{
		GhostID: ctx.CurrentAgent.ID,
		Pid:     int32(pid),
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgTerminate,
		Data: data,
	}, defaultTimeout)
	termResp := &ghostpb.Terminate{}
	err = proto.Unmarshal(resp.Data, termResp)
	if err != nil {
		fmt.Printf(Warn+"Error: %v\n", err)
		return
	}
	if termResp.Err != "" {
		fmt.Printf(Warn+"Error: %s\n", termResp.Err)
		return
	}
	fmt.Printf(Info+"Process %d has been terminated\n", pid)
}

func getPIDByName(name string, ctx ShellContext, rpc RPCServer) int {
	data, _ := proto.Marshal(&ghostpb.PsReq{GhostID: ctx.CurrentAgent.ID})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgPsReq,
		Data: data,
	}, defaultTimeout)
	ps := &ghostpb.Ps{}
	proto.Unmarshal(resp.Data, ps)
	for _, proc := range ps.Processes {
		if proc.Executable == name {
			return int(proc.Pid)
		}
	}
	return -1
}

func migrate(args []string, ctx ShellContext, rpc RPCServer) {

	if len(args) != 1 {
		fmt.Printf(Warn + "You must provide a PID to migrate to")
		return
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf(Warn+"Error: %v", err)
	}
	config := getActiveGhostConfig(ctx)
	ctrl := make(chan bool)
	msg := fmt.Sprintf("Migrating into %d ...", pid)
	go spin.Until(msg, ctrl)
	data, _ := proto.Marshal(&clientpb.MigrateReq{
		Pid:     uint32(pid),
		Config:  config,
		GhostID: ctx.CurrentAgent.ID,
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgMigrate,
		Data: data,
	}, 45*time.Minute)
	ctrl <- true
	<-ctrl
	if resp.Err != "" {
		fmt.Printf(Warn+"%s\n", resp.Err)
	} else {
		fmt.Printf("\n"+Success+"Successfully migrated to %d\n", pid)
	}
}

func procFilters(args []string) (opts map[string]interface{}) {
	opts = make(map[string]interface{}, 0)

	for _, arg := range args {

		// Process type
		if strings.Contains(arg, "exe") {
			vals := strings.Split(arg, "=")
			opts["exe"] = vals[1]
		}
		// Process
		if strings.Contains(arg, "proc") {
			vals := strings.Split(arg, "=")
			opts["proc"] = vals[1]
		}
		// Owner
		if strings.Contains(arg, "owner") {
			vals := strings.Split(arg, "=")
			opts["owner"] = vals[1]
		}
		// Process ID
		if strings.Contains(arg, "pid") {
			vals := strings.Split(arg, "=")
			timeout, _ := strconv.Atoi(vals[1])
			opts["pid"] = timeout
		}
		// Timeout
		if strings.Contains(arg, "timeout") {
			vals := strings.Split(arg, "=")
			timeout, _ := strconv.Atoi(vals[1])
			opts["timeout"] = timeout
		}

		// Special Arguments
		if strings.Contains(arg, "args") {
			desc := regexp.MustCompile(`\b(args){1}.*"`)
			result := desc.FindStringSubmatch(strings.Join(args, " "))
			opts["args"] = strings.Trim(strings.TrimPrefix(result[0], "args="), "\"")
		}
	}

	return opts

}

var (
	// Stylizes known processes in the `ps` command
	knownProcs = map[string]string{
		"ccSvcHst.exe": tui.RED, // SEP
		"cb.exe":       tui.RED, // Carbon Black
	}
)
