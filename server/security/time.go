package security

import "time"

// TimeScheduler - Manages the implant connection schedule:
// It can restrict/start/stop implant communications based on these constraints
//
// Because one listener (or therefore, service) can handle multiple ghost implants,
// each implant object has its own scheduler, and the server, by querying it, can
// stop communications with one implant without stopping the others.
type TimeScheduler struct {
	Deadline time.Duration
	Start    time.Time
	Stop     time.Time
}
