package c2

// HandleGhostRegistration - Of all the process starting from TCP handshake to
// complete registration and usage of the ghost implant by users, this function
// is the first one that all implants, independently from their transports, target OS,
// have in common.
// Generally, security details linked to the transport mechanism are already dealt with.
// The Session paramater is a transport-layer connection, to which we register everything.
func HandleGhostRegistration(sess *Session) {

	// Custom C2 -----------------------------------------------------------------------------------

	// MTLS/DNS read/write loops

	// RPC frameworks ------------------------------------------------------------------------------

	// Base RPC methods allow to exchange registration and information messages
	// client := generic.NewClient(sess)

	// Register RPC services/handlers if the ghost reverse-calls us (we are the server)

	// If bind, then either we wait for registration message to come in, or we request it.

	// This should include the logging infrastructure

	// Implant Registration ------------------------------------------------------------------------

	// Populate new ghostpb object with all registration info, and register user/module interfaces
	// This means, at this point, that although all OS-specific commands are technically available,
	// much of the implant state/information is not disseminated in the ghost object that will be
	// further used by consoles/modules.
	// registrar := &ghostpb.Ghost{}
	// ghosts.NewGhost(registrar)

	// Register/check ghost owner & permissions

	// Register/instantiate/populate OS specific objects in the registered ghost.

	// Send registration notification to user consoles

	// Send all necessary information to Database

	// Network ------------------------------------------------------------------------------------

	// Register implant address to the routing table, for automatic finding of it when we start
	// various handlers, and other needs.
}
