package authentication

import (
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

// Authenticate - Perform full authentication process to server
func Authenticate() (ok bool, user serverpb.User, token string) {

	// Prompt, store and send password (as a hash)

	// Send authentication request

	// Check answer, with success and token

	// If error, try again (five tries)

	// If success, store token, return

	return
}
