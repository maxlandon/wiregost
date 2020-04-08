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

package jobs

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// JobsKillCmd - kill one or all background jobs
type JobsKillCmd struct {
	Positional struct {
		JobID int `description:"Job ID" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var JobsKill JobsKillCmd

func RegisterJobsKill() {
	jobs := CommandParser.Find(constants.Jobs)
	jobs.AddCommand(constants.JobsKill, "", "", &JobsKill)

	kill := jobs.Find(constants.JobsKill)
	kill.ShortDescription = "Kill an active background job"
}

// Execute - kill one or all background jobs
func (jk *JobsKillCmd) Execute(args []string) error {
	return killJob(int32(jk.Positional.JobID), Context.Server.RPC)
}

// JobsKillAllCmd - kill all active background jobs
type JobsKillAllCmd struct{}

var JobsKillAll JobsKillAllCmd

func RegisterJobsKillAll() {
	jobs := CommandParser.Find(constants.Jobs)
	jobs.AddCommand(constants.JobsKillAll, "", "", &JobsKillAll)

	all := jobs.Find(constants.JobsKillAll)
	all.ShortDescription = "Kill all active background jobs"
}

// Execute - kill all active background jobs
func (ka *JobsKillAllCmd) Execute(args []string) error {

	jobs := GetJobs(Context.Server.RPC)
	if jobs == nil {
		return nil
	}
	for _, job := range jobs.Active {
		killJob(job.ID, Context.Server.RPC)
	}

	return nil
}

func killJob(jobID int32, rpc RPCServer) error {

	fmt.Printf("\n"+Info+"Killing job #%d ...\n", jobID)
	data, _ := proto.Marshal(&clientpb.JobKillReq{ID: jobID})
	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgJobKill,
		Data: data,
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError, "%s\n", resp.Err)
		return nil
	}
	jobKill := &clientpb.JobKill{}
	proto.Unmarshal(resp.Data, jobKill)

	if jobKill.Success {
		fmt.Printf(Success+"Successfully killed job #%d\n", jobKill.ID)
	} else {
		fmt.Printf(Error+"Failed to kill job #%d, %s\n", jobKill.ID, jobKill.Err)
	}

	return nil
}
