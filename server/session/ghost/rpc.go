package ghost

import (
	"io"
	"sync"
	"time"

	transportpb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport"
	"github.com/sirupsen/logrus"
)

// The RPC layer should be applicable independently and concurrently for all
// implant channels, therefore the code should preferably wrap an io.ReadWriteCloser
// and all I/O should be handled accordingly.

// RPC - Contains all infrastructure needed for sending/receiving RPC requests/responses
// from this implant (or one of its channels).
type RPC struct {
	Send      chan *transportpb.Envelope            // All outbound (requests to implant)
	Resp      map[uint64]chan *transportpb.Envelope // All inbound (implant responses)
	RespMutex *sync.RWMutex                         // Concurrency
	log       *logrus.Entry                         // Logger used only for debug purposes
}

// SetupStreamRPC - An implant has spawned a new channel (or setting up its main one), and wants to be
// able to send/receive RPC requests/responses to implant.
// NOTE: The implant has already registered, therefore all RPC boilerplate code should only manage
// normal requests, including kill/sleep ones. This boilerplate code has already been added to the main
// channel, which MUST be a RPC-enabled one.
func (c *Client) SetupStreamRPC(stream io.ReadWriteCloser) (err error) {

	// Log that we are upgrading a connection a to ghost session

	// Defer cleanup if needed

	//
	return
}

// Request - All RPC requests to the implant remote side are made through this function.
func (r *RPC) Request(msgType int32, timeout time.Duration, req []byte) (res []byte, err error) {

	// Setup envelope and concurrency safety

	// Defer setup/cleanup/security

	// Wait for response concurrently and return

	return
}
