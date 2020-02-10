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
	"github.com/maxlandon/wiregost/server/core"
)

func rpcJobs(_ []byte, timeout time.Duration, resp RPCResponse) {
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

func rpcJobKill(data []byte, timeout time.Duration, resp RPCResponse) {
	jobKillReq := &clientpb.JobKillReq{}
	err := proto.Unmarshal(data, jobKillReq)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	job := core.Jobs.Job(int(jobKillReq.ID))
	jobKill := &clientpb.JobKill{ID: int32(job.ID)}
	if job != nil {
		job.JobCtrl <- true
		jobKill.Success = true
	} else {
		jobKill.Success = false
		jobKill.Err = "Invalid Job ID"
	}
	data, err = proto.Marshal(jobKill)
	resp(data, err)
}
