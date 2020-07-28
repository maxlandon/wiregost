package route

// InitRouting - Setups and starts all routing infrastructure for this implant.
// The function has access to many objects previously setup, and the behaviour
// of the function may vary depending on security, authorisations, compiled routes, etc...
func InitRouting() {

	// Check pre-compiled authorisations

	// If needed, open a dedicated muxed stream over which we send routing requests/responses.
	// Start a concurrent handler, for managing requests.

	// Check if we have various route listeners to open.

}
