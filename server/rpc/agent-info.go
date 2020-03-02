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
	"time"

	"github.com/golang/protobuf/proto"

	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
)

func rpcIfconfig(req []byte, timeout time.Duration, resp Response) {
	ifconfigReq := &ghostpb.IfconfigReq{}
	err := proto.Unmarshal(req, ifconfigReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	sliver := (*core.Wire.Ghosts)[ifconfigReq.GhostID]
	if sliver == nil {
		resp([]byte{}, err)
		return
	}

	data, _ := proto.Marshal(&ghostpb.IfconfigReq{})
	data, err = sliver.Request(ghostpb.MsgIfconfigReq, timeout, data)
	resp(data, err)
}
