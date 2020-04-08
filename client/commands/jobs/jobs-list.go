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
	"sort"
	"strconv"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

// JobsCmd - List active background jobs
type JobsCmd struct{}

var Jobs JobsCmd

func RegisterJobs() {
	CommandParser.AddCommand(constants.Jobs, "", "", &Jobs)

	jobs := CommandParser.Find(constants.Jobs)
	CommandMap[MAIN_CONTEXT] = append(CommandMap[MAIN_CONTEXT], jobs)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], jobs)
	jobs.ShortDescription = "List active background jobs"
	jobs.SubcommandsOptional = true
}

// Execute - List active background jobs
func (j *JobsCmd) Execute(args []string) error {

	jobs := GetJobs(Context.Server.RPC)
	if jobs == nil {
		return nil
	}
	activeJobs := map[int32]*clientpb.Job{}
	for _, job := range jobs.Active {
		activeJobs[job.ID] = job
	}
	if 0 < len(activeJobs) {
		printJobs(activeJobs)
	} else {
		fmt.Printf(Info + "No active jobs\n")
	}

	return nil
}

// GetJobs - Exported so that shell can use it when refreshing
func GetJobs(rpc RPCServer) *clientpb.Jobs {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgJobs,
		Data: []byte{},
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}
	jobs := &clientpb.Jobs{}
	proto.Unmarshal(resp.Data, jobs)
	return jobs
}

func printJobs(jobs map[int32]*clientpb.Job) {

	tab := util.NewTable()
	headers := []string{"ID", "Name", "Protocol", "Port", "Description"}
	widths := []int{3, 10, 10, 5, 50}
	tab.SetColumns(headers, widths)
	tab.SetColWidth(50)

	var keys []int
	for _, job := range jobs {
		keys = append(keys, int(job.ID))
	}
	sort.Ints(keys) // Fucking Go can't sort int32's, so we convert to/from int's

	for _, k := range keys {
		job := jobs[int32(k)]
		description := util.AutoWrap(job.Description)
		tab.Append([]string{strconv.Itoa(int(job.ID)), job.Name, job.Protocol, strconv.Itoa(int(job.Port)), description})
	}

	tab.Output()
}
