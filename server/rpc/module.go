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
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/module"
)

func rpcModuleSetOption(data []byte, timeout time.Duration, resp RPCResponse) {
	optionReq := &clientpb.SetOptionReq{}
	err := proto.Unmarshal(data, optionReq)
	if err != nil {
		resp(data, err)
	}

	// Find module
	path := strings.Join(optionReq.Path, "/")
	wsID := uint(optionReq.WorkspaceID)
	stack := (*module.Stacks)[wsID]
	mod := (*stack.Loaded)[path]
	mod.SetOption(optionReq.Name, optionReq.Value)

	option := &clientpb.SetOption{
		Success: true,
		Err:     "",
	}
	data, err = proto.Marshal(option)
	resp(data, err)
}

func rpcModuleRun(data []byte, timeout time.Duration, resp RPCResponse) {
	modReq := &clientpb.ModuleActionReq{}
	err := proto.Unmarshal(data, modReq)
	if err != nil {
		resp(data, err)
	}

	// Find module
	path := strings.Join(modReq.Path, "/")
	wsID := uint(modReq.WorkspaceID)
	stack := (*module.Stacks)[wsID]
	mod := (*stack.Loaded)[path]

	res, err := mod.Run(modReq.Action)

	modRun := &clientpb.ModuleAction{}
	if err != nil {
		modRun.Sucess = false
		modRun.Err = err.Error()
	} else {
		modRun.Sucess = true
		modRun.Result = res
	}

	data, err = proto.Marshal(modRun)
	resp(data, err)
}
