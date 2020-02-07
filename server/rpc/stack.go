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
		fmt.Println("error at unmarshalling")
		resp(data, err)
	}

	// Find module
	path := strings.Join(stackReq.Path, "/")
	mod, err := module.GetModule(path)
	if err != nil {
		fmt.Println("Error finding module")
		resp(data, err)
	}

	fmt.Println(mod)
	// Load it on the workspace stack
	wsID := uint(stackReq.WorkspaceID)
	stack := (*module.Stacks)[wsID]
	stack.LoadModule(path)
	fmt.Println(mod)
	// (*stack.Loaded)[path] = mod
	// testmod := stack.Loaded[path]

	// Send back module
	stackUse := &clientpb.Stack{
		Path: stackReq.Path,
		// Modules: []*clientpb.Module{mod.ToProtobuf()},
		Modules: []*clientpb.Module{mod.ToProtobuf()},
		Err:     "",
	}

	data, err = proto.Marshal(stackUse)
	resp(data, err)
}

func rpcStackPop(data []byte, timeout time.Duration, resp RPCResponse) {
	stackReq := &clientpb.StackReq{}
	err := proto.Unmarshal(data, stackReq)
	if err != nil {
		resp(data, err)
	}

	// data, err := proto.Marshal()

	resp(data, err)
}

func rpcStackList(data []byte, timeout time.Duration, resp RPCResponse) {
	stackReq := &clientpb.StackReq{}
	err := proto.Unmarshal(data, stackReq)
	if err != nil {
		resp(data, err)
	}

	// data, err := proto.Marshal()

	resp(data, err)
}
