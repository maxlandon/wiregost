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

package commands

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/maxlandon/wiregost/client/util"
	. "github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/olekukonko/tablewriter"
)

func RegisterJobCommands() {

	jobs := &Command{
		Name: "jobs",
		SubCommands: []string{
			"kill",
			"kill-all",
		},
		Handle: func(r *Request) error {
			if len(r.Args) == 0 {
				fmt.Println()
				listJobs(*r.context, r.context.Server.RPC)
			} else {
				switch r.Args[0] {
				case "kill":
					if len(r.Args) == 1 {
						fmt.Printf("\n" + Error + "Provide one or more Job IDs\n")
						return nil
					} else {
						for _, arg := range r.Args[1:] {
							idInt, _ := strconv.Atoi(arg)
							id := int32(idInt)
							killJob(id, r.context.Server.RPC)
						}
					}
				case "kill-all":
					killAllJobs(r.context.Server.RPC)
				}
			}
			return nil
		},
	}

	AddCommand("main", jobs)
	AddCommand("module", jobs)
}

func listJobs(ctx ShellContext, rpc RPCServer) {

	jobs := GetJobs(rpc)
	if jobs == nil {
		return
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
}

// GetJobs - Exported so that shell can use it when refreshing
func GetJobs(rpc RPCServer) *clientpb.Jobs {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgJobs,
		Data: []byte{},
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return nil
	}
	jobs := &clientpb.Jobs{}
	proto.Unmarshal(resp.Data, jobs)
	return jobs
}

func killAllJobs(rpc RPCServer) {
	jobs := GetJobs(rpc)
	if jobs == nil {
		return
	}
	for _, job := range jobs.Active {
		killJob(job.ID, rpc)
	}
}

func killJob(jobID int32, rpc RPCServer) {

	fmt.Printf("\n"+Info+"Killing job #%d ...\n", jobID)
	data, _ := proto.Marshal(&clientpb.JobKillReq{ID: jobID})
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgJobKill,
		Data: data,
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError, "%s\n", resp.Err)
		return
	}
	jobKill := &clientpb.JobKill{}
	proto.Unmarshal(resp.Data, jobKill)

	if jobKill.Success {
		fmt.Printf(Success+"Successfully killed job #%d\n", jobKill.ID)
	} else {
		fmt.Printf(Error+"Failed to kill job #%d, %s\n", jobKill.ID, jobKill.Err)
	}
}

func printJobs(jobs map[int32]*clientpb.Job) {

	tab := util.Table()
	tab.SetHeader([]string{"ID", "Name", "Protocol", "Port", "Description"})
	tab.SetColWidth(50)
	tab.SetColMinWidth(4, 30)
	tab.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)
	var keys []int
	for _, job := range jobs {
		keys = append(keys, int(job.ID))
	}
	sort.Ints(keys) // Fucking Go can't sort int32's, so we convert to/from int's

	for _, k := range keys {
		job := jobs[int32(k)]
		description := util.Wrap(job.Description, 100)
		tab.Append([]string{strconv.Itoa(int(job.ID)), job.Name, job.Protocol, strconv.Itoa(int(job.Port)), description})
		// tab.Append([]string{strconv.Itoa(int(job.ID)), job.Name, job.Protocol, strconv.Itoa(int(job.Port)), job.Description})
	}

	tab.Render()
}
