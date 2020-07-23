package generic

import (
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
)

// Client - All generic RPC handlers.
type Client struct {
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
