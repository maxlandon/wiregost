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
	"github.com/maxlandon/wiregost/server/db"
	"github.com/maxlandon/wiregost/server/generate"
)

func rpcListProfiles(data []byte, timeout time.Duration, resp RPCResponse) {
	profiles := &clientpb.Profiles{List: []*clientpb.Profile{}}
	for name, config := range generate.Profiles() {
		profiles.List = append(profiles.List, &clientpb.Profile{
			Name:   name,
			Config: config.ToProtobuf(),
		})
	}
	data, err := proto.Marshal(profiles)
	resp(data, err)
}

func rpcDeleteProfile(data []byte, timeout time.Duration, resp RPCResponse) {

	profileReq := &clientpb.Profile{}
	err := proto.Unmarshal(data, profileReq)
	if err != nil {
		resp(data, err)
	}

	bucket, err := db.GetBucket(generate.ProfilesBucketName)
	if err != nil {
		profileReq.Name = err.Error()
		data, err = proto.Marshal(profileReq)
		resp(data, err)
	}

	err = bucket.Delete(profileReq.Name)
	if err != nil {
		profileReq.Name = err.Error()
		data, err = proto.Marshal(profileReq)
		resp(data, err)
	} else {
		data, err = proto.Marshal(profileReq)
		resp(data, err)
	}
}
