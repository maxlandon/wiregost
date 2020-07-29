package security

import "github.com/maxlandon/wiregost/ghost/log"

var (
	secLog = log.GhostLog("security", "core")
)

// SetupSecurity - The entry point to all security checks performed by an implant at startup.
func SetupSecurity() {

	// Check/set limits
}
