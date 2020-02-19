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

package core

import (
	"sync"

	pb "github.com/maxlandon/wiregost/protobuf/client"
)

var (
	Jobs = &jobs{
		Active: &map[int]*Job{},
		mutex:  &sync.RWMutex{},
	}

	jobID = new(int)
)

// Job - Is a background Job object
type Job struct {
	ID          int
	Name        string
	Description string
	Err         string
	Protocol    string
	Port        uint16
	JobCtrl     chan bool
}

// ToProtobuf - Returns the protobuf version of the object
func (j *Job) ToProtobuf() *pb.Job {
	return &pb.Job{
		ID:          int32(j.ID),
		Name:        j.Name,
		Description: j.Description,
		Protocol:    j.Protocol,
		Port:        int32(j.Port),
	}
}

// Holds refs to all active jobs
type jobs struct {
	Active *map[int]*Job
	mutex  *sync.RWMutex
}

// AddJob - Add a job to the job group
func (j *jobs) AddJob(job *Job) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	(*j.Active)[job.ID] = job
}

// RemoveJob - Remove a job to the job group
func (j *jobs) RemoveJob(job *Job) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	delete((*j.Active), job.ID)
}

// Job - Get Job by ID
func (j *jobs) Job(jobID int) *Job {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	return (*j.Active)[jobID]
}

func GetJobID() int {
	newID := (*jobID) + 1
	(*jobID)++
	return newID
}
