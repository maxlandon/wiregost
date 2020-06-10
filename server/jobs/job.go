package jobs

import (
	"sync"

	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

// Job - A background job running in Wiregost
type Job struct {
	Proto   *serverpb.Job // Protobuf serialization
	JobCtrl chan bool     // Asynchronous job control
}

var (
	// Jobs - All background tasks, (compilers, proxies, etc) running on the Wiregost server
	Jobs = &jobs{
		Active: &map[uint32]*Job{},
		mutex:  &sync.RWMutex{},
	}

	jobID = new(uint32)
)

// jobs - Holds references to all active jobs
type jobs struct {
	Active *map[uint32]*Job
	mutex  *sync.RWMutex
}

// AddJob - Add a job to the job group
func (j *jobs) AddJob(job *Job) {
}

// RemoveJob - Remove a job from the job group
func (j *jobs) RemoveJob(jobID uint32) {
}

// Job - Get a job by ID
func (j *jobs) Job(jobID uint32) (job *Job) {
	return
}

// GetJobID - Generate an ID for a new job
func GetJobID() (id uint32) {
	return
}
