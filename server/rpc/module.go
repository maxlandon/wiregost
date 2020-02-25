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
	"github.com/maxlandon/wiregost/server/core"
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
	stack := (*module.Stacks)[wsID][optionReq.User]
	mod := (*stack.Loaded)[path]
	mod.SetOption(optionReq.Name, optionReq.Value)

	option := &clientpb.SetOption{
		Success: true,
		Err:     "",
	}
	data, err = proto.Marshal(option)
	resp(data, err)
	time.Sleep(time.Millisecond * 50)

	optionUpdated := fmt.Sprintf("%s %s %s", optionReq.Name, optionReq.Value, path)
	optionBytes := []byte(optionUpdated)

	for _, c := range *core.Clients.Connections {
		if c.User == optionReq.User {
			core.EventBroker.Publish(core.Event{
				Client:    c,
				EventType: "module",
				Data:      optionBytes,
			})
			// One push is enough
			break
		}
	}
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
	stack := (*module.Stacks)[wsID][modReq.User]
	mod := (*stack.Loaded)[path]

	var res string
	var modErr error
	if modReq.Profile != "" {
		action := modReq.Action + " " + modReq.Profile
		res, modErr = mod.Run(action)
	} else {
		res, modErr = mod.Run(modReq.Action)
	}

	modRun := &clientpb.ModuleAction{}
	if modErr != nil {
		modRun.Success = false
		modRun.Err = modErr.Error()
	} else {
		modRun.Success = true
		modRun.Result = res
	}

	// Send updated module
	modRun.Updated = mod.ToProtobuf()

	data, err = proto.Marshal(modRun)
	resp(data, err)
}
