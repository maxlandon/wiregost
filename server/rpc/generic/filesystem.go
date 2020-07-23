package generic

import (
	"context"

	corepb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost/core"
	wContext "github.com/maxlandon/wiregost/server/context"
)

// Ls - Send a Ls request to implant
func (c *Client) Ls(ctx context.Context, req corepb.LsRequest) (res corepb.Ls, err error) {

	// Check permissions
	in := wContext.GetMetadata(ctx)

	// Check if custom DNS/MTLS/HTTP
	if C2Custom(in.Ghost) {

		return
	}

	// Else, call through RPCX/gRPC
	if C2RPC(in.Ghost) {

		return
	}

	return
}
