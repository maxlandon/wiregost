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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/c2"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/db"
)

func rpcJobs(_ []byte, timeout time.Duration, resp Response) {
	jobs := &clientpb.Jobs{
		Active: []*clientpb.Job{},
	}
	for _, job := range *core.Jobs.Active {
		jobs.Active = append(jobs.Active, &clientpb.Job{
			ID:          int32(job.ID),
			Name:        job.Name,
			Description: job.Description,
			Protocol:    job.Protocol,
			Port:        int32(job.Port),
		})
	}
	data, err := proto.Marshal(jobs)
	if err != nil {
		rpcLog.Errorf("Error encoding rpc response %v", err)
		resp([]byte{}, err)
		return
	}
	resp(data, err)
}

func rpcJobKill(data []byte, timeout time.Duration, resp Response) {
	jobKillReq := &clientpb.JobKillReq{}
	err := proto.Unmarshal(data, jobKillReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	job := core.Jobs.Job(int(jobKillReq.ID))
	jobKill := &clientpb.JobKill{ID: int32(job.ID)}
	// kill job
	if job != nil {
		job.JobCtrl <- true
		jobKill.Success = true
	} else {
		jobKill.Success = false
		jobKill.Err = "Invalid Job ID"
	}

	// If persistent listener, delete config
	var persist = fmt.Sprintf("%s[P]%s ", tui.GREEN, tui.RESET)
	if job.Protocol != "" && strings.HasPrefix(job.Description, persist) {
		bucket, _ := db.GetBucket(c2.ListenerBucketName)
		ls, _ := bucket.List(c2.ListenerNamespace)

		// listeners := []*c2.ListenerConfig{}
		for _, listener := range ls {
			rawListener, err := bucket.Get(listener)
			if err != nil {
				fmt.Println(err)
			}
			config := &c2.ListenerConfig{}
			err = json.Unmarshal(rawListener, config)
			if err != nil {
				fmt.Println(err)
			}
			if config.LPort == job.Port && config.Description == job.Description && config.Name == job.Name {
				bucket.Delete(listener)
			}
		}
	}

	data, err = proto.Marshal(jobKill)
	resp(data, err)
}
