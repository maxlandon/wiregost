package main

import (
	"github.com/maxlandon/wiregost/client/assets"
	client "github.com/maxlandon/wiregost/client/console"
)

// Console executable entry
func main() {

	// Load server connection configuration (check files in ~/.wiregost first, then binary)
	assets.LoadServerConfig()

	// Start the client console
	client.Console.Start()
}
