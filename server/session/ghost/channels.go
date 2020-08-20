package ghost

import (
	"bytes"
	"io"

	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/sirupsen/logrus"
)

// Channel - The client-side peer of a concurrent process running on the Ghost implant.
// This process might be any command, synchronous or not, or any core functionality of
// the implant that needs/allows concurrent running, such as routers, shells, etc.
// Some of these channels MAY NOT be used/interacted with by users: they are simply "jobs"
// for routing and listeners.
type Channel struct {
	ID         string               // A UUID-as-string for this channel.
	Name       string               // An optional name for this channel ("shellTest", "routingWork", etc)
	Type       serverpb.ChannelType // This channel supports RPC or byte streams only.
	WorkingDir string               // The channel has its own primary context

	// stream - A channel stream is generally a logical connection that has been
	// "muxed" from the implant's physical connection. Triggering the opening of
	// a new channel with its dedicated stream is always done through the Client's
	// main Interactive stream.
	// For things like pseudo-command shells, this stream is used in an asynchronous
	// way: it pushes output as it comes, and does not wait each time for all the
	// output to go out first before sending it back to the server.
	// This stream might also implement the RPC infrastructure needed to perform
	// core actions through Protobuf requests.
	stream io.ReadWriteCloser

	// Log - Each channel may optionally have a dedicated logger, for things like
	// command history.
	Log *logrus.Entry

	// History - Every channel has its own command history, which might be optionally
	// persisted across implant reboots. This field might be used for console history
	// completion per-channel, for instance.
	History []string

	// Buffer - This buffer is used to store output that might have been procuced while
	// the user is away from the Session/channel, and that we want to be able to access
	// it later.
	Buffer bytes.Buffer
}

// NewUserChannel - A user has requested to open a new concurrent "sub-session" in the ghost
// implant. A request has been forwarded to the remote peer, and the Session has muxed
// the physical/logical connection. This channel's only purpose is to be operated by the user,
// and not for being used as a routing/transport mechanism.
func (s *Client) NewUserChannel(chanType int) (ch *Channel, err error) {

	// Instantiation & data (ID, name, etc)

	// Send Request to implant, wait acknowledge.
	// This means the ghost implant has called Accept(), blocking for him.

	// Stream multiplexing: we mux the main stream, add to chan object

	// Add channel to Session, and potentially register things needed.

	return
}

// GetChannelShell - This function is used to wire a ghost implant channel (its client-side object)
// with a user console. It should provide two -real Go- channels: 1 for reading user commands, and
// one for pushing command output back to the console.
//
// Usually, this function is used when the user opens a new channel with Shell type,
// which allows for pretty much any command (string) and returns only strings. Therefore,
// this function only wires byte streams to others: no RPC boiler code is needed here.
func (s *Client) GetChannelShell(chanID string) {

	// NOTE: If the channel is already up and its type is CORE, we "downgrade" the stream
	// for handling only byte streams and []string (for commands sent)
}

// GetChannelCore - This function is the same as GetChannelShell(), but insted add some RPC boilerplate.
func (s *Client) GetChannelCore(chanID string) {

	// NOTE: If the channel is already up and its type is SHELL, we "upgrade" the stream
	// for handling RPC calls.
}
