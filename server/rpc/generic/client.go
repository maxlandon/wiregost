package generic

import (
	corepb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost/core"
	"github.com/maxlandon/wiregost/server/context"
)

// Client - All generic RPC handlers.
type Client struct {
}

// Ls - Send a Ls request to implant
func (c *Client) Ls(ctx context.RPCContext, req corepb.LsRequest) (err error) {

	// Check permissions

	// Check if custom DNS/MTLS/HTTP

	//
	return
}
