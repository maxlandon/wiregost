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

package rpc

import (
	"fmt"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/module"
)

func rpcStackUse(data []byte, timeout time.Duration, resp RPCResponse) {
	stackReq := &clientpb.StackReq{}
	err := proto.Unmarshal(data, stackReq)
	if err != nil {
		resp(data, err)
	}

	// Find the workspace stack
	wsID := uint(stackReq.WorkspaceID)
	stack := (*module.Stacks)[wsID][stackReq.User]

	path := strings.Join(stackReq.Path, "/")

	if mod, found := (*stack.Loaded)[path]; !found {
		// If module not found, load it and send it
		err = stack.LoadModule(path)
		if err != nil {
			stackErr := &clientpb.Stack{
				Path: stackReq.Path,
				Err:  err.Error(),
			}

			data, err = proto.Marshal(stackErr)
			resp(data, err)
			return
		}
		mod = (*stack.Loaded)[path]
		module := []*clientpb.Module{mod.ToProtobuf()}
		stackUse := &clientpb.Stack{
			Path:    stackReq.Path,
			Modules: module,
			Err:     "",
		}

		data, err = proto.Marshal(stackUse)
		resp(data, err)
		// If found, send it
	} else {
		mod = (*stack.Loaded)[path]
		module := []*clientpb.Module{mod.ToProtobuf()}
		stackUse := &clientpb.Stack{
			Path:    stackReq.Path,
			Modules: module,
			Err:     "",
		}

		data, err = proto.Marshal(stackUse)
		resp(data, err)
	}
}

func rpcStackPop(data []byte, timeout time.Duration, resp RPCResponse) {
	stackReq := &clientpb.StackReq{}
	err := proto.Unmarshal(data, stackReq)
	if err != nil {
		resp(data, err)
	}

	// Find workspace stack
	wsID := uint(stackReq.WorkspaceID)
	stack := (*module.Stacks)[wsID][stackReq.User]

	if stackReq.All {
		for k, _ := range *stack.Loaded {
			stack.PopModule(k)
		}
	} else {
		stack.PopModule(strings.Join(stackReq.Path, "/"))
	}

	// Get next module on top of stack
	modules := []*clientpb.Module{}
	for _, v := range *stack.Loaded {
		modules = append(modules, v.ToProtobuf())
	}
	next := []string{}
	if len(modules) != 0 {
		next = modules[0].Path
		modules = []*clientpb.Module{modules[0]}
	}

	// Send it back
	stackPop := &clientpb.Stack{
		Path:    next,
		Modules: modules,
		Err:     "",
	}

	data, err = proto.Marshal(stackPop)
	resp(data, err)
}

func rpcStackList(data []byte, timeout time.Duration, resp RPCResponse) {
	stackReq := &clientpb.StackReq{}
	err := proto.Unmarshal(data, stackReq)
	if err != nil {
		resp(data, err)
	}

	// Find workspace stack
	wsID := uint(stackReq.WorkspaceID)
	stack := (*module.Stacks)[wsID][stackReq.User]
	for _, s := range (*module.Stacks)[wsID] {
		fmt.Println(s)
	}
	fmt.Println(stack)

	modules := []*clientpb.Module{}

	for _, v := range *stack.Loaded {
		modules = append(modules, v.ToProtobuf())
	}

	stackList := &clientpb.Stack{
		Path:    stackReq.Path,
		Modules: modules,
		Err:     "",
	}

	data, err = proto.Marshal(stackList)
	resp(data, err)
}
