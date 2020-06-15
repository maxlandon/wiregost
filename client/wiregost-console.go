package main

import (
	client "github.com/maxlandon/wiregost/client/console"
)

// Console executable entry
func main() {

	// Start the client console. Configuration loading, authentication, and connection info fetching is handled by this function.
	client.Console.Start()
}
