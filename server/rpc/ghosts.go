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

	"github.com/gogo/protobuf/proto"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/generate"
)

func rpcGhostBuilds(_ []byte, timeout time.Duration, resp Response) {
	configs, err := generate.GhostConfigMap()
	if err != nil {
		resp([]byte{}, err)
		return
	}
	ghostBuilds := &clientpb.GhostBuilds{
		Configs: map[string]*clientpb.GhostConfig{},
	}
	for name, cfg := range configs {
		ghostBuilds.Configs[name] = cfg.ToProtobuf()
	}
	data, err := proto.Marshal(ghostBuilds)
	resp(data, err)
}

func rpcListCanaries(_ []byte, timeout time.Duration, resp Response) {
	jsonCanaries, err := generate.ListCanaries()
	if err != nil {
		resp([]byte{}, err)
	}
	rpcLog.Infof("Found %d canaries", len(jsonCanaries))
	canaries := []*clientpb.DNSCanary{}
	for _, canary := range jsonCanaries {
		canaries = append(canaries, canary.ToProtobuf())
	}
	data, err := proto.Marshal(&clientpb.Canaries{
		Canaries: canaries,
	})
	resp(data, err)
}
