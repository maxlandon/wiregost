package channels

import (
	"io"
)

// Channel - A goroutine object that is stored and manageable by the implant user.
// This is simply boilerplate to provide more explicit control on implant concurrent processes,
// such as open shells, programs running, etc...
// It is somehow the equivalent of jobs on the Wiregost server.
// This is also an attempt at mimicking the channels of Metasploit's Meterpreters, in which
// many different processes can be run concurrently, such as system shells.
// This can be used also when running concurrent processes such as routers: if we switch to a router's
// channels, we might get a shell refreshed with currently routed connections (like gost)
type Channel struct {
	ID          uint32    // ID of the Channel
	Process     string    // Process name string (generally from ps)
	Description string    // Further description for the process
	Err         string    // Potential errors caught
	JobCtrl     chan bool // Asynchronous job control
	Stdout      io.ReadCloser
	Stdin       io.WriteCloser
}

// SetupChannels - Inits the goroutine management system of the implant.
func SetupChannels() {

}
