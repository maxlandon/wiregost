package generic

import (
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	"github.com/maxlandon/wiregost/server/c2"
)

// Client - All generic RPC handlers.
type Client struct {
	Session *c2.Session
}

// NewClient - Once a ghost transport layer (Session) is set up, we register it to a new RPC client.
func NewClient(sess *c2.Session) (cli *Client) {
	cli = &Client{
		Session: sess,
	}
	return
}

// C2Custom - The RPC stubs used are custom made (DNS/MTLS/HTTPS)
func C2Custom(ghost ghostpb.Ghost) bool {

	// Here, we check the current transport type of the ghost implant

	return false
}

// C2RPC - The RPC stubs are RPCX stubs (KCP, QUIC, HTTP, TCP)
func C2RPC(ghost ghostpb.Ghost) bool {

	return false
}
