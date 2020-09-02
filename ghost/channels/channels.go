package channels

import (
	"io"

	"github.com/sirupsen/logrus"

	corepb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost/core"
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
	ID          uint32             // ID of the Channel
	Name        string             // A name for this channel
	Type        corepb.ChannelType // This channel supports RPC or byte streams only.
	Process     string             // Process name string (generally from ps)
	WorkingDir  string             // Channel current working directory (abstract)
	IsUsed      bool               // Equivalent value for the client-side
	Description string             // Further description for the process
	Err         string             // Potential errors caught
	JobCtrl     chan bool          // Asynchronous job control
	// Stdout      io.ReadCloser
	// Stdin       io.WriteCloser
	Stream io.ReadWriteCloser // Stdin, stdout & stderr channels
	Log    *logrus.Entry      // A logger for this channel
}

// SetupChannels - Inits the goroutine management system of the implant.
func SetupChannels() {

	// We init the channel map

	// Instantiate a new channel, which will be the implant's main, RPC enabled.

	// We don't pass it a stream yet, because we did not reach back to the server.
}

// StartMainChannel - Init the channel that will be used for all core functionality from implant
// startup to kill/sleep. We might later use other functions for starting specialized channels.
func StartMainChannel(stream io.ReadWriteCloser) (main *Channel, err error) {

	main = &Channel{
		Stream: stream,
	}

	return
}
